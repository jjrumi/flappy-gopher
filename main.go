package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("Could not initialise SDL: %v", err)
	}
	defer sdl.Quit()

	if err = ttf.Init(); err != nil {
		return fmt.Errorf("Could not initalise TTF: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("Could not create Window: %v", err)
	}
	defer w.Destroy() // Free both window and renderer.

	if err = drawTitle(r, "Happy Gopher"); err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	scene, err := NewScene(r)
	if err != nil {
		return fmt.Errorf("Could not create scene: %v", err)
	}
	defer scene.destroy()

	events := make(chan sdl.Event)
	errc := scene.run(events, r)
	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}

	return nil
}

func drawTitle(r *sdl.Renderer, text string) error {
	r.Clear()

	f, err := ttf.OpenFont("res/fonts/Flappy.ttf", 20)
	if err != nil {
		return fmt.Errorf("Could not open font: %v", err)
	}
	defer f.Close()

	s, err := f.RenderUTF8_Solid(text, sdl.Color{R: 255, G: 125, B: 0})
	if err != nil {
		return fmt.Errorf("Could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("Could not create texture: %v", err)
	}
	defer t.Destroy()

	if err = r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	r.Present()

	return nil
}
