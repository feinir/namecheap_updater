package main

import (
	"flag"
	"namecheapupdater/common"
)

func main() {
	config := &common.Config{}
	configFile := flag.String("config", "", "yaml config file")
	//platform := flag.String("p", "namecheap", "ddns platform")
	host := flag.String("host", "", "host (subdomains) to be updated")
	domain := flag.String("domain", "", "domain name to be updated")
	ip := flag.String("ip", "", "ip addres to be updated")
	password := flag.String("password", "", "password to be used")
	updateFrequency := flag.Float64("freq", 5, "update frequency (minutes, default 5)")
	flag.Parse()
	if len(*configFile) > 0 {
		cfg, err := common.ReadConfig(*configFile)
		if err != nil {
			flag.Usage()
			return
		}
		cfg.ConfigFile = *configFile
		config = cfg
		common.Start(config)
	} else {
		if len(*host) > 0 && len(*domain) > 0 && len(*password) > 0 {
			domainHost := common.DomainHost{}
			domainHost.Host = *host
			domainHost.DomainName = *domain
			domainHost.IP = *ip
			domainHost.Password = *password
			domainHost.UpdateFrequency = *updateFrequency
			config.DomainHosts = append(config.DomainHosts, domainHost)
			common.Start(config)
		} else {
			flag.Usage()
			return
		}
	}
}
