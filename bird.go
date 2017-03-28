package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

const (
	gravity   = 0.2
	jumpSpeed = 5
)

type bird struct {
	textures []*sdl.Texture
	time     int

	y, speed float64
}

func NewBird(r *sdl.Renderer) (*bird, error) {
	b := &bird{y: 300, speed: 0}

	for i := 1; i < 5; i++ {
		path := fmt.Sprintf("res/imgs/frame-%d.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("Could not load bird texture: %v", err)
		}
		b.textures = append(b.textures, texture)
	}

	return b, nil
}

func (b *bird) Paint(r *sdl.Renderer) error {
	b.time++
	b.y -= b.speed
	if b.y < 0 {
		b.y = 0
		b.speed = -b.speed
	}
	b.speed += gravity

	rect := &sdl.Rect{X: 10, Y: (600 - int32(b.y)) - 43/2, W: 50, H: 43}

	i := b.time / 10 % len(b.textures)
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird texture: %v", err)
	}

	return nil
}

func (b *bird) Jump() {
	b.speed = -jumpSpeed
}

func (b *bird) Destroy() {
	for _, t := range b.textures {
		t.Destroy()
	}
}
