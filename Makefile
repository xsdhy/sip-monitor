


sip1:
	GOOS=linux GOARCH=amd64 go build -o sbc
	scp ./sbc sbc:/data/apps/sbc

sbc:
	GOOS=linux GOARCH=amd64 go build -o sip-monitor
	scp ./sip-monitor sbc:/data/apps/sip-monitor/