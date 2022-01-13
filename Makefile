.PHONY: build
build:
	go build .

.PHONY: producer
producer:
	./weather-alert producer

.PHONY: consumer
consumer:
	./weather-alert consumer