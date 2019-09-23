package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/01-edu/z01"
)

type pixel struct {
	r, g, b uint32
}

//lightest to darkest for the acsii
const asciiBrightness = "`^\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

//slice for the brightness for the acsii
var (
	brightness     []uint32
	boundX, boundY int
)

func getImage(img string) [][]pixel {
	//create a slice of slices for the pixels (bidimentional slice)
	var images [][]pixel

	//Walk walks the file tree using the path, start from the img and apply the function
	filepath.Walk(img, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		image := loadImage(path)
		pixels := getPixels(image)
		images = append(images, pixels)
		return nil
	})
	return images
}

//decode the image and returns it
func loadImage(filename string) image.Image {
	f, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	//it will exetude this if the function finally returns the statment
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

//return a slice of pixels, fetchs all the pixels
func getPixels(img image.Image) []pixel {
	//fetch the bound of the image to use the height and the width of the image
	bound := img.Bounds()
	boundX = bound.Dx()
	boundY = bound.Dy()

	fmt.Printf("amout of pixels: %d x %d\n", bound.Dx(), bound.Dy())
	dime := bound.Dx() * bound.Dy()

	//the length/dimention of all the pixels in the image
	pixels := make([]pixel, dime)

	for i := 0; i < dime; i++ {
		//the x is a pixel, if the image as 480*503 it will have 241440 pixels
		//every time the i encrements it will give the pixels position
		//if we do 3%480=3 and if we do 483%480=3
		x := i % bound.Dx()
		//the y is the row, if we do 3/480=0.0... and if we do 483/480=1.006.. it's gona give the next row of the image
		y := i / bound.Dx()
		// At returns the color of the pixel at (x, y) of the image
		r, g, b, _ := img.At(x, y).RGBA()

		//average := (r + g + b) / 3
		average := math.Sqrt(0.299*math.Pow(float64(r), 2) + 0.587*math.Pow(float64(g), 2) + 0.114*math.Pow(float64(b), 2)) // <- this is the best way to get the brightness

		brightness = append(brightness, uint32(average/257))
		pixels[i].r = r
		pixels[i].g = g
		pixels[i].b = b
	}
	return pixels
}

func main() {
	getImage("Homer-simpson.jpg")
	br := []rune(asciiBrightness)
	//mapping the acsii with the brigthness using formula
	baseBri := 127
	baseChar := 32

	for i := 0; i < len(brightness); i++ {
		x := i % boundX
		if x == 0 {
			fmt.Println()
		}
		formula := (baseChar * int(brightness[i])) / baseBri
		z01.PrintRune(br[formula])
	}
}
