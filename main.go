package main

import (
	"fmt"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"go.uber.org/zap"
)

var (
	fullLogger          *zap.Logger
	logger              *zap.SugaredLogger
	name                = "GOIR"
	window              *sdl.Window
	renderer            *sdl.Renderer
	font                *ttf.Font
	winWidth, winHeight int32 = 1920, 1080
	days                      = [...]string{"Måndag", "Tisdag", "Onsdag", "Torsdag", "Fredag", "Lördag", "Söndag"}
	months                    = [...]string{"januari", "februari", "mars", "april", "maj", "juni", "juli", "augusti", "september", "oktober", "november", "december"}
	temperatureOut            = "0"
)

func init() {
	fullLogger, _ = zap.NewProduction()
	defer fullLogger.Sync() // flushes buffer, if any
	logger = fullLogger.Sugar()
}

func main() {
	logger.Infof("Starting goir...")
	_ = createAndStartMQTT("192.168.1.3:1883", "shiprock", "hass")

	initGraphics()
	items := createItems()
	run(items)

	renderer.Destroy()
	window.Destroy()

}

func run(items []item) {
	var runFlag = true
	var event sdl.Event
	surfacesToFree := make([]*sdl.Surface, 0)
	texturesToDestroy := make([]*sdl.Texture, 0)

	for runFlag {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				logger.Infof("key = %v", t.Keysym)
				if t.Keysym.Scancode == 41 {
					runFlag = false
				}
			}
		}

		renderer.Clear()
		for _, item := range items {
			surface, rect, err := item.getSurfaceAndRect()
			if err != nil {
				logger.Error(err)
			}
			texture, err := renderer.CreateTextureFromSurface(surface)
			renderer.Copy(texture, nil, rect)
			surfacesToFree = append(surfacesToFree, surface)
			texturesToDestroy = append(texturesToDestroy, texture)

		}
		renderer.Present()

		for _, surface := range surfacesToFree {
			surface.Free()
		}
		for _, texture := range texturesToDestroy {
			texture.Destroy()
		}

		sdl.Delay(500)
		surfacesToFree = nil
		texturesToDestroy = nil
	}
}

func initGraphics() {
	sdl.Init(sdl.INIT_EVERYTHING)
	var err error

	if err = ttf.Init(); err != nil {
		logger.Fatal("Failed to initialize TTF: %s\n", err)
	}

	if window, err = sdl.CreateWindow(name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_FULLSCREEN_DESKTOP); err != nil {
		logger.Fatalf("Failed to create window: %s\n", err)
	}

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE); err != nil {
		logger.Fatalf("Failed to create renderer: %v\n", err)
	}

	if font, err = ttf.OpenFont("fonts/Signika-Regular.ttf", 140); err != nil {
		logger.Fatalf("Failed to open font: %v\n", err)
	}

	sdl.ShowCursor(sdl.DISABLE)

}

func format(t time.Time) string {
	return fmt.Sprintf("%s %02d %s %d", days[t.Weekday()-1], t.Day(), months[t.Month()-1], t.Year())
}

func getFullDate() string {
	t := time.Now()
	return format(t)
}

func getTime() string {
	t := time.Now()
	return t.Format("15.04")
}

func getTempOut() string {
	return temperatureOut + "°"
}

func getColor() sdl.Color {
	return sdl.Color{R: 255, G: 255, B: 200, A: 200}
}

func getMixedColor() sdl.Color {
	c1 := sdl.Color{R: 255, G: 0, B: 0, A: 0}
	c2 := sdl.Color{R: 0, G: 0, B: 255, A: 0}
	return blendColor(c1, c2, 0.5)
}

func createItems() []item {
	items := make([]item, 0)
	dateItem := &textItem{
		color: getColor,
		text:  getFullDate,
		x:     25,
		y:     25,
	}
	items = append(items, dateItem)

	timeItem := &textItem{
		color: getColor,
		text:  getTime,
		x:     25,
		y:     180,
	}
	items = append(items, timeItem)

	tempItem := &textItem{
		color: getMixedColor,
		text:  getTempOut,
		x:     25,
		y:     360,
	}
	items = append(items, tempItem)

	return items
}

type color struct {
	R, G, B float64
}

// https://stackoverflow.com/questions/22607043/color-gradient-algorithm, Mark Ransom
func blendColor(c1, c2 sdl.Color, frac float64) sdl.Color {
	gamma := 0.43
	c1Lin := color{R: fromsRGB(c1.R), G: fromsRGB(c1.G), B: fromsRGB(c1.B)}
	c2Lin := color{R: fromsRGB(c2.R), G: fromsRGB(c2.G), B: fromsRGB(c2.B)}
	c1Brightness := math.Pow(float64(c1Lin.R+c1Lin.G+c1Lin.B), gamma)
	c2Brightness := math.Pow(float64(c2Lin.R+c2Lin.G+c2Lin.B), gamma)
	intensity := lerp(c1Brightness, c2Brightness, frac)
	rc := color{R: lerp(c1Lin.R, c2Lin.R, frac), G: lerp(c1Lin.R, c2Lin.R, frac), B: lerp(c1Lin.R, c2Lin.R, frac)}
	sum := rc.R + rc.G + rc.B
	rc2 := color{R: rc.R * intensity / sum, G: rc.G * intensity / sum, B: rc.B * intensity / sum}
	ret := sdl.Color{R: tosRGB(rc2.R), G: tosRGB(rc2.G), B: tosRGB(rc2.B), A: 0}
	return ret

}

func tosRGBf(x float64) float64 {
	if x <= 0.0031308 {
		return 12.92 * x
	}
	return (1.055*(math.Pow(x, (1/2.4))) - 0.055)
}

func tosRGB(x float64) uint8 {
	return uint8(255.9999 * tosRGBf(x))
}

func fromsRGB(x uint8) float64 {
	var y float64
	xfloat := float64(x) / 255.0
	if xfloat <= 0.04045 {
		y = xfloat / 12.92
	} else {
		y = math.Pow(((xfloat + 0.055) / 1.055), 2.4)
	}
	return y
}

func lerp(a, b float64, frac float64) float64 {
	return a*(1-frac) + b*frac
}
