package utils

import (
	"image/gif"
	"image/png"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

func WaitForSignal() chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	return ch
}

// random integer between
func RandomBetween(start int, end int) int {
	if end < start {
		return rand.Intn(start-end) + end
	}

	return rand.Intn(end-start) + start
}

// random choice element
func RandChoice[T any](l []T) T {
	randIndex := RandomBetween(0, len(l)-1)

	return l[randIndex]

}

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}

	return logger

}

func PngToGif(gifWriter io.Writer, pngReader io.Reader) error {
	pngImg, err := png.Decode(pngReader)
	if err != nil {
		return err
	}

	return gif.Encode(gifWriter, pngImg, nil)
}

func StringToBoolean(b string) bool {
	b = strings.ToLower(b)

	switch b {
	case "true":
		return true
	case "false":
		return false

	}

	return false

}
