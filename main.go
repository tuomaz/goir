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

	run()

	renderer.Destroy()
	window.Destroy()
}

func run() {
	var runFlag = true
	var event sdl.Event
	//var surface *sdl.Surface
	//var solid *sdl.Surface
	var clock *sdl.Surface
	var clockTexture *sdl.Texture
	var clock2 *sdl.Surface
	var clockTexture2 *sdl.Texture

	for runFlag {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				logger.Infof("key = %v", t.Keysym)
				runFlag = false
			}
		}
		var err error
		t := time.Now()

		if clock, err = font.RenderUTF8Blended(format(t), sdl.Color{R: 255, G: 255, B: 200, A: 200}); err != nil {
			logger.Fatalf("Failed to render text: %s\n", err)
		}

		if clock2, err = font.RenderUTF8Blended(t.Format(" 15.04"), sdl.Color{R: 255, G: 255, B: 200, A: 200}); err != nil {
			logger.Fatalf("Failed to render text: %s\n", err)
		}

		if clockTexture, err = renderer.CreateTextureFromSurface(clock); err != nil {
			logger.Fatalf("Failed to create texture from surface: %s\n", err)
		}

		if clockTexture2, err = renderer.CreateTextureFromSurface(clock2); err != nil {
			logger.Fatalf("Failed to create texture from surface: %s\n", err)
		}

		r3 := &sdl.Rect{
			H: clock.H,
			W: clock.W,
			X: 25,
			Y: 25,
		}

		r4 := &sdl.Rect{
			H: clock2.H,
			W: clock2.W,
			X: 25,
			Y: 160,
		}

		renderer.Clear()
		renderer.Copy(clockTexture, nil, r3)
		renderer.Copy(clockTexture2, nil, r4)
		renderer.Present()
		clock.Free()
		clockTexture.Destroy()
		sdl.Delay(100)

		/*
			if solid, err = font.RenderUTF8Solid("TEST", sdl.Color{255, 0, 0, 255}); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to render text: %s\n", err)
			}
			defer solid.Free()

			if surface, err = window.GetSurface(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get window surface: %s\n", err)
			}

			if err = solid.Blit(nil, surface, nil); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to put text on window surface: %s\n", err)
			}

			// Show the pixels for a while
			window.UpdateSurface()
		*/
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

	if font, err = ttf.OpenFont("fonts/Signika-Regular.ttf", 132); err != nil {
		logger.Fatalf("Failed to open font: %v\n", err)
	}

	sdl.ShowCursor(sdl.DISABLE)

}

func format(t time.Time) string {
	return fmt.Sprintf("%s %02d %s %d", days[t.Weekday()], t.Day(), months[t.Month()-1], t.Year())
}
