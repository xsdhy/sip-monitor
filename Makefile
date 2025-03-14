



qa:
	GOOS=linux GOARCH=amd64 go build -o sip-monitor
	scp ./sip-monitor voiceqa:/data/apps/sip-monitor/