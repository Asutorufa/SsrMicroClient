package config

import (
	"encoding/json"
	"os"
)

var (
	pathSeparator = string(os.PathSeparator)
	ConPath       = Path + pathSeparator + "yuhaiinConfig.json"
)

// Setting setting json struct
type Setting struct {
	PythonPath         string `json:"pythonPath"`
	SsrPath            string `json:"ssrPath"`
	PidFile            string `json:"pidPath"`
	LogFile            string `json:"logPath"`
	FastOpen           bool   `json:"fastOpen"`
	Works              string `json:"works"`
	LocalAddress       string `json:"localAddress"`
	LocalPort          string `json:"localPort"`
	TimeOut            string `json:"timeOut"`
	HttpProxy          bool   `json:"httpProxy"`
	Bypass             bool   `json:"bypass"`
	HttpProxyAddress   string `json:"httpProxyAddress"`
	Socks5ProxyAddress string `json:"socks5ProxyAddress"`
	RedirProxyAddress  string `json:"redir_proxy_address"`
	BypassFile         string `json:"bypassFile"`
	DnsServer          string `json:"dnsServer"`
	UdpTrans           bool   `json:"udpTrans"`
	AutoStartSsr       bool   `json:"autoStartSsr"`
	IsPrintLog         bool   `json:"is_print_log"`
	IsDNSOverHTTPS     bool   `json:"is_dns_over_https"`
	DNSAcrossProxy     bool   `json:"dns_across_proxy"`
	UseLocalDNS        bool   `json:"use_local_dns"`
	BlackIcon          bool   `json:"black_icon"`
}

// SettingInitJSON init setting json file
func SettingInitJSON() error {
	pa := &Setting{
		AutoStartSsr:       true,
		BypassFile:         Path + string(os.PathSeparator) + "yuhaiin.conf",
		DnsServer:          "1.0.0.1:53",
		Bypass:             true,
		HttpProxy:          true,
		HttpProxyAddress:   "127.0.0.1:8188",
		Socks5ProxyAddress: "127.0.0.1:1080",
		RedirProxyAddress:  "127.0.0.1:8088",
		IsPrintLog:         false,
		IsDNSOverHTTPS:     false,
		DNSAcrossProxy:     false,
		SsrPath:            " ",
		BlackIcon:          false,

		// not use now
		PythonPath:   "",
		LocalAddress: "0.0.0.0",
		LocalPort:    "0",
		TimeOut:      "1000",
		UseLocalDNS:  false,
		UdpTrans:     true,
		PidFile:      "",
		LogFile:      "",
		FastOpen:     true,
		Works:        "8",
	}
	if err := SettingEnCodeJSON(pa); err != nil {
		return err
	}
	return nil
}

// SettingDecodeJSON decode setting json to struct
func SettingDecodeJSON() (*Setting, error) {
	pa := &Setting{}
	file, err := os.Open(ConPath)
	if err != nil {
		return &Setting{}, err
	}
	if json.NewDecoder(file).Decode(&pa) != nil {
		return &Setting{}, err
	}
	return pa, nil
}

// SettingEnCodeJSON encode setting struct to json
func SettingEnCodeJSON(pa *Setting) error {
	file, err := os.Create(ConPath)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")
	if err := enc.Encode(&pa); err != nil {
		return err
	}
	return nil
}
