package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
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

	p, err := NewPipes(r)
	if err != nil {
		return nil, err
	}

	return &scene{bg: bckTexture, bird: b, pipes: p}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()

				if s.bird.isDead() {
					drawTitle(r, "Game Over")
					time.Sleep(3 * time.Second)
					s.restart()
				}

				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.jump()
	case *sdl.MouseMotionEvent:
	default:
		log.Printf("Unknown event %T", event)
	}
	return false
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not background texture: %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return err
	}

	if err := s.pipes.paint(r); err != nil {
		return err
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}
