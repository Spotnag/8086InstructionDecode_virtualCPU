package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
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

var (
	intputFileFlag     = flag.String("input", "", "binary file to read")
	comparisonFileFlag = flag.String("compareFile", "", "ASM file to compare with")
)

func main() {
	flag.Parse()
	decodedFile, err := decodeFile()
	if err != nil {
		fmt.Println("error decoding file", err)
	}
	if *comparisonFileFlag != "" {
		compareFiles(decodedFile)
	}
	fmt.Print(decodedFile)
}

func compareFiles(decodedFile string) {
	comparisonFile, err := os.Open(*comparisonFileFlag)
	if err != nil {
		fmt.Printf("Failed to open comparison file: %v", err)
	}

	reader1 := bufio.NewReader(strings.NewReader(decodedFile))
	reader2 := bufio.NewReader(comparisonFile)

	for {
		line1, err1 := reader1.ReadString('\n')
		line2, err2 := reader2.ReadString('\n')

		if line1 != line2 {
			fmt.Printf("lines do not match. DecodedFile Line: \n %v \n ComparisonFile Line: \n %v \n", line1, line2)
		}

		if err1 != nil {
			if err1.Error() != "EOF" {
				fmt.Printf("error while reading f1: %v", err1)
			}
		}

		if err2 != nil {
			if err2.Error() != "EOF" {
				fmt.Printf("error while reading file2: %v", err2)
			}
		}

		if err1 == io.EOF && err2 == io.EOF {
			break
		}
	}

}

func decode(buffer []byte) string {
	opcode := (buffer[0] & 0b11111100) >> 2
	dFlag := (buffer[0] & 0b00000010) >> 1
	wFlag := (buffer[0] & 0b00000001)
	// mod := (buffer[1] & 0b11000000) >> 6
	reg := (buffer[1] & 0b00111000) >> 3
	rm := (buffer[1] & 0b00000111)

	if dFlag == 1 {
		return fmt.Sprint(opCodes[opcode], " ", registers[reg][wFlag]+", ", registers[rm][wFlag])
	} else {
		return fmt.Sprint(opCodes[opcode], " ", registers[rm][wFlag]+", ", registers[reg][wFlag])
	}

}

func decodeFile() (string, error) {
	file, err := os.Open(*intputFileFlag)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var sb strings.Builder
	sb.WriteString("bits 16")
	sb.WriteString("\n")

	buffer := make([]byte, 2)
	for {
		_, err := file.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return "", fmt.Errorf("error reading file: %v", err)
		}
		if len(buffer) != 2 {
			break
		}
		decodedLine := decode(buffer)
		sb.WriteString(decodedLine)
		sb.WriteString("\n")
	}
	return sb.String(), nil
}
