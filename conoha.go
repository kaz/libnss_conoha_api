package main

import (
	"encoding/gob"
	"net"
	"os"
	"strings"
)

const (
	DOMAIN     = ".conoha"
	CACHE_PATH = "/var/cache/nss_conoha_cache"
)

func LookupInstance(host string, v6 bool) []string {
	if !strings.HasSuffix(host, DOMAIN) {
		return nil
	}
	host = strings.TrimSuffix(host, DOMAIN)

	cli := &ConohaClient{}
	cache, _ := os.Open(CACHE_PATH)
	if err := gob.NewDecoder(cache).Decode(cli); err != nil {
		cli, err = NewClient(os.Getenv("NSS_CONOHA_REGION"), os.Getenv("NSS_CONOHA_TENANT_ID"), os.Getenv("NSS_CONOHA_USERNAME"), os.Getenv("NSS_CONOHA_PASSWORD"))
		if err != nil {
			return nil
		}

		cache, _ = os.Create(CACHE_PATH)
		if err := gob.NewEncoder(cache).Encode(cli); err != nil {
			return nil
		}
	}

	servers, err := cli.Servers()
	if err != nil {
		return nil
	}

	result := make([]string, 0, 1)
	for _, srv := range servers {
		if srv.Metadata.Tag == host {
			for _, addrs := range srv.Addresses {
				for _, addr := range addrs {
					ipaddr := net.ParseIP(addr.Addr)
					if !v6 && addr.Version == 4 {
						result = append(result, string(ipaddr[len(ipaddr)-4:]))
					}
					if v6 && addr.Version == 6 {
						result = append(result, string(ipaddr))
					}
				}
			}
		}
	}

	return result
}
