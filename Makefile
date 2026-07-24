.PHONY: predict
predict:
	go build -ldflags "-s -w" -o ./bin ./predict
