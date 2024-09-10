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

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	convolved := ConvolveConcurent(img, 14)
	fmt.Printf("Elapsed (ms): %d\n", time.Since(now).Milliseconds())

	output, err := os.Create("output.jpg")
	if err != nil {
		panic(err)
	}

	if err = jpeg.Encode(output, convolved, nil); err != nil {
		panic(err)
	}
}
