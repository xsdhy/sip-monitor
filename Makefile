


sip1:
	GOOS=linux GOARCH=amd64 go build -o sbc
	scp ./sbc sip1:/data/apps/sbc