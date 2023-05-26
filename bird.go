package main

import (
	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type bird struct {
	life     *Life
	time     int
	birds    []*sdl.Texture
	x, y     int32
	gravit_y float32
	speed_y  float32
	dead     bool
}

type Life struct {
	text   string
	mylife int
}

func (l *Life) init() {
	l.text = "Life : "
	l.mylife = 5
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

		if b.life.mylife > 1 {
			go playDeathSound()
			b.paintKilled(r)
			life := b.life.mylife - 1
			pipes.speed_x *= 1.1
			b.init()
			b.life.mylife = life
		} else {
			go playGameOverSound()
			drawTitle(r, "Game Over")
			return fmt.Errorf("cannot update further, you died")
		}
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
	b.paintLife(r)
	return nil
}

func (b *bird) paintLife(r *sdl.Renderer) error {
	f, err := ttf.OpenFont("./res/fonts/GoFlap.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not init font: %v", err)
	}
	defer f.Close()

	textColor := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	surface, err := f.RenderUTF8Blended(b.life.text+strconv.Itoa(b.life.mylife), textColor)
	if err != nil {
		return fmt.Errorf("Error rendering text:", err)

	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Error creating texture:", err)
	}
	defer texture.Destroy()

	textRect := sdl.Rect{X: 100, Y: 200, W: surface.W, H: surface.H}
	r.Copy(texture, nil, &textRect)
	return nil
}

func (b *bird) paintKilled(r *sdl.Renderer) error {
	f, err := ttf.OpenFont("./res/fonts/GoFlap.ttf", 25)
	if err != nil {
		return fmt.Errorf("could not init font: %v", err)
	}
	defer f.Close()

	textColor := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	surface, err := f.RenderUTF8Blended("KILLED", textColor)
	if err != nil {
		return fmt.Errorf("Error rendering text:", err)

	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Error creating texture:", err)
	}
	defer texture.Destroy()

	textRect := sdl.Rect{X: 200, Y: 200, W: surface.W, H: surface.H}
	r.Copy(texture, nil, &textRect)

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
	b.life = &Life{}
	b.life.init()
}

func (b *bird) destroy() {
	for _, t := range b.birds {
		t.Destroy()
	}
}
