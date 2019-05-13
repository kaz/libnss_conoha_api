install: libnss_conoha.so.2
	install -m 0755 $< /lib

libnss_conoha.so.2: *.go
	go build -buildmode=c-shared -o $@ *.go
	strip $@

test:
	go run *.go

clean:
	rm -f libnss_conoha.so.*
