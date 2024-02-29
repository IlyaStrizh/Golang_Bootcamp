package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

const (
	width  = 300
	height = 300
	name   = "amazing_logo.png"
)

func createLogo() error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, fillImage()); err != nil {
		return nil
	}

	return nil
}

func fillImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Заполнение изображения
	patternSize := 38
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if (x/patternSize+y/patternSize)%2 == 0 {
				img.Set(x, y, color.RGBA{210, 180, 140, 255})
				//img.Set(x, y, color.RGBA{140, 180, 210, 255})
			} else {
				img.Set(x, y, color.RGBA{100, 30, 0, 255})
				//img.Set(x, y, color.RGBA{0, 30, 100, 255})
			}
		}
	}

	return img
}

func main() {
	if err := createLogo(); err != nil {
		log.Println(err)
	}
}
