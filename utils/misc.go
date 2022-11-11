package utils

import (
	"errors"
	"image"
	"image/gif"
	"image/png"
	"io"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

var ErrUnEqualBounds = errors.New("images with unequal bounds")

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

// check an target is in a slice
func IsInSlice[T comparable](elements []T, target T) bool {
	for _, e := range elements {
		if e == target {
			return true
		}
	}

	return false
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

func CompareTwoImage(img1, img2 image.Image) (dist float64, diffcnt float32, err error) {
	if img1.Bounds() != img2.Bounds() {
		err = ErrUnEqualBounds
		return
	}

	imgMin, imgMax := img1.Bounds().Min, img1.Bounds().Max

	var errsum uint64 = 0
	var rc, gc, bc, tr, tg, tb uint64 = 0, 0, 0, 0, 0, 0

	for x := imgMin.X; x < imgMax.X; x++ {
		for y := imgMin.Y; y < imgMax.Y; y++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()

			// sum rgb values
			rdiff := sqrtDiff(r1, r2)
			gdiff := sqrtDiff(g1, g2)
			bdiff := sqrtDiff(b1, b2)
			adiff := sqrtDiff(a1, a2)
			errsum += rdiff + gdiff + bdiff + adiff

			// count diff red
			if rdiff > 0 {
				rc++
			}

			// count diff green
			if gdiff > 0 {
				gc++
			}

			// count diff blue
			if bdiff > 0 {
				bc++
			}

			// total pixel
			tb++
			tr++
			tg++
		}

	}

	dist = math.Sqrt(float64(errsum))
	diffcnt = float32(rc+bc+gc) / float32(tb+tr+tg)

	return
}

func sqrtDiff(x, y uint32) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}
