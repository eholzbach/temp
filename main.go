package main

import (
	"flag"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/marpaia/graphite-golang"
	"github.com/quhar/bme280"
	"golang.org/x/exp/io/i2c"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func errFunc(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
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

func sendGraphite(g *graphite.Graphite, temperature float64, humidity float64, pressure float64) {
	name, _ := os.Hostname()

	metric := fmt.Sprintf("temp.%s.%s", name, "temperature")
	a := strconv.FormatFloat(temperature, 'f', 6, 64) 
        g.SimpleSend(metric, a)

	metric = fmt.Sprintf("temp.%s.%s", name, "humidity")
	a = strconv.FormatFloat(humidity, 'f', 6, 64)
	g.SimpleSend(metric, a)

        metric = fmt.Sprintf("temp.%s.%s", name, "pressure")
	a = strconv.FormatFloat(pressure, 'f', 6, 64)
	g.SimpleSend(metric, a)
}

func sensorData(sensor *i2c.Device) (float64, float64, float64, error) {
	b := bme280.New(sensor, bme280.TempUnit(1))
	err := b.Init()

	temperature, pressure, humidity, err := b.EnvData()
	if err != nil {
		return 0, 0, 0, err
	}
	return temperature, pressure, humidity, err
}

func main() {
	// parse flags
	enableOled := flag.Bool("oled", false, "enable SSD1306 OLED")
	enableGraphite := flag.String("graphite", "none", "send to graphite, requires fqdn:port")
	flag.Parse()

	// init sensor
	sensor, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x76)
	errFunc(err)

	// load font if oled is enabled
	var font *truetype.Font

	if *enableOled {
		font, err = loadFont("./visitor1.ttf")
		errFunc(err)
	}

	// init graphite if enabled
	var metricHost string
	var metricPort int 
	var g *graphite.Graphite

	if *enableGraphite != "none" {
		a := strings.Split(*enableGraphite, ":")
		metricHost = a[0]
		metricPort, err = strconv.Atoi(a[1])
		errFunc(err)

		g, err = graphite.NewGraphite(metricHost, metricPort)
	} else {
		g = graphite.NewGraphiteNop(metricHost, metricPort)
	}
	if err != nil {
		g = graphite.NewGraphiteNop(metricHost, metricPort)
	}

	// collect data
	for {
		temperature, pressure, humidity, err := sensorData(sensor)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

		if *enableGraphite != "none" {
			go sendGraphite(g, temperature, humidity, pressure)
		}

		if *enableOled {
			go updateOled(temperature, humidity, pressure, font)
		}

		time.Sleep(10 * time.Second)
	}
}
