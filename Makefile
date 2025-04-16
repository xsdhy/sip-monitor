build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o sip-monitor

sbc:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o sip-monitor
	scp ./sip-monitor sbc:/data/apps/sip-monitor/sip-monitor.new


sbcmini:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o sip-monitor
	rm sip-monitor-lb
	upx -9 -o sip-monitor-lb sip-monitor 
	scp ./sip-monitor-lb sbc:/data/apps/sip-monitor/sip-monitor.new

cti:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o sip-monitor
	scp ./sip-monitor sip1:/data/apps/sbc/