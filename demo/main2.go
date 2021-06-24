package main

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	epd "github.com/justmiles/epd/lib/epd7in5v2"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"log"
	"strings"
)

func main() {
	demo2("adfasdf")
	//demo1()
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

// DisplayText accepts a string text and displays it on the screen
func demo2(text string) error {

	width := 800
	height := 480

	// Create new logo context
	dc := gg.NewContext(width, height)

	// Set Background Color
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Set font color
	dc.SetColor(color.Black)

	dc.Fill()
	dc.SetRGB(0, 0, 0)

	var (
		maxWidth, maxHeight           float64 = float64(width), float64(height)
		fontSize                      float64 = 300  // initial font size
		fontSizeReduction             float64 = 0.95 // reduce the font size by this much until message fits in the display
		fontSizeMinimum               float64 = 10   // Smallest font size before giving up
		lineSpacing                   float64 = 1
		measuredWidth, measuredHeight float64
	)

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return err
	}
	for {
		face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
		dc.SetFontFace(face)

		stringLines := dc.WordWrap(text, maxWidth)

		measuredWidth, measuredHeight = dc.MeasureMultilineString(strings.Join(stringLines, "\n"), lineSpacing)

		// If the message fits within the frame, let's break. Otherwise reduce the font size and try again
		if measuredWidth < maxWidth && measuredHeight <= maxHeight {
			break
		} else {
			fontSize = fontSize * fontSizeReduction
		}

		if fontSize < fontSizeMinimum {
			return fmt.Errorf("unable to fit text on screen: \n %s", text)
		}
		// TODO: debug logging: fmt.Printf("font size: %v\n", fontSize)
	}

	dc.DrawStringWrapped(text, 0, (maxHeight-measuredHeight)/2-(fontSize/4), 0, 0, maxWidth, lineSpacing, gg.AlignCenter)
	buf := convertImage(dc.Image())

	epd, e := epd.NewRaspberryPiHat()
	if e != nil {
		log.Fatalln(e)
	}
	epd.HardwareInit()
	epd.Clear()


	epd.Display(buf)

	epd.Sleep()
	return nil
}
