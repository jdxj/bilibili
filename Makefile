file="bilibili.out"

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(file) *.go

send: build
	scp $(file) root@hd1h.ssh.aaronkir.xyz:bilibili

clean:
	rm -vf $(file)
