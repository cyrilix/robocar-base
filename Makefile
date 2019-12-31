.PHONY: test

test:
	go test -race -mod vendor ./cli ./mqttdevice ./service ./testtools ./types

