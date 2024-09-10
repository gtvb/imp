package main

import (
	"image"
	"os"
	"testing"
)

const KernelRadius = 7
const TestFilePath = "assets/office.jpg"

func BenchmarkConvolveSeq(b *testing.B) {
    file, _ := os.Open(TestFilePath)
    defer file.Close()

    image, _, _ := image.Decode(file)

    for n := 0; n < b.N; n++ {
        ConvolveSeq(image, KernelRadius)
    }
}

func BenchmarkConvolveConcurrent(b *testing.B) {
    file, _ := os.Open(TestFilePath)
    defer file.Close()

    image, _, _ := image.Decode(file)

    for n := 0; n < b.N; n++ {
        ConvolveConcurent(image, KernelRadius)
    }
}
