package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

const MemoryOffset = 0x200

type Chip8 struct {
	Memory     [4096]byte // 4Kb de Ram
	V          [16]byte   // Registradores V0 a VF
	I          uint16     // Registrador de índice
	PC         uint16     // Program Counter
	Stack      [16]uint16 // Pilha para sub-rotinas
	SP         uint16     // Stack Pointer
	DelayTimer byte
	SoundTimer byte
	Display    [64 * 32]byte // 0 para apagado, 1 para aceso
	VideoOut   Renderer
}

func (c *Chip8) Init(videoDisplay Renderer) {
	c.PC = MemoryOffset
	c.VideoOut = videoDisplay
}

func (c *Chip8) Run() {
	ticker := time.NewTicker(time.Second / 500) // 500Hz
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Cycle()
		}
	}
}

func (c *Chip8) LoadROM(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	copy(c.Memory[MemoryOffset:], data)

	return nil
}

func (c *Chip8) Cycle() {
	// 1. FETCH: Busca o opcode de 2 bytes
	// Pegamos o byte na posição PC e deslocamos 8 bits para a esquerda
	// Depois fazemos um OR (|) com o byte na posição PC + 1

	highByte := uint16(c.Memory[c.PC])
	lowByte := uint16(c.Memory[c.PC+1])

	opcode := (highByte << 8) | lowByte

	fmt.Printf("PC: 0x%03X | Opcode: 0x%04X\n", c.PC, opcode)

	// 2. DECODE & EXECUTE: (Vamos preparar o terreno aqui)
	c.Execute(opcode)

	// 3. INCREMENTAR PC:
	// Por padrão, cada instrução tem 2 bytes.
	// Algumas instruções de pulo alteram o PC diretamente,
	// mas o padrão é avançar para a próxima.
	c.PC += 2
}

func (c *Chip8) Execute(opcode uint16) {
	// Extraímos os componentes comuns para facilitar a leitura
	x := (opcode & 0x0F00) >> 8 // O segundo nibble
	y := (opcode & 0x00F0) >> 4 // O terceiro nibble
	n := opcode & 0x000F        // O último nibble
	// nnn := opcode & 0x0FFF     // Os últimos 12 bits (endereço)
	nn := byte(opcode & 0x00FF) // Os últimos 8 bits (valor imediato)

	switch opcode & 0xF000 {
	case 0x0000:
		// Limpa o array do display colocando 0 em tudo
		for i := range c.Display {
			c.Display[i] = 0
		}
		c.VideoOut.Clear()

	case 0x1000:
		// 1NNN: Jump para o endereço NNN
		address := opcode & 0x0FFF
		c.PC = address
		// Importante: Como mudamos o PC manualmente,
		// temos que evitar o PC += 2 no final do Cycle!
		c.PC -= 2

	case 0xA000:
		// ANNN: Seta o registrador I para NNN
		address := opcode & 0x0FFF
		c.I = address
		fmt.Printf("Instrução: I = 0x%X\n", address)

	case 0x6000:
		// 0x6XNN (Set Register VX)
		c.V[x] = nn

	case 0xF000:
		// 0xFX29 (Set I to Sprite Location)
		switch nn {
		case 0x29:
			// c.V[x] contém o caractere (ex: 0x0B).
			// Multiplicamos por 5 porque cada sprite tem 5 bytes de altura.
			c.I = uint16(0x50) + (uint16(c.V[x]) * 5)
		}

	case 0xD000:
		// 0xDXYN (Draw Sprite)
		xCoord := c.V[x] % 64 // Wrap around (opcional, depende da ROM)
		yCoord := c.V[y] % 32
		c.V[0xF] = 0 // Reseta o flag de colisão

		for row := uint16(0); row < n; row++ {
			spriteByte := c.Memory[c.I+row]
			for col := uint16(0); col < 8; col++ {
				// Verifica se o bit específico do spriteByte está ligado (1)
				// Começamos do bit mais significativo (0x80)
				if (spriteByte & (0x80 >> col)) != 0 {
					// Cálculo do índice no array linear da tela
					idx := (uint16(xCoord) + col) + ((uint16(yCoord) + row) * 64)

					// Se o índice estourar a tela, paramos ou fazemos wrap
					if idx < 64*32 {
						if c.Display[idx] == 1 {
							c.V[0xF] = 1 // Colisão detectada!
						}
						c.Display[idx] ^= 1 // XOR: Liga se estava desligado, desliga se estava ligado
					}
				}
			}
		}
		c.VideoOut.Draw(c.Display)

	default:
		fmt.Printf("Opcode não implementado: 0x%X\n", opcode)
	}
}

func (c *Chip8) MemoryDump() {
	fmt.Println(hex.Dump(c.Memory[:]))
}
