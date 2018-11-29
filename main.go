package main

import (
	"fmt"
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
)

func init() {
	fullLogger, _ = zap.NewProduction()
	defer fullLogger.Sync() // flushes buffer, if any
	logger = fullLogger.Sugar()
}

func main() {
	logger.Infof("Starting goir...")
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
				runFlag = false
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

		sdl.Delay(5000)
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

func getColor() sdl.Color {
	return sdl.Color{R: 255, G: 255, B: 200, A: 200}
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

	return items
}
