package main

import (
	"encoding/hex"
	"fmt"
	"os"
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
}

func (c *Chip8) Init() {
	c.PC = MemoryOffset
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
	// Extrai o primeiro dígito hexadecimal (ex: de 0xA2F0, pega o 'A')
	firstNibble := opcode & 0xF000

	switch firstNibble {
	case 0x0000:
		// Limpa o array do display colocando 0 em tudo
		for i := range c.Display {
			c.Display[i] = 0
		}
		fmt.Println("Tela limpa!")

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

	default:
		fmt.Printf("Opcode não implementado: 0x%X\n", opcode)
	}
}

func (c *Chip8) MemoryDump() {
	fmt.Println(hex.Dump(c.Memory[:]))
}
