package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
	score *Score
	level *Level
}

type Score struct {
	text    string
	myscore int
}

type Level struct {
	text    string
	mylevel int
}

func (l *Level) init() {
	l.text = "Level : "
	l.mylevel = 1
}

func (l *Level) update(score int) {
	l.mylevel = 1 + score/10
}

func (s *Score) init() {
	s.text = "Score : "
	s.myscore = 0
}

func NewScene(r *sdl.Renderer) (*scene, error) {
	var score Score
	score.init()

	var level Level
	level.init()

	bg, err := img.LoadTexture(r, "./res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load bg: %v", err)
	}

	var bird bird
	err = bird.loadTheBird(r)
	if err != nil {
		return nil, fmt.Errorf("could not load bird: %v", err)
	}

	var pipes pipes
	err = pipes.loadThePipes(r)
	if err != nil {
		return nil, fmt.Errorf("could not load pipes: %v", err)
	}

	return &scene{bg: bg, bird: &bird, pipes: &pipes, level: &level, score: &score}, nil
}

func (s *scene) run(r *sdl.Renderer, events <-chan sdl.Event) <-chan error {
	errch := make(chan error)

	go func() {
		defer close(errch)
		done := true
		tick := time.Tick(time.Millisecond * 10)
		for done {
			select {
			case e := <-events:
				// fmt.Println("event")
				done = s.handlerEvent(e)

			case <-tick:
				// fmt.Println("tick")
				err := s.paint(r)
				if err != nil {
					errch <- err
					return
				}
			case <-time.After(time.Second * 5):
				// fmt.Println("its working")
				return
			}
		}
	}()
	return errch
}

func (s *scene) handlerEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent, *sdl.MouseMotionEvent, *sdl.TextEditingEvent:
		go playClickSound()
		s.bird.speed_y = -5
		return false
	default:
		log.Printf("unknkown event %T", e)
		return false
	}
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()
	err := r.Copy(s.bg, nil, nil)
	if err != nil {
		return fmt.Errorf("could not copy backgrouond: %v", err)
	}

	err = s.bird.paintTheBird(r, s.pipes)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not copy bird: %v", err)
	}

	err = s.pipes.paintThePipes(r, s)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not copy pipes: %v", err)
	}

	s.addLevel(r)
	s.addScore(r)

	r.Present()
	return nil
}

func (s *scene) destroy() error {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
	return nil
}

func (s *scene) addLevel(r *sdl.Renderer) error {
	f, err := ttf.OpenFont("./res/fonts/GoFlap.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not init font: %v", err)
	}
	defer f.Close()

	textColor := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	surface, err := f.RenderUTF8Blended(s.level.text+strconv.Itoa(s.level.mylevel), textColor)
	if err != nil {
		return fmt.Errorf("Error rendering text:", err)

	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Error creating texture:", err)
	}
	defer texture.Destroy()

	textRect := sdl.Rect{X: 100, Y: 100, W: surface.W, H: surface.H}
	r.Copy(texture, nil, &textRect)
	return nil
}

func (s *scene) addScore(r *sdl.Renderer) error {
	f, err := ttf.OpenFont("./res/fonts/GoFlap.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not init font: %v", err)
	}
	defer f.Close()

	textColor := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	surface, err := f.RenderUTF8Blended(s.score.text+strconv.Itoa(s.score.myscore), textColor)
	if err != nil {
		return fmt.Errorf("Error rendering text:", err)

	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Error creating texture:", err)
	}
	defer texture.Destroy()

	textRect := sdl.Rect{X: 100, Y: 150, W: surface.W, H: surface.H}
	r.Copy(texture, nil, &textRect)
	return nil
}
