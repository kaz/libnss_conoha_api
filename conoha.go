package main

import (
	"strings"
)

func lookup(host string, v6 bool) []string {
	if !strings.HasSuffix(host, ".conoha") {
		return nil
	}
	if v6 {
		return []string{
			string([]byte{1, 1, 4, 5, 1, 4, 1, 9, 1, 9, 8, 1, 0, 8, 9, 3}),
			string([]byte{1, 1, 4, 5, 1, 4, 1, 9, 1, 9, 8, 1, 0, 8, 9, 30}),
		}
	}
	return []string{
		string([]byte{1, 14, 5, 14}),
		string([]byte{1, 14, 5, 140}),
	}
}
