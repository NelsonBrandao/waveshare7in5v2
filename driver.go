package waveshare7in5v2

import (
	"image"
	"log"

	"github.com/stianeikeland/go-rpio/v4"
)

type Epd struct {
	dc   rpio.Pin
	cs   rpio.Pin
	rst  rpio.Pin
	busy rpio.Pin

	bounds     image.Rectangle
	bufferSize int
	pixelWidth int
}

func New() (*Epd, error) {
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
	}

	return d, nil
}

func (e *Epd) Init() {
	e.reset()

	e.initDisplay()
}

func (e *Epd) Bounds() image.Rectangle {
	return e.bounds
}

func (e *Epd) GetBuffer(img image.Image) []byte {
	buffer := make([]byte, e.bufferSize)

	for y := 0; y < e.bounds.Dy(); y++ {
		for x := 0; x < e.bounds.Dx(); x += PIXEL_SIZE {
			// Start with back
			var pixel byte = 0x00

			// Iterate and append over the next 8 pixels
			for px := 0; px < PIXEL_SIZE; px++ {
				if isWhite(img.At(x+px, y)) {
					pixel |= (0x80 >> px)
				}
			}

			buffer[(y*e.pixelWidth + x/PIXEL_SIZE)] = pixel
		}
	}

	return buffer
}

func (e *Epd) UpdateFrame(buffer []byte) {
	log.Println("Updating frame")

	e.sendCommandWithData(DISPLAY_START_TRANSMISSION_1, buffer)

	e.sendCommandWithData(DISPLAY_START_TRANSMISSION_2, buffer)
	log.Println("Updating frame. Done")
}

func (e *Epd) Refresh() {
	log.Println("Refreshing display")

	e.sendCommand(DISPLAY_REFRESH)
	wait(100)
	e.waitUntilIdle()

	log.Println("Refreshing display. Done")
}

func (e *Epd) UpdateFrameAndRefresh(buffer []byte) {
	e.UpdateFrame(buffer)
	e.Refresh()
}

func (e *Epd) DisplayImage(img image.Image) {
	buffer := e.GetBuffer(img)

	e.UpdateFrameAndRefresh(buffer)
}

func (e *Epd) Clear() {
	log.Println("Clearing display")

	var buffer = make([]byte, e.bufferSize)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 0x00
	}

	e.UpdateFrameAndRefresh(buffer)

	log.Println("Clearing display. Done")
}

func (e *Epd) Sleep() {
	log.Println("Putting display to sleep")
	e.sendCommand(POWER_OFF)
	e.waitUntilIdle()

	e.sendCommandWithData(DEEP_SLEEP, []byte{0xa5})

	wait(2000)

	log.Println("Putting display to sleep. Done")
}

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
	// e.sendCommandWithData(VCOM_DC, []byte{0x24})

	e.sendCommand(POWER_ON)
	wait(100)
	e.waitUntilIdle()

	e.sendCommandWithData(PANEL_SETTING, []byte{0x1f})

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
