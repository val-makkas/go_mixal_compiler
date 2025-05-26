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

func printUsage() {
	fmt.Println("MIXAL Compiler - Compiles custom language to MIXAL assembly")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s <source_file>\n", os.Args[0])
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Printf("  %s examples/simple.txt\n", os.Args[0])
	fmt.Printf("  %s my_program.src\n", os.Args[0])
	fmt.Println()
	fmt.Println("Output:")
	fmt.Println("  Creates <source_file>.mixal with the generated assembly code")
}
