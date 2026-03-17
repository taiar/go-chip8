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
}

func (c *Chip8) LoadROM(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	copy(c.Memory[MemoryOffset:], data)

	return nil
}

func (c *Chip8) MemoryDump() {
	fmt.Println(hex.Dump(c.Memory[:]))
}
