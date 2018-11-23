package main

import (
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

	for runFlag {
		logger.Infof("1")
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				logger.Infof("key = %v", t.Keysym)
				runFlag = false
			}
		}
		var err error
		t := time.Now()
		logger.Infof("ts=%v", t.Format("15:04:05"))

		if clock, err = font.RenderUTF8Blended("TEST", sdl.Color{R: 255, G: 255, B: 200, A: 200}); err != nil {
			logger.Fatalf("Failed to render text: %s\n", err)
		}
		logger.Infof("2")
		if clockTexture, err = renderer.CreateTextureFromSurface(clock); err != nil {
			logger.Fatalf("Failed to create texture from surface: %s\n", err)
		}

		r3 := &sdl.Rect{
			H: clock.H,
			W: clock.W,
			X: 100,
			Y: 100,
		}

		renderer.Clear()
		renderer.Copy(clockTexture, nil, r3)
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
		logger.Infof("3")
	}
}

func initGraphics() {
	sdl.Init(sdl.INIT_EVERYTHING)
	var err error

	if err = ttf.Init(); err != nil {
		logger.Fatal("Failed to initialize TTF: %s\n", err)
	}

	if window, err = sdl.CreateWindow(name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_RESIZABLE); err != nil {
		logger.Fatalf("Failed to create window: %s\n", err)
	}

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE); err != nil {
		logger.Fatalf("Failed to create renderer: %v\n", err)
	}

	if font, err = ttf.OpenFont("fonts/Signika-Regular.ttf", 144); err != nil {
		logger.Fatalf("Failed to open font: %v\n", err)
	}

	sdl.ShowCursor(sdl.DISABLE)

}
