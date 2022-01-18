.PHONY: build
build:
	go build .

.PHONY: producer
producer: build
	./weather-alert producer

.PHONY: consumer
consumer: build
	./weather-alert consumer