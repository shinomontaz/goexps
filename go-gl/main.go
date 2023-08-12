package main

import "github.com/go-gl/gl/v3.3-core/gl"

func main() {
	window := ... // Open a window.
	window.MakeContextCurrent()

	// Important! Call gl.Init only under the presence of an active OpenGL context,
	// i.e., after MakeContextCurrent.
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
}