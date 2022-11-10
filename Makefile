BINARY_NAME=mcrc_tgbot
OUTPUT_FOLDER=build/
.PHONY: all build
all: build
build:
	GOOS=linux CGO_ENABLE=0 go build -ldflags=" -s -w" -o ${OUTPUT_FOLDER}${BINARY_NAME} main.go
	upx -9 -f ${OUTPUT_FOLDER}${BINARY_NAME}

clean:
	go clean || rm -f ${OUTPUT_FOLDER}${BINARY_NAME}