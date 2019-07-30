package socks5server

import (
	"io"
	"log"
	"net"
	"strconv"

	"../cidrmatch"
	"../dns"
	"../socks5ToHttp"
)

// ServerSocks5 <--
type ServerSocks5 struct {
	Server             string
	Port               string
	conn               net.Listener
	ToHTTP             bool
	HTTPServer         string
	HTTPPort           string
	Username           string
	Password           string
	ToShadowsocksr     bool
	ShadowsocksrServer string
	ShadowsocksrPort   string
	Socks5Server       string
	Socks5Port         string
	Bypass             bool
	cidrmatch          *cidrmatch.CidrMatch
	CidrFile           string
	DNSServer          string
	dnscache           dns.DnsCache
}

// Socks5 <--
func (socks5Server *ServerSocks5) Socks5() error {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	// socks5Server.dns = map[string]bool{}
	socks5Server.dnscache = dns.DnsCache{
		DNSServer: socks5Server.DNSServer,
	}
	var err error
	socks5Server.cidrmatch, err = cidrmatch.NewCidrMatchWithMap(socks5Server.CidrFile)
	if err != nil {
		return err
	}
	socks5Server.conn, err = net.Listen("tcp", socks5Server.Server+":"+socks5Server.Port)
	if err != nil {
		// log.Panic(err)
		return err
	}

	for {
		client, err := socks5Server.conn.Accept()
		if err != nil {
			// log.Panic(err)
			return err
		}

		go func() {
			// log.Println(runtime.NumGoroutine())
			if client == nil {
				return
			}
			defer client.Close()
			socks5Server.handleClientRequest(client)
		}()
	}
}

func (socks5Server *ServerSocks5) handleClientRequest(client net.Conn) {

	var b [1024]byte
	_, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}

	if b[0] == 0x05 { //只处理Socks5协议
		client.Write([]byte{0x05, 0x00})
		if b[1] == 0x01 {
			// 对用户名密码进行判断
			if b[2] == 0x02 {
				_, err = client.Read(b[:])
				if err != nil {
					log.Println(err)
					return
				}
				username := b[2 : 2+b[1]]
				password := b[3+b[1] : 3+b[1]+b[2+b[1]]]
				if socks5Server.Username == string(username) && socks5Server.Password == string(password) {
					client.Write([]byte{0x01, 0x00})
				} else {
					client.Write([]byte{0x01, 0x01})
					return
				}
			}
		}

		n, err := client.Read(b[:])
		if err != nil {
			log.Println(err)
			return
		}

		var host, port, hostTemplate string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
			hostTemplate = "ip"
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
			hostTemplate = "domain"
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
			hostTemplate = "ip"
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		switch b[1] {
		case 0x01:
			switch socks5Server.Bypass {
			case true:
				// var isMatched bool

				// if _, exist := socks5Server.dns.Load(host); exist == false {
				// 	if hostTemplate != "ip" {
				// 		// ip, err := net.LookupHost(host)
				// 		ip, isSuccess := dns.DNSv4(socks5Server.DNSServer, host)
				// 		if isSuccess == true {
				// 			isMatched = socks5Server.cidrmatch.MatchWithMap(ip[0])
				// 		} else {
				// 			isMatched = false
				// 		}
				// 	} else {
				// 		isMatched = socks5Server.cidrmatch.MatchWithMap(host)
				// 	}
				// 	// if len(socks5Server.dns) > 10000 {
				// 	// 	i := 0
				// 	// 	for key := range socks5Server.dns {
				// 	// 		delete(socks5Server.dns, key)
				// 	// 		i++
				// 	// 		if i > 0 {
				// 	// 			break
				// 	// 		}
				// 	// 	}
				// 	// }
				// 	socks5Server.dns.Store(host, isMatched)
				// 	fmt.Println(runtime.NumGoroutine(), "connect:"+net.JoinHostPort(host, port), isMatched)
				// } else {
				// 	isMatchedTemp, _ := socks5Server.dns.Load(host)
				// 	isMatched = isMatchedTemp.(bool)
				// 	fmt.Println(runtime.NumGoroutine(), "use cache", "connect:"+net.JoinHostPort(host, port), isMatched)
				// }

				switch socks5Server.dnscache.Match(host, hostTemplate, socks5Server.cidrmatch.MatchWithMap) {
				case false:
					if socks5Server.ToHTTP == true {
						socks5Server.toHTTP(client, host, port)
					} else if socks5Server.ToShadowsocksr == true {
						socks5Server.toSocks5(client, net.JoinHostPort(host, port), b[:n])
					} else {
						socks5Server.toTCP(client, net.JoinHostPort(host, port))
					}
				case true:
					socks5Server.toTCP(client, net.JoinHostPort(host, port))
				}

			case false:
				if socks5Server.ToHTTP == true {
					socks5Server.toHTTP(client, host, port)
				} else if socks5Server.ToShadowsocksr == true {
					socks5Server.toSocks5(client, net.JoinHostPort(host, port), b[:n])
				} else {
					socks5Server.toTCP(client, net.JoinHostPort(host, port))
				}
			}

		case 0x02:
			log.Println("bind 请求 " + net.JoinHostPort(host, port))

		case 0x03:
			log.Println("udp 请求 " + net.JoinHostPort(host, port))
			socks5Server.udp(client, net.JoinHostPort(host, port))
		}
	}
}

