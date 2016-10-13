package lcd 

import (
	"image"
	"image/draw"	
	"fmt"
	"github.com/goiot/devices/monochromeoled"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/exp/io/i2c"
	"io/ioutil"
)

func LoadFont(fontFile string) (*truetype.Font, error) {
	b, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	return truetype.Parse(b)
}

func writeText(font *truetype.Font, text string, text2 string, size float64) (image.Image, error) {
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

func Screen(text string, text2 string, font *truetype.Font, size float64) {
	image, err := writeText(font, text, text2, size)
	if err != nil {
		fmt.Println(err)
	}

	d, err := monochromeoled.Open(&i2c.Devfs{Dev: "/dev/i2c-1"})
	if err != nil {
		panic(err)
	}
	defer d.Close()

	if err := d.SetImage(0, 0, image); err != nil {
		panic(err)
	}
	if err := d.Draw(); err != nil {
		panic(err)
	}
}

