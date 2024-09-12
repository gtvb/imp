package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"time"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
    defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

    now := time.Now()
	convolved := ConvolveConcurent(img, 14)
    fmt.Printf("Elapsed: %f (seconds)\n", time.Since(now).Seconds())

	output, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}

	if err = jpeg.Encode(output, convolved, nil); err != nil {
		panic(err)
	}
}
