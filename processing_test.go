package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"testing"
)

const KernelRadius = 7

var filePaths = []string{"assets/city.jpg", "assets/office.jpg", "assets/beach.jpg"}
var images []image.Image

type ImageBench struct {
	i             image.Image
	name          string
	width, height int
}

func init() {
	var img image.Image

	for _, path := range filePaths {
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}

		img, err = jpeg.Decode(file)
		if err != nil {
			panic(err)
		}

		images = append(images, img)
	}
}

func BenchmarkConvolveSeq(b *testing.B) {
	var benckmarks []ImageBench

	for i := 0; i < len(filePaths); i++ {
		bounds := images[i].Bounds()

		imgBench := ImageBench{}
		imgBench.i = images[i]
		imgBench.name = filePaths[i]
		imgBench.width = bounds.Dx()
		imgBench.height = bounds.Dy()

		benckmarks = append(benckmarks, imgBench)
	}

	for _, bench := range benckmarks {
		b.Run(fmt.Sprintf("Image %s (%dx%d) - Sequential", bench.name, bench.width, bench.height), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				ConvolveSeq(bench.i, KernelRadius)
			}
		})
	}
}

func BenchmarkConvolveConcurrent(b *testing.B) {
	var benckmarks []ImageBench

	for i := 0; i < len(filePaths); i++ {
		bounds := images[i].Bounds()

		imgBench := ImageBench{}
		imgBench.i = images[i]
		imgBench.name = filePaths[i]
		imgBench.width = bounds.Dx()
		imgBench.height = bounds.Dy()

		benckmarks = append(benckmarks, imgBench)
	}

	for _, bench := range benckmarks {
		b.Run(fmt.Sprintf("Image %s (%dx%d) - Concurrent", bench.name, bench.width, bench.height), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				ConvolveConcurent(bench.i, KernelRadius)
			}
		})
	}
}
