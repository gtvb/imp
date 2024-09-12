package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"time"
)

func main() {
    var kernelRadius int
    var concurrent bool

    flag.IntVar(&kernelRadius, "r", 7, "Raio do kernel de convolução")
    flag.BoolVar(&concurrent, "c", false, "Habilita a convolução concorrente")
    flag.Parse()

    args := flag.Args()

	file, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
    defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

    var processedImage *image.RGBA
    now := time.Now()
    if concurrent {
        fmt.Println("Utilizando convolução concorrente...")
	    processedImage = ConvolveConcurent(img, kernelRadius)
    } else {
        fmt.Println("Utilizando convolução sequencial...")
	    processedImage = ConvolveSeq(img, kernelRadius)
    }
    fmt.Printf("Duração: %f (segundos)\n", time.Since(now).Seconds())

	output, err := os.Create(args[1])
	if err != nil {
		panic(err)
	}

	if err = jpeg.Encode(output, processedImage, nil); err != nil {
		panic(err)
	}
}
