package main

import (
	"image"
	"image/color"
	"math"
	"sync"
)

const StdDev = float64(10)

type Pixel struct {
	R int
	G int
	B int
	A int
}

func generateKernel(radius int) []float64 {
	kSize := (2 * radius) + 1
	kernel := make([]float64, kSize)

	eFactor := 1 / math.Sqrt(2*math.Pi*(StdDev*StdDev))
	xDiv := -(1 / (2 * StdDev * StdDev))

	sum := float64(0)
	for x := -radius; x <= radius; x++ {
		res := eFactor * math.Pow(math.E, (float64(x*x)*xDiv))
		kernel[x+radius] = res
		sum += res
	}

	for i := 0; i < kSize; i++ {
		kernel[i] /= sum
	}

	return kernel
}

func toPixelArr(img image.Image) [][]Pixel {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func pixelToRgba(pixel Pixel) (uint32, uint32, uint32, uint32) {
	return uint32(pixel.R * 257), uint32(pixel.G * 257), uint32(pixel.B * 257), uint32(pixel.A * 257)
}

func convolvePixel(pixels []Pixel, kernel []float64) Pixel {
	kernelSum := float64(0)
	for _, ki := range kernel {
		kernelSum += ki
	}

	var avgR, avgG, avgB, avgA float64 = 0, 0, 0, 0
	for i := 0; i < len(kernel); i++ {
		avgR += (kernel[i] * float64(pixels[i].R))
		avgG += (kernel[i] * float64(pixels[i].G))
		avgB += (kernel[i] * float64(pixels[i].B))
		avgA += (kernel[i] * float64(pixels[i].R))
	}

	avgR /= kernelSum
	avgG /= kernelSum
	avgB /= kernelSum
	avgA /= kernelSum

	return Pixel{
		R: int(avgR),
		G: int(avgG),
		B: int(avgB),
		A: int(avgA),
	}
}

// ConvolveSeq uses a two pass approach to convolve the pixels
// of an image using a Gaussian kernel. It does it all sequantially
// pixel by pixel, wihtout leveraging goroutines
func ConvolveSeq(img image.Image, radius int) *image.RGBA {
	kernel := generateKernel(radius)
	resultingImg := image.NewRGBA(img.Bounds())

	pixels := toPixelArr(img)

	// First pass
	var pixelBuffer []Pixel
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[0]); x++ {
			// Create the pixel buffer
			for i := -radius; i <= radius; i++ {
				if x+i < 0 || x+i >= len(pixels[0]) {
					pixelBuffer = append(pixelBuffer, Pixel{255, 255, 255, 255})
					continue
				}

				pixelBuffer = append(pixelBuffer, pixels[y][x+i])
			}

			newPixel := convolvePixel(pixelBuffer, kernel)
			c := color.RGBA{
				R: uint8(newPixel.R),
				G: uint8(newPixel.G),
				B: uint8(newPixel.B),
				A: uint8(newPixel.A),
			}

			resultingImg.Set(x, y, c)
			pixelBuffer = nil
		}
	}

	// Second pass
	pixels = toPixelArr(resultingImg)
	for x := 0; x < len(pixels[0]); x++ {
		for y := 0; y < len(pixels); y++ {
			// Create the pixel buffer
			for i := -radius; i <= radius; i++ {
				if y+i < 0 || y+i >= len(pixels) {
					pixelBuffer = append(pixelBuffer, Pixel{255, 255, 255, 255})
					continue
				}

				pixelBuffer = append(pixelBuffer, pixels[y+i][x])
			}

			newPixel := convolvePixel(pixelBuffer, kernel)
			c := color.RGBA{
				R: uint8(newPixel.R),
				G: uint8(newPixel.G),
				B: uint8(newPixel.B),
				A: uint8(newPixel.A),
			}

			resultingImg.Set(x, y, c)
			pixelBuffer = nil
		}
	}

	return resultingImg
}

func ConvolveConcurent(img image.Image, radius int) *image.RGBA {
	kernel := generateKernel(radius)
	resultingImg := image.NewRGBA(img.Bounds())
	pixels := toPixelArr(img)

	var wg sync.WaitGroup

    // Primeira passada: iteramos sobre as linhas e aplicamos a convolução nos
    // pixels
	for y := 0; y < len(pixels); y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			var pixelBuffer []Pixel
			for x := 0; x < len(pixels[0]); x++ {
                // Criar o buffer de pixels que contém os pixels
                // vizinhos de acordo com o raio do kernel
				for i := -radius; i <= radius; i++ {
					if x+i < 0 || x+i >= len(pixels[0]) {
						pixelBuffer = append(pixelBuffer, Pixel{255, 255, 255, 255})
						continue
					}
					pixelBuffer = append(pixelBuffer, pixels[y][x+i])
				}

				newPixel := convolvePixel(pixelBuffer, kernel)
				c := color.RGBA{
					R: uint8(newPixel.R),
					G: uint8(newPixel.G),
					B: uint8(newPixel.B),
					A: uint8(newPixel.A),
				}

				resultingImg.Set(x, y, c)
				pixelBuffer = nil
			}
		}(y)
	}

	wg.Wait()

    // Segunda passada: iteramos sobre as colunas e aplicamos a convolução nos
    // pixels.
	pixels = toPixelArr(resultingImg) // Update pixels after the first pass
	for x := 0; x < len(pixels[0]); x++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			var pixelBuffer []Pixel
			for y := 0; y < len(pixels); y++ {
                // Criar o buffer de pixels que contém os pixels
                // vizinhos de acordo com o raio do kernel
				for i := -radius; i <= radius; i++ {
					if y+i < 0 || y+i >= len(pixels) {
						pixelBuffer = append(pixelBuffer, Pixel{255, 255, 255, 255})
						continue
					}
					pixelBuffer = append(pixelBuffer, pixels[y+i][x])
				}

				newPixel := convolvePixel(pixelBuffer, kernel)
				c := color.RGBA{
					R: uint8(newPixel.R),
					G: uint8(newPixel.G),
					B: uint8(newPixel.B),
					A: uint8(newPixel.A),
				}

				resultingImg.Set(x, y, c)
				pixelBuffer = nil
			}
		}(x)
	}

	wg.Wait()

	return resultingImg
}
