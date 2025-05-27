package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <name>")
		os.Exit(1)
	}

	sourceFile := os.Args[1]

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		fmt.Printf("Error: Source file '%s' does not exist.\n", sourceFile)
		os.Exit(1)
	}

	compiler := NewCompiler()

	fmt.Printf("Compiling source file: %s\n", sourceFile)
	fmt.Println("=====================================")

	err := compiler.Compile(sourceFile)
	if err != nil {
		fmt.Printf("Compilation failed: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Compilation successful!")
	fmt.Println("=====================================")
}
