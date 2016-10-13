package main

import (
	"fmt"
	"github.com/eholzbach/temp/lcd"
	"github.com/quhar/bme280"
	"golang.org/x/exp/io/i2c"
	"math"
	"time"
	"flag"
	"strconv"
)

func main() {

	const temperature string = "temperature"
	const humidity string = "humidity"
	const pressure string = "pressure"

	fontFile := flag.String("font", "./visitor1.ttf", "derp")
	font, _ := lcd.LoadFont(*fontFile)	

	device, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x76)
	if err != nil {
		fmt.Println(err)
		return
	}

	var header string

	for {
		b := bme280.New(device, bme280.TempUnit(1))
		err = b.Init()
		t, p, h, err := b.EnvData()
		if err != nil {
			fmt.Println(err)
			return
		}

		header = "temperature"
		lcd.Screen(header, round(t), font, 90)
		fmt.Println(round(t))
		time.Sleep(4 * time.Second)

		header = "humidity"
		lcd.Screen(header, round(h), font, 90)
		fmt.Println(round(h))
		time.Sleep(4 * time.Second)

		header = "pressure"
		lcd.Screen(header, round(p), font, 60)
		fmt.Println(round(p))
		time.Sleep(3 * time.Second)

	}
}

func round(a float64) string {
	if a < 0 {
		b := (math.Ceil(a - 0.5))
		return strconv.FormatFloat(b, 'f', 1, 64)
	}
	c := int(math.Floor(a + 0.5))
	d := strconv.Itoa(c)
	return d
}
