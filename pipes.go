package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipes struct {
	time     int
	speed_x  float32
	texture  *sdl.Texture
	y        int32
	allPipes []*pipe
}

type pipe struct {
	w        int32
	x        int32
	h        int32
	inverted bool
}

func (ps *pipes) loadThePipes(r *sdl.Renderer) error {
	ps.init()
	t, err := img.LoadTexture(r, "./res/imgs/pipe.png")
	if err != nil {
		return fmt.Errorf("could not load bg: %v", err)
	}
	ps.texture = t

	go func() {
		for {
			ps.allPipes = append(ps.allPipes, newPipe())
			time.Sleep(time.Second)
		}
	}()

	return nil
}

func (ps *pipes) paintThePipes(r *sdl.Renderer, s *scene) error {
	if success := ps.update(s); !success {
		return fmt.Errorf("cannot update further, you died")
	}

	for _, p := range ps.allPipes {
		var flip sdl.RendererFlip
		rect := &sdl.Rect{
			X: p.x,
			Y: 0,
			W: p.w,
			H: p.h,
		}
		if p.inverted {
			rect.Y = 1080 - p.h
		} else {
			flip = sdl.FLIP_VERTICAL
		}
		err := r.CopyEx(ps.texture, nil, rect, 0, nil, flip)
		if err != nil {
			return fmt.Errorf("could not copy the pipes on texture: %v", err)
		}
	}

	return nil
}

func (ps *pipes) update(s *scene) bool {
	ps.time++
	for _, p := range ps.allPipes {
		p.x -= int32(ps.speed_x)
	}
	var newpipes []*pipe
	for _, p := range ps.allPipes {
		if p.x+p.w > 0 {
			newpipes = append(newpipes, p)
		} else {
			s.score.myscore += 1
			s.level.update(s.score.myscore)
		}
	}
	ps.allPipes = newpipes
	return true
}

func (ps *pipes) init() {
	ps.time = 0
	ps.speed_x = 1
	ps.y = 320
}

func newPipe() *pipe {
	p := pipe{}
	p.h = int32(rand.Intn(400)) + 100
	p.inverted = rand.Float32() > .5
	p.w = 52
	p.x = 1920
	return &p
}

func (ps *pipes) destroy() {
	ps.texture.Destroy()
}
