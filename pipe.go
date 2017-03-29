package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type pipes struct {
	mu      sync.RWMutex
	texture *sdl.Texture
	pipes   []*pipe
}

func NewPipes(r *sdl.Renderer) (*pipes, error) {
	texture, err := img.LoadTexture(r, "res/imgs/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("Could not load pipe texture: %v", err)
	}

	ps := &pipes{texture: texture}

	go func() {
		tick := time.Tick(2 * time.Second)
		for {
			select {
			case <-tick:
				ps.mu.Lock()
				ps.pipes = append(ps.pipes, NewPipe())
				ps.mu.Unlock()
			}
		}
	}()

	return ps, nil
}

func (ps *pipes) update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var rem []*pipe
	for _, p := range ps.pipes {
		p.update()

		if p.x+p.w > 0 {
			rem = append(rem, p)
		}
	}

	ps.pipes = rem
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(r, ps.texture); err != nil {
			return err
		}
	}

	return nil
}

func (ps *pipes) touch(b *bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		p.touch(b)
	}
}

func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *pipes) destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.texture.Destroy()
}

type pipe struct {
	mu       sync.RWMutex
	x        int32
	h        int32
	w        int32
	speed    int32
	inverted bool
}

func NewPipe() *pipe {
	return &pipe{
		x:        800,
		h:        100 + int32(rand.Intn(300)),
		w:        50,
		speed:    1,
		inverted: rand.Float32() > 0.5,
	}
}

func (p *pipe) update() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.x -= p.speed
}

func (p *pipe) touch(b *bird) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	b.touch(p)
}

func (p *pipe) paint(r *sdl.Renderer, t *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	flip := sdl.FLIP_NONE
	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := r.CopyEx(t, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("could not copy pipe texture: %v", err)
	}

	return nil
}
