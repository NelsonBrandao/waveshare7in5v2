// Package waveshare7in5v2 implements a driver for the waveshare 7in5 V2 e-Paper display
// to be used on a Raspberry Pi board.
//
// A simple driver Epd is implemented that closely follows the official C/Python examples
// provided by waveshare. There is also a Canvas that implements draw.Image allowing to use
// any compatible package to draw to the display.
//
// Datasheet:   https://www.waveshare.com/w/upload/6/60/7.5inch_e-Paper_V2_Specification.pdf
// C code:      https://github.com/waveshare/e-Paper/blob/master/RaspberryPi_JetsonNano/c/lib/e-Paper/EPD_7in5_V2.c
// Python code: https://github.com/waveshare/e-Paper/blob/master/RaspberryPi_JetsonNano/python/lib/waveshare_epd/epd7in5_V2.py
package waveshare7in5v2

import (
	"image"
	"log"

	"github.com/stianeikeland/go-rpio/v4"
)

// The driver to interact with the e-paper display
type Epd struct {
	dc   rpio.Pin
	cs   rpio.Pin
	rst  rpio.Pin
	busy rpio.Pin

	bounds     image.Rectangle
	bufferSize int
	pixelWidth int

	fasterNormalRefresh bool
}

func New(fasterNormalRefresh bool) (*Epd, error) {
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		return nil, err
	}

	rpio.SpiChipSelect(0)

	dc := rpio.Pin(25)
	cs := rpio.Pin(8)
	rst := rpio.Pin(17)
	busy := rpio.Pin(24)

	dc.Output()
	cs.Output()
	rst.Output()
	busy.Input()

	bounds := image.Rect(0, 0, EPD_WIDTH, EPD_HEIGHT)
	pixelWidth := EPD_WIDTH / PIXEL_SIZE
	bufferSize := pixelWidth * EPD_HEIGHT

	d := &Epd{
		dc:   dc,
		cs:   cs,
		rst:  rst,
		busy: busy,

		bounds:     bounds,
		bufferSize: bufferSize,
		pixelWidth: pixelWidth,

		fasterNormalRefresh: fasterNormalRefresh,
	}

	return d, nil
}

// Powers on the screen after power off or sleep.
func (e *Epd) Init() {
	e.reset()

	e.initDisplay()
}

// Returns the current screen bounds
func (e *Epd) Bounds() image.Rectangle {
	return e.bounds
}

// Converts an image into a buffer array ready to be sent to the display.
// Due to the display only supporting 2 colors a threshold is applied to convert the image to pure black and white.
// The returned buffer is ready to be sent using UpdateFrame.
func (e *Epd) GetBuffer(img image.Image, threshold uint8) []byte {
	buffer := make([]byte, e.bufferSize)

	for y := 0; y < e.bounds.Dy(); y++ {
		for x := 0; x < e.bounds.Dx(); x += PIXEL_SIZE {
			// Start with white
			var pixel byte = 0x00

			// Iterate and append over the next 8 pixels
			for px := 0; px < PIXEL_SIZE; px++ {
				if isBlack(img.At(x+px, y), threshold) {
					pixel |= (0x80 >> px)
				}
			}

			buffer[(y*e.pixelWidth + x/PIXEL_SIZE)] = pixel
		}
	}

	return buffer
}

// Updates the internal display buffer.
func (e *Epd) UpdateFrame(buffer []byte) {
	log.Println("Updating frame")

	e.sendCommandWithData(DISPLAY_START_TRANSMISSION_1, buffer)

	e.sendCommandWithData(DISPLAY_START_TRANSMISSION_2, buffer)
	log.Println("Updating frame. Done")
}

// Refreshes the display by sending the internal buffer to the screen.
func (e *Epd) Refresh() {
	log.Println("Refreshing display")

	if e.fasterNormalRefresh {
		e.setLut()
	} else {
		e.sendCommandWithData(PANEL_SETTING, []byte{0x1f})
	}

	e.sendCommand(DISPLAY_REFRESH)
	wait(100)
	e.waitUntilIdle()

	log.Println("Refreshing display. Done")
}

func (e *Epd) RefreshQuick() {
	log.Println("Refreshing display quick")

	if !e.fasterNormalRefresh {
		e.sendCommandWithData(PANEL_SETTING, []byte{0x3f}) // use custom LUT
	}
	e.setLutQuick()
	e.sendCommand(DISPLAY_REFRESH)
	wait(100)
	e.waitUntilIdle()

	log.Println("Refreshing display quick. Done")
}

