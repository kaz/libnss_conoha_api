package main

import (
	"fmt"
	"net"
)

func main() {
	addrs, err := net.LookupIP("mmk.conoha")
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		fmt.Println(addr.String())
	}
}
