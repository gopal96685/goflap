package main

import (
	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type bird struct {
	time     int
	birds    []*sdl.Texture
	x, y     int32
	gravit_y float32
	speed_y  float32
	dead     bool
}

func (b *bird) loadTheBird(r *sdl.Renderer) error {
	b.init()
	for i := 1; i <= 4; i++ {
		path := "./res/imgs/bird_frame_" + strconv.Itoa(i) + ".png"
		birdCurrent, err := img.LoadTexture(r, path)
		if err != nil {
			return fmt.Errorf("could not load bg: %v", err)
		}
		b.birds = append(b.birds, birdCurrent)
	}
	return nil
}

func (b *bird) paintTheBird(r *sdl.Renderer, pipes *pipes) error {
	if success := b.update(pipes); !success {
		// return fmt.Errorf("cannot update further, you died")
		drawTitle(r, "Game over")
		pipes.speed_x *= 2
		b.init()
	}
	rect := &sdl.Rect{
		X: b.x,
		Y: 1080 - b.y,
		W: 50,
		H: 43,
	}

	index := b.time / 10 % len(b.birds)
	err := r.Copy(b.birds[index], nil, rect)
	if err != nil {
		return fmt.Errorf("could not copy the bird on texture: %v", err)
	}
	return nil
}

func (b *bird) update(pipes *pipes) bool {
	b.time++
	b.speed_y += b.gravit_y
	b.y -= int32(b.speed_y)
	if b.y < 0 || b.y > 1080 {
		b.die()
		return false
	}
	if b.detectCollision(pipes) {
		return false
	}
	return true
}

func (b *bird) detectCollision(pipes *pipes) bool {
	for _, pipe := range pipes.allPipes {
		if pipe.inverted {
			if b.x >= pipe.x && b.x <= pipe.x+pipe.w {
				if b.y <= pipe.h {
					return true
				}
			}
		} else {
			if b.x >= pipe.x && b.x <= pipe.x+pipe.w {
				if b.y >= 1080-pipe.h {
					return true
				}
			}
		}

	}
	return false
}

func (b *bird) die() {
	b.dead = true
	b.speed_y = 0
}

func (b *bird) init() {
	b.y = 1080/2 - 43/2
	b.x = 10
	b.speed_y = 0
	b.gravit_y = .1
	b.dead = false
	b.time = 0
}

func (b *bird) destroy() {
	for _, t := range b.birds {
		t.Destroy()
	}
}