func (socks5Server *ServerSocks5) connect() {
	// do something
}

func (socks5Server *ServerSocks5) udp(client net.Conn, domain string) {
	// log.Println()
	server, err := net.Dial("udp", domain)
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
	//进行转发
	// httpConnect := make([]byte, 1024)
	// n, _ := client.Read(httpConnect[:])
	// log.Println(string(httpConnect))
	// server.Write(httpConnect[:n])
	go io.Copy(server, client)
	io.Copy(client, server)

}

func (socks5Server *ServerSocks5) toTCP(client net.Conn, domain string) {
	server, err := net.Dial("tcp", domain)
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
	//进行转发
	// httpConnect := make([]byte, 1024)
	// n, _ := client.Read(httpConnect[:])
	// log.Println(string(httpConnect))
	// server.Write(httpConnect[:n])
	go io.Copy(server, client)
	io.Copy(client, server)
}

func (socks5Server *ServerSocks5) toHTTP(client net.Conn, host, port string) {
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
	server, err := net.Dial("tcp", socks5Server.HTTPServer+":"+socks5Server.HTTPPort)
	if err != nil {
		log.Println(err)
	}
	defer server.Close()
	// if port == "443" {
	server.Write([]byte("CONNECT " + host + ":" + port + " HTTP/1.1\r\n\r\n"))
	httpConnect := make([]byte, 1024)
	server.Read(httpConnect[:])
	log.Println(string(httpConnect))
	// }
	// n, _ := client.Read(httpConnect[:])
	// log.Println(string(httpConnect))
	// server.Write(httpConnect[:n])
	go io.Copy(server, client)
	io.Copy(client, server)
}

func (socks5Server *ServerSocks5) toShadowsocksr(client net.Conn) {
	server, err := net.Dial("tcp", socks5Server.ShadowsocksrServer+":"+socks5Server.ShadowsocksrPort)
	if err != nil {
		log.Println(err)
	}
	defer server.Close()
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
	// 转发
	// httpConnect := make([]byte, 1024)
	// n, _ := client.Read(httpConnect[:])
	// log.Println(string(httpConnect))
	// server.Write(httpConnect[:n])
	go io.Copy(server, client)
	io.Copy(client, server)
}

func (socks5Server *ServerSocks5) toSocks5(client net.Conn, host string, b []byte) {
	socks5Conn, err := (&socks5ToHttp.Socks5Client{
		Server:  socks5Server.Socks5Server,
		Port:    socks5Server.Socks5Port,
		Address: host}).NewSocks5ClientOnlyFirstVerify()
	if err != nil {
		log.Println(err)
		socks5Server.toTCP(client, host)
		return
	}

	defer socks5Conn.Close()
	socks5Conn.Write(b)
	// client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
	// 转发
	// httpConnect := make([]byte, 1024)
	// n, _ := client.Read(httpConnect[:])
	// log.Println(string(httpConnect))
	// server.Write(httpConnect[:n])

	go io.Copy(client, socks5Conn)
	io.Copy(socks5Conn, client)
}