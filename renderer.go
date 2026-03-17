package main

import "fmt"

type Renderer interface {
	Draw(grid [64 * 32]byte)
	Clear()
}

type TerminalRenderer struct{}

func (t *TerminalRenderer) Clear() {
	// Comando ANSI para limpar o terminal e resetar o cursor
	fmt.Print("\033[H\033[2J")
}

func (t *TerminalRenderer) Draw(grid [64 * 32]byte) {
	t.Clear() // Limpa antes de desenhar o novo frame
	var output string

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			pixel := grid[y*64+x]
			if pixel == 1 {
				output += "██" // Pixel aceso
			} else {
				output += "  " // Pixel apagado
			}
		}
		output += "\n"
	}
	fmt.Print(output)
}
