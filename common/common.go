package common

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var StopFlag bool = false

// var WaitGroup sync.WaitGroup
var config *Config

func CallStop() {
	StopFlag = true
}

func Start(cfg *Config) {
	fmt.Println("Start app")
	config = cfg
	//	defer WaitGroup.Done()
	for i := range config.DomainHosts {
		if config.DomainHosts[i].UpdateFrequency <= 0 {
			config.DomainHosts[i].UpdateFrequency = 5
		}
	}
	for !StopFlag {
		go Update()
		time.Sleep(time.Second * 40)
		// http timeout 30
	}
}

func Update() {
	for i := range config.DomainHosts {
		host := config.DomainHosts[i]
		if time.Now().Sub(host.LastUpdateTime).Minutes() >= host.UpdateFrequency {
			if len(host.LastUpIP) <= 0 {
				//第一次启动
				host = UpdateOne(host)
			} else {
				nowIP := GetHttpIP(host)
				if nowIP != GetDNSIp(&host) {
					host = UpdateOne(host)
				}
			}
			host.LastUpdateTime = time.Now()
			config.DomainHosts[i] = host
		}
	}
}

func UpdateOne(host DomainHost) DomainHost {
	var ip string
	if len(ip) <= 0 {
		ip = GetHttpIP(host)
	} else {
		ip = host.IP
	}
	if len(ip) > 0 {
		host.LastUpIP = ip
		result, err := updateIP(host.Host, host.DomainName, ip, host.Password)
		if err != nil {
			print(result)
			host.LastStatus = false
			return host
		}
		print(result)
		host.LastStatus = true
	} else {
		fmt.Println("error, ip is null")
		host.LastStatus = false
	}
	return host
}

func GetDNSIp(host *DomainHost) string {
	fullHost := host.Host
	if !strings.Contains(fullHost, host.DomainName) {
		fullHost = fullHost + "." + host.DomainName
	}
	return ResolveDomain(fullHost)
}

func ResolveDomain(name string) string {
	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		fmt.Println("Resolution error", err.Error())
		return ""
	}
	fmt.Println("Resolved Address is: ", "["+name+"]"+addr.String())
	return addr.String()
}

func GetHttpIP(host DomainHost) string {
	url := host.GetIpUrl
	if len(url) <= 0 {
		url = config.CommonGetIpUrl
	}
	if len(url) <= 0 {
		url = "http://checkip.synology.com/"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	resp, err := sendRequest(req)
	if err != nil {
		return ""
	}
	result := findIptStr(string(resp))
	println("Find IP Address From HTTP is: ", result+"["+url+"]")
	return result
}

func updateIP(host, domain, ip, password string) (string, error) {
	req, err := createRequest(host, domain, ip, password)
	if err != nil {
		return "request err", err
	}
	result, err := sendRequest(req)
	if err != nil {
		return "send err", err
	}
	return result, nil
}

// https://dynamicdns.park-your-domain.com/update?host=[host]&domain=[domain_name]&password=[ddns_password]&ip=[your_ip]
func createRequest(host, domain, ip, password string) (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://dynamicdns.park-your-domain.com/update", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("host", host)
	q.Add("domain", domain)
	q.Add("password", password)
	q.Add("ip", ip)
	req.URL.RawQuery = q.Encode()
	return req, nil
}

func sendRequest(req *http.Request) (string, error) {
	println("Send data: ", req.URL.String()+"\n")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("ERROR", err)
		}
	}()

	result, err := parse(resp.Body)
	if err != nil {
		return "", err
	}
	println("Recv data: ", result+"\n")
	return result, nil
}

func parse(body io.Reader) (string, error) {
	bodyStr, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(bodyStr), nil
}

func findIptStr(src string) string {
	flysnowRegexp := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
	params := flysnowRegexp.FindStringSubmatch(src)

	for _, param := range params {
		if len(param) > 0 {
			return param
		}
	}
	return ""
}
