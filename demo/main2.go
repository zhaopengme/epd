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
	"time"
)

func main() {
	demo2("bbbbb")
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
	EPD_WIDTH := 800
	EPD_HEIGHT := 480
	var widthByte, heightByte int

	if EPD_WIDTH%8 == 0 {
		widthByte = (EPD_WIDTH / 8)
	} else {
		widthByte = (EPD_WIDTH/8 + 1)
	}

	heightByte = EPD_HEIGHT

	var byteToSend byte = 0x00
	var bgColor = 1

	buffer := bytes.Repeat([]byte{0x00}, widthByte*heightByte)

	for j := 0; j < EPD_HEIGHT; j++ {
		for i := 0; i < EPD_WIDTH; i++ {
			bit := bgColor

			if i < img.Bounds().Dx() && j < img.Bounds().Dy() {
				bit = color.Palette([]color.Color{color.White, color.Black}).Index(img.At(i, j))
			}

			if bit == 1 {
				byteToSend |= 0x80 >> (uint32(i) % 8)
			}

			if i%8 == 7 {
				buffer[(i/8)+(j*widthByte)] = byteToSend
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

	dc.SavePNG("a.png")
	buf := convertImage(dc.Image())

	epd, e := epd.NewRaspberryPiHat()
	if e != nil {
		log.Fatalln(e)
	}
	epd.HardwareInit()
	epd.Clear()

	epd.Display(buf)

	time.Sleep(10 * time.Second)

	epd.Clear()

	time.Sleep(5 * time.Second)

	epd.Sleep()
	return nil
}

func getBuffer(image image.Image) []byte {
	width := 800
	height := 480

	size := (width * height) / 8
	data := make([]byte, size)
	for i := range data {
		data[i] = 255
	}

	imageWidth := image.Bounds().Dx()
	imageHeight := image.Bounds().Dy()

	if imageWidth == width && imageHeight == height {
		fmt.Println("Vertical")
		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {
				if isBlack(image, x, y) {
					shift := uint32(x % 8)
					data[(x+y*width)/8] &= ^(0x80 >> shift)
				}
			}
		}
	} else if imageWidth == height && imageHeight == width {
		fmt.Println("Horizontal")
		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {
				newX := y
				newY := height - x - 1
				if isBlack(image, x, y) {
					shift := uint32(y % 8)
					data[(newX+newY*width)/8] &= ^(0x80 >> shift)
				}
			}
		}
	} else {
		fmt.Println("Invalid image size")
	}
	return data
}

func getRGBA(image image.Image, x, y int) (int, int, int, int) {
	r, g, b, a := image.At(x, y).RGBA()
	r = r / 257
	g = g / 257
	b = b / 257
	a = a / 257

	return int(r), int(g), int(b), int(a)
}

func isBlack(image image.Image, x, y int) bool {
	r, g, b, a := getRGBA(image, x, y)
	offset := 10
	return r < 255-offset && g < 255-offset && b < 255-offset && a > offset
}
