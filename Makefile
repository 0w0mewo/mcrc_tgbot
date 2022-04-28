BINARY_NAME=mcrc_tgbot

build:
	GOOS=linux CGO_ENABLE=0 go build -ldflags=" -s -w" -o ${BINARY_NAME} main.go
	upx -9 -f ${BINARY_NAME}

clean:
	go clean || rm -f ${BINARY_NAME}