// main.go
package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	if err := run(); err != nil {
		// fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func playClickSound() {

	clickSound, err := mix.LoadWAV("./res/music/keyEvent.mp3") // Replace with the path to your click sound file
	if err != nil {
		fmt.Printf("Failed to load sound file: %v\n", err)
		return
	}
	defer clickSound.Free()

	clickSound.Play(-1, 1)

	sdl.Delay(250)
}

func playDeathSound() {

	clickSound, err := mix.LoadWAV("./res/music/death.mp3") // Replace with the path to your click sound file
	if err != nil {
		fmt.Printf("Failed to load sound file: %v\n", err)
		return
	}
	defer clickSound.Free()

	clickSound.Play(-1, 1)

	sdl.Delay(500)
}

func playGameOverSound() {

	clickSound, err := mix.LoadWAV("./res/music/gameOver.mp3") // Replace with the path to your click sound file
	if err != nil {
		fmt.Printf("Failed to load sound file: %v\n", err)
		return
	}
	defer clickSound.Free()

	clickSound.Play(-1, 1)

	sdl.Delay(4000)
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		return fmt.Errorf("could not init TTF: %v", err)
	}

	err = mix.OpenAudio(48000, sdl.AUDIO_S16, 2, 4096)
	if err != nil {
		return fmt.Errorf("failed to open audio: %v", err)
	}
	defer mix.CloseAudio()

	err = mix.Init(int(mix.MP3))
	if err != nil {
		return fmt.Errorf("failed to initialize SDL mixer: %v", err)
	}
	defer mix.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(1920, 1080, sdl.WINDOWEVENT_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()

	_ = r

	if err = drawTitle(r, "goflap"); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	s, err := NewScene(r)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.destroy()

	events := make(chan sdl.Event)
	runtime.LockOSThread()
	for {
		select {
		case err := <-s.run(r, events):
			return err
		case events <- sdl.WaitEvent():

		}
	}
}

func drawTitle(r *sdl.Renderer, title string) error {
	r.Clear()
	f, err := ttf.OpenFont("./res/fonts/GoFlap.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not init font: %v", err)
	}
	defer f.Close()

	s, err := f.RenderUTF8Solid(title, sdl.Color{
		R: 255,
		G: 100,
		B: 0,
		A: 255,
	})
	if err != nil {
		return fmt.Errorf("could not write title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture : %v", err)
	}
	defer t.Destroy()

	err = r.Copy(t, nil, nil)
	if err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	r.Present()
	time.Sleep(time.Second * 2)
	return nil
}
