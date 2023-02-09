# Waveshare 7in5 v2 Driver

[![Go Reference](https://pkg.go.dev/badge/github.com/NelsonBrandao/waveshare7in5v2.svg)](https://pkg.go.dev/github.com/NelsonBrandao/waveshare7in5v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/NelsonBrandao/waveshare7in5v2)](https://goreportcard.com/report/github.com/NelsonBrandao/waveshare7in5v2)

A driver written in Go for Waveshare 7in5 v2 e-Paper display to be used on a Raspberry Pi board.

Follow the [official documentation](https://www.waveshare.com/wiki/7.5inch_e-Paper_HAT_Manual#Overview) for how to setup and connect the display.

### Usage

#### Epd driver

A simple driver Epd is implemented that closely follows the official C/Python examples provided by waveshare.

```go
// Create the edp instance
epd, err := waveshare7in5v2.New()
if err != nil {
  fmt.Println("Failed to initialize driver:", err)
}

// Init the display
epd.Init()

// Create an image with the same size has the screen
pattern := image.NewRGBA(epd.Bounds())
draw.Draw(pattern, pattern.Bounds(), image.White, image.ZP, draw.Src)
drawer := &font.Drawer{
  Dst:  pattern,
  Src:  image.Black,
  Face: basicfont.Face7x13,
  Dot:  fixed.P(pattern.Bounds().Dx()/2, pattern.Bounds().Dy()/2),
}
drawer.DrawString("Hello World!")

// Display the image on the screen
epd.DisplayImage(pattern)
waitForInput()

// Clear the screen before sleeping
epd.Clear()

// Set the display do deep sleep
epd.Sleep()

// Close the connection and cleanup
epd.Close()
```

#### Canvas
Canvas implements `draw.Image` allowing to use any compatible package to draw directly to the display.

```go
// Create the edp instance
epd, err := waveshare7in5v2.New()
if err != nil {
  fmt.Println("Failed to initialize driver:", err)
}

// Init the display
epd.Init()

// Create the canvas instance
canvas := waveshare7in5v2.NewCanvas(epd);

// Draw directly into the canvas. This will store any changes into a buffer
// until the display is refreshed
draw.Draw(canvas, canvas.Bounds(), image.White, image.ZP, draw.Src)
drawer := &font.Drawer{
  Dst:  canvas,
  Src:  image.Black,
  Face: basicfont.Face7x13,
  Dot:  fixed.P(canvas.Bounds().Dx()/2, canvas.Bounds().Dy()/2),
}
drawer.DrawString("Hello World!")

// Flushes any changes done locally and updates the display
canvas.Refresh()
waitForInput()

// Clear the screen before sleeping
canvas.Clear()

// Set the display do deep sleep
epd.Sleep()

// Close the connection and cleanup
epd.Close()
```
