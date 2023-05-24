package main

import (
	"fmt"
	"os"
)

var opCodes = map[byte]string{
	0b00100010: "mov",
}

var registers = map[byte][]string{
	0b000: {"al", "ax"},
	0b001: {"cl", "cx"},
	0b010: {"dl", "dx"},
	0b011: {"bl", "bx"},
	0b100: {"ah", "sp"},
	0b101: {"ch", "bp"},
	0b110: {"dh", "si"},
	0b111: {"bh", "di"},
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("need a single file to process")
		return
	}
	decodeFile(os.Args[1])
}

func decode(buffer []byte) {
	opcode := (buffer[0] & 0b11111100) >> 2
	dFlag := (buffer[0] & 0b00000010) >> 1
	wFlag := (buffer[0] & 0b00000001)
	// mod := (buffer[1] & 0b11000000) >> 6
	reg := (buffer[1] & 0b00111000) >> 3
	rm := (buffer[1] & 0b00000111)

	if dFlag == 1 {
		fmt.Println(opCodes[opcode], registers[reg][wFlag], ",", registers[rm][wFlag])
	} else {
		fmt.Println(opCodes[opcode], registers[rm][wFlag], ",", registers[reg][wFlag])
	}

}

func decodeFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file: ", err)
		return
	}
	defer file.Close()
	fmt.Println("bits 16")

	buffer := make([]byte, 2)
	for {
		_, err := file.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("error reading file: ", err)
			return
		}
		if len(buffer) != 2 {
			break
		}
		decode(buffer)
	}
}