func (e *Epd) setLut() {
	e.sendCommandWithData(0x20, getLutVcom())
	e.sendCommandWithData(0x21, getLutWw())
	e.sendCommandWithData(0x22, getLutBw())
	e.sendCommandWithData(0x23, getLutWb())
	e.sendCommandWithData(0x24, getLutBb())
}

func (e *Epd) setLutQuick() {
	e.sendCommandWithData(0x20, getLutVcomFast())
	e.sendCommandWithData(0x21, getLutWwFast())
	e.sendCommandWithData(0x22, getLutBwFast())
	e.sendCommandWithData(0x23, getLutWbFast())
	e.sendCommandWithData(0x24, getLutBbFast())
}

// Updates the internal display buffer and refresh the screen in sequence.
func (e *Epd) UpdateFrameAndRefresh(buffer []byte) {
	e.UpdateFrame(buffer)
	e.Refresh()
}

func (e *Epd) UpdateFrameAndRefreshQuick(buffer []byte) {
	e.UpdateFrame(buffer)
	e.RefreshQuick()
}

// Allows to easily send an image.Image directly to the screen.
func (e *Epd) DisplayImage(img image.Image) {
	buffer := e.GetBuffer(img, 199)

	e.UpdateFrameAndRefresh(buffer)
}

func (e *Epd) DisplayImageQuick(img image.Image) {
	buffer := e.GetBuffer(img, 199)

	e.UpdateFrameAndRefreshQuick(buffer)
}

// Clear the buffer and updates the screen right away.
func (e *Epd) Clear() {
	log.Println("Clearing display")

	var buffer = make([]byte, e.bufferSize)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 0x00
	}

	e.UpdateFrameAndRefresh(buffer)

	log.Println("Clearing display. Done")
}

// Puts the display to sleep and powers off. This helps ensure the display longevity
// since keeping it powered on for long periods of time can damage the screen.
// After Sleep the display needs to be woken up by running Init again
func (e *Epd) Sleep() {
	log.Println("Putting display to sleep")
	e.sendCommand(POWER_OFF)
	e.waitUntilIdle()

	e.sendCommandWithData(DEEP_SLEEP, []byte{0xa5})

	wait(2000)

	log.Println("Putting display to sleep. Done")
}

// Powers off the display and closes the SPI connection.
func (e *Epd) Close() {
	e.cs.Write(rpio.Low)
	e.dc.Write(rpio.Low)
	e.rst.Write(rpio.Low)

	rpio.SpiEnd(rpio.Spi0)
	rpio.Close()
}

func (e *Epd) initDisplay() {
	// According to the spec this should be:
	// 0x17, 0x17, 0x27, 0x27 (Strength 3,3,5,5)
	// According to the C source code this should be:
	// 0x27, 0x27, 0x2f, 0x17 (Strength 5,5,6,3)
	e.sendCommandWithData(BOOSTER_SOFT_START, []byte{
		0x27, // Strength 5
		0x27, // Strength 5
		0x2f, // Strength 6
		0x17, // Strength 3
	})

	e.sendCommandWithData(POWER_SETTING, []byte{
		0x07, // Border LDO disabled
		0x17,
		0x3f,
		0x3f,
	})

	// Source code
	//e.sendCommandWithData(VCOM_DC, []byte{0x24})

	e.sendCommand(POWER_ON)
	wait(100)
	e.waitUntilIdle()

	if e.fasterNormalRefresh {
		e.sendCommandWithData(PANEL_SETTING, []byte{0x3f}) // use custom LUT
	} else {
		e.sendCommandWithData(PANEL_SETTING, []byte{0x1f})
	}

	e.sendCommandWithData(RESOLUTION_SETTING, []byte{0x03, 0x20, 0x01, 0xe0})

	e.sendCommandWithData(DUAL_SPI, []byte{0x00})

	e.sendCommandWithData(TCON, []byte{0x22})

	e.sendCommandWithData(VCOM_DATA_INTERVAL_SETTING, []byte{0x10, 0x07})
}

func (e *Epd) reset() {
	e.rst.Write(rpio.High)
	wait(200)
	e.rst.Write(rpio.Low)
	wait(2)
	e.rst.Write(rpio.High)
	wait(200)
}
