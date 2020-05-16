file=bilibili.out

send: build
	scp ./$(file) root@hd1h.ssh.aaronkir.xyz:bilibili
	scp ./config.json root@hd1h.ssh.aaronkir.xyz:bilibili
build:
	go build -o $(file) *.go
clean:
	rm -vf ./$(file)
