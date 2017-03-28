package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	bg   *sdl.Texture
	bird *bird

	time int
}

func NewScene(r *sdl.Renderer) (*scene, error) {
	bckTexture, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("Could not load background texture: %v", err)
	}

	b, err := NewBird(r)
	if err != nil {
		return nil, err
	}
	return &scene{bg: bckTexture, bird: b}, nil
}

func (s *scene) Run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.HandleEvent(e); done {
					return
				}
			case <-tick:
				if err := s.Paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) HandleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.Jump()
	case *sdl.MouseMotionEvent:
	default:
		log.Printf("Unknown event %T", event)
	}
	return false
}

func (s *scene) Paint(r *sdl.Renderer) error {
	s.time++
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not background texture: %v", err)
	}

	if err := s.bird.Paint(r); err != nil {
		return err
	}

	r.Present()
	return nil
}

func (s *scene) Destroy() {
	s.bg.Destroy()
	s.bird.Destroy()
}
