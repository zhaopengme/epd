package main

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	epd "github.com/justmiles/epd/lib/epd7in5v2"
	"image"
	"image/color"
	"log"
)

func main() {
	demo1()
}

func demo1() {

	epd, e := epd.NewRaspberryPiHat()
	if e != nil {
		log.Fatalln(e)
	}
	epd.HardwareInit()
	epd.Clear()

	// Create new logo context
	dc := gg.NewContext(800, 480)

	dc.SetColor(color.White)
	dc.DrawRectangle(0, 0, 800, 480)
	dc.Fill()

	dc.DrawCircle(100, 100, 50)
	dc.SetRGB(0, 0, 0)
	dc.Fill()

	buf := convertImage(dc.Image())
	fmt.Println(buf)
	//dc.SavePNG("ab.png")

	epd.Display(buf)

	epd.Sleep()
}

// Convert converts the input image into a ready-to-display byte buffer.
func convertImage(img image.Image) []byte {
	var byteToSend byte = 0x00
	var bgColor = 1

	buffer := bytes.Repeat([]byte{byteToSend}, (800/8)*480)
	max := (800 / 8) * 480
	for j := 0; j < 800; j++ {
		for i := 0; i < 480; i++ {
			bit := bgColor

			if i < img.Bounds().Dx() && j < img.Bounds().Dy() {
				bit = color.Palette([]color.Color{color.White, color.Black}).Index(img.At(i, j))
			}

			if bit == 1 {
				byteToSend |= 0x80 >> (uint32(i) % 8)
			}

			if i%8 == 7 {
				n := (i / 8) + (j * (800 / 8))
				if n >= max {
					n = max - 1
				}
				buffer[n] = byteToSend
				byteToSend = 0x00
			}
		}
	}

	return buffer
}
