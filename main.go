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

//darkest to lightest for the acsii
const (
	asciiBrightness = "`^\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
)

//slice for the brightness for the acsii
var (
	brightness []uint32
	boundX     int
	arg        map[string]bool
	color      map[string]string
)

func getImage(img string) {
	//Walk walks the file tree using the path, start from the img and apply the function
	filepath.Walk(img, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		image := loadImage(path)
		getPixels(image)
		return nil
	})
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
func getPixels(img image.Image) {
	//fetch the bound of the image to use the height and the width of the image
	bound := img.Bounds()
	boundX = bound.Dx()

	//fmt.Printf("amout of pixels: %d x %d\n", bound.Dx(), bound.Dy())
	//the length/dimention of all the pixels in the image
	dime := bound.Dx() * bound.Dy()

	for i := 0; i < dime; i++ {
		//the x is a pixel, if the image as 480*503 it will have 241440 pixels
		//every time the i encrements it will give the pixels position
		//if we do 3%480=3 and if we do 483%480=3
		x := i % bound.Dx()
		//the y is the row, if we do 3/480=0.0... and if we do 483/480=1.006.. it's gona give the next row of the image
		y := i / bound.Dx()
		// At returns the color of the pixel at (x, y) of the image
		r, g, b, _ := img.At(x, y).RGBA()

		//average := (r + g + b) / 3 <- other way but not that good
		average := math.Sqrt(0.299*math.Pow(float64(r), 2) + 0.587*math.Pow(float64(g), 2) + 0.114*math.Pow(float64(b), 2)) // <- this is the best way to get the brightness

		brightness = append(brightness, uint32(average/257))
	}
}

func colors() {
	if color == nil {
		color = make(map[string]string)
	}
	color["blue"] = "\033[0;34m"
	color["red"] = "\033[0;31m"
	color["green"] = "\033[0;32m"
	color["purple"] = "\033[0;35m"
	color["brown"] = "\033[0;33m"

}
func reverseColor(value int) int {
	a := 33
	b := 32
	valueReverse := (value * b) / a
	return valueReverse
}

func argVerification(baseBri, baseChar, i int) (int, bool) {
	for k, v := range arg {
		if v == true && k == "--up" {
			return (baseChar * int(brightness[len(brightness)-1-i])) / baseBri, false
		} else if v == true && k == "--reverseColor" {
			value := (baseChar * int(brightness[i])) / baseBri
			return reverseColor(value), false
		} else if v == true && k == "--color" {
			colors()
			return (baseChar * int(brightness[i])) / baseBri, true
		}
	}
	return (baseChar * int(brightness[i])) / baseBri, false
}

func main() {
	if len(os.Args) >= 2 {
		if arg == nil {
			arg = make(map[string]bool)
			arg["--up"] = false
			arg["--reverseColor"] = false
			arg["--color"] = false
		}
		for i := 1; i < len(os.Args); i++ {
			arg[os.Args[i]] = true
		}
	}
	getImage("b.jpg")
	br := []rune(asciiBrightness)
	colorArg := os.Args[len(os.Args)-1]
	//mapping the acsii with the brigthness using formula
	//this will be the middle of the brightness slice and the middle of the asciiBrightness slice
	baseBri := 127
	baseChar := 32

	for i := 0; i < len(brightness); i++ {
		x := i % boundX
		if x == 0 {
			z01.PrintRune('\n')
		}
		formula, ok := argVerification(baseBri, baseChar, i)
		if ok == true {
			fmt.Print(color[colorArg] + string(br[formula]))
			fmt.Print(color[colorArg] + string(br[formula]))
		} else {
			z01.PrintRune(br[formula])
			z01.PrintRune(br[formula])
		}
	}
}
