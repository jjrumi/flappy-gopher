package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

const (
	gravity   = 0.05
	jumpSpeed = 4
)

type bird struct {
	mu       sync.RWMutex
	textures []*sdl.Texture
	time     int

	x, y int32
	w, h int32

	speed float64
	dead  bool
}

func NewBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i < 5; i++ {
		path := fmt.Sprintf("res/imgs/frame-%d.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("Could not load bird texture: %v", err)
		}
		textures = append(textures, texture)
	}

	return &bird{textures: textures, x: 10, y: 300, w: 50, h: 43, speed: 0}, nil
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.time++
	b.y -= int32(b.speed)
	if b.y < 0 {
		b.dead = true
	}
	b.speed += gravity
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	rect := &sdl.Rect{X: b.x, Y: (600 - b.y) - b.h/2, W: b.w, H: b.h}

	i := b.time / 10 % len(b.textures)
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird texture: %v", err)
	}

	return nil
}

func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.dead
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.speed = -jumpSpeed
}

func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if p.x > b.x+b.w { // Pipe is too far right.
		return
	}
	if p.x+p.w < b.x { // Pipe is too far left.
		return
	}
	if !p.inverted && p.h < b.y-b.h/2 { // Pipe is too low.
		return
	}

	if p.inverted && p.h > b.y-b.h/2 { // Pipe is too low.
		return
	}
	b.dead = true
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, t := range b.textures {
		t.Destroy()
	}
}
