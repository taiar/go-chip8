package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

func main() {

	type Chip8 struct {
		Memory 		[4096]byte 	// 4Kb de Ram
		V			[16]byte 	 	// Registradores V0 a VF
		I			uint16			// Registrador de índice
		PC			uint16			// Program Counter
		Stack		[16]uint16	// Pilha para sub-rotinas
		SP			uint16			// Stack Pointer
		DelayTimer	byte
		SoundTimer	byte
	}

	func (c *Chip8) LoadROM(filename string) error {
		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		memoryOffset := 0x200

		for i := 0; i < len(data); i++ {
			c.Memory[memoryOffset + 1] = data[i]
		}

		return nil
	}

	// bytes, err := os.ReadFile("./roms/test.ch8")
	// if err != nil {
	// 	panic(err)
	// }

	// var memory [4096]byte
	// memoryOffset := 512

	// // lê programa em memória
	// for i := 0; i < len(bytes); i++ {
	// 	memory[i+memoryOffset] = bytes[i]
	// }

	fmt.Println(string(bytes))             // Unicode string
	fmt.Println(hex.EncodeToString(bytes)) // Raw hex bytes
	fmt.Println(hex.Dump(bytes))           // Nice bytes
	fmt.Println("------------------------")
	fmt.Println(hex.Dump(memory[:]))
}

