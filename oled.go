package main

import (
	"fmt"
	"github.com/goiot/devices/monochromeoled"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/exp/io/i2c"
	"image"
	"image/draw"
	"io/ioutil"
	"os"
	"time"
)

func loadFont(fontFile string) (*truetype.Font, error) {
	a, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	return truetype.Parse(a)
}

func drawText(font *truetype.Font, text string, text2 string, size float64) (image.Image, error) {
	dst := image.NewRGBA(image.Rect(0, 0, 128, 64))
	draw.Draw(dst, dst.Bounds(), image.Black, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.White)
	c.SetFont(font)
	c.SetFontSize(size)

	_, err := c.DrawString(text2, freetype.Pt(4, 65))
	if err != nil {
		return nil, err
	}

	c.SetFontSize(20)

	_, err = c.DrawString(text, freetype.Pt(0, 15))
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func printDisplay(text string, text2 string, font *truetype.Font, size float64) {
	image, err := drawText(font, text, text2, size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	d, err := monochromeoled.Open(&i2c.Devfs{Dev: "/dev/i2c-1"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

	defer d.Close()

	if err := d.SetImage(0, 0, image); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	if err := d.Draw(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}

func updateOled(temperature float64, humidity float64, pressure float64, font *truetype.Font) {
	printDisplay("temperature", round(temperature), font, 90)
	time.Sleep(6 * time.Second)

	printDisplay("humidity", round(humidity), font, 90)
	time.Sleep(2 * time.Second)

	printDisplay("pressure", round(pressure), font, 60)
	time.Sleep(2 * time.Second)
}
