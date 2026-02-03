run-eth:
	go build && ./coin-ping --token=ETH,2400,3000 --interval=1

run-eth-bg:
	go build && nohup ./coin-ping --token=ETH,2400,3000 --interval=1 > /dev/null 2>&1 &
