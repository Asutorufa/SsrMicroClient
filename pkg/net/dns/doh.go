package dns

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/proxy"

	"github.com/Asutorufa/yuhaiin/pkg/net/utils"
)

type DoH struct {
	DNS
	*utils.ClientUtil

	Subnet *net.IPNet
	Proxy  func(domain string) (net.Conn, error)

	host string
	port string
	url  string

	cache      *utils.LRU
	httpClient *http.Client
}

func NewDoH(host string, subnet *net.IPNet, p proxy.Proxy) DNS {
	if subnet == nil {
		_, subnet, _ = net.ParseCIDR("0.0.0.0/0")
	}
	dns := &DoH{
		Subnet: subnet,
		cache:  utils.NewLru(200, 20*time.Minute),
	}

	dns.setServer(host)

	if p == nil {
		dns.setProxy(func(s string) (net.Conn, error) {
			return dns.ClientUtil.GetConn()
		})
	} else {
		dns.setProxy(p.Conn)
	}

	return dns
}

// LookupIP .
// https://tools.ietf.org/html/rfc8484
func (d *DoH) LookupIP(domain string) (ip []net.IP, err error) {
	if x, _ := d.cache.Load(domain); x != nil {
		return x.([]net.IP), nil
	}
	if ip, err = d.search(domain); len(ip) != 0 {
		d.cache.Add(domain, ip)
	}
	return
}

func (d *DoH) search(domain string) ([]net.IP, error) {
	DNS, err := dnsCommon(
		domain,
		d.Subnet,
		func(data []byte) ([]byte, error) {
			return d.post(data)
		},
	)
	if err != nil || len(DNS) == 0 {
		return nil, fmt.Errorf("doh resolve domain %s failed: %v", domain, err)
	}
	return DNS, nil
}

func (d *DoH) setServer(host string) {
	d.url = "https://" + host
	uri, err := url.Parse("//" + host)
	if err != nil {
		d.host = host
		d.port = "443"
	} else {
		d.host = uri.Hostname()
		d.port = uri.Port()
		if d.port == "" {
			d.port = "443"
		}
		if uri.Path == "" {
			d.url += "/dns-query"
		}
	}

	d.ClientUtil = utils.NewClientUtil(d.host, d.port)
}

func (d *DoH) setProxy(p func(string) (net.Conn, error)) {
	d.Proxy = p
	d.httpClient = &http.Client{
		Transport: &http.Transport{
			//Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				switch network {
				case "tcp":
					return d.Proxy(addr)
				default:
					return net.Dial(network, addr)
				}
			},
			DisableKeepAlives: false,
		},
		Timeout: 10 * time.Second,
	}
}

func (d *DoH) get(dReq []byte) (body []byte, err error) {
	query := strings.Replace(base64.URLEncoding.EncodeToString(dReq), "=", "", -1)
	urls := d.url + "?dns=" + query
	res, err := d.httpClient.Get(urls)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return
}

// https://www.cnblogs.com/mafeng/p/7068837.html
func (d *DoH) post(dReq []byte) (body []byte, err error) {
	resp, err := d.httpClient.Post(d.url, "application/dns-message", bytes.NewReader(dReq))
	if err != nil {
		return nil, fmt.Errorf("doh post failed: %v", err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("doh read body failed: %v", err)
	}
	return
}

func (d *DoH) Resolver() *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(context.Context, string, string) (net.Conn, error) {
			return dohDial(d.url, d.httpClient), nil
		},
	}
}

type dohResolverDial struct {
	host       string
	deadline   time.Time
	buffer     *bytes.Buffer
	httpClient *http.Client
}

func dohDial(host string, client *http.Client) net.Conn {
	return &dohResolverDial{
		host:       host,
		buffer:     bytes.NewBuffer(nil),
		httpClient: client,
	}
}

func (d *dohResolverDial) Write(data []byte) (int, error) {
	return d.WriteTo(data, nil)
}

func (d *dohResolverDial) Read(data []byte) (int, error) {
	n, err := d.buffer.Read(data)
	return n, err
}

func (d *dohResolverDial) WriteTo(data []byte, _ net.Addr) (int, error) {
	if time.Now().After(d.deadline) {
		return 0, fmt.Errorf("timeout")
	}

	resp, err := d.httpClient.Post(d.host, "application/dns-message", bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("post failed: %v", err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read resp body failed: %v", err)
	}
	defer resp.Body.Close()

	d.buffer.Truncate(0)
	_, err = d.buffer.Write(res)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (d *dohResolverDial) ReadFrom(data []byte) (n int, addr net.Addr, err error) {
	n, err = d.buffer.Read(data)
	return n, nil, err
}

func (d *dohResolverDial) Close() error {
	return nil
}

func (d *dohResolverDial) SetDeadline(t time.Time) error {
	d.deadline = t
	return nil
}

func (d *dohResolverDial) SetReadDeadline(t time.Time) error {
	return nil
}

func (d *dohResolverDial) SetWriteDeadline(t time.Time) error {
	return nil
}

func (d *dohResolverDial) LocalAddr() net.Addr {
	return nil
}
func (d *dohResolverDial) RemoteAddr() net.Addr {
	return nil
}
