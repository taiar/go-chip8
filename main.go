package main

func main() {
	var chip Chip8

	chip.LoadROM("./roms/test.ch8")
	chip.MemoryDump()

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

	// fmt.Println(string(bytes))             // Unicode string
	// fmt.Println(hex.EncodeToString(bytes)) // Raw hex bytes
	// fmt.Println(hex.Dump(bytes))           // Nice bytes
	// fmt.Println("------------------------")
	// fmt.Println(hex.Dump(memory[:]))
}
