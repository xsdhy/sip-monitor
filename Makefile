

sbc:
	GOOS=linux GOARCH=amd64 go build -o sip-monitor
	scp ./sip-monitor sbc:/data/apps/sip-monitor/sip-monitor.new

cti:
	GOOS=linux GOARCH=amd64 go build -o sip-monitor
	scp ./sip-monitor sip1:/data/apps/sbc/


qa:
	GOOS=linux GOARCH=amd64 go build -o sip-monitor
	scp ./sip-monitor voiceqa:/data/apps/sip-monitor/