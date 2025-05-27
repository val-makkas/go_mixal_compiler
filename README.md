# Go MIXAL Compiler

A compiler written in Go that translates a simple procedural programming language to MIXAL assembly code for Donald Knuth's MIX computer architecture.

### Running Tests

The `examples/` directory contains test cases:

```bash
# Test simple arithmetic
./mixal_compiler_go examples/simple.txt

# Test complex expressions
./mixal_compiler_go examples/complex.txt

# Test multiplication
./mixal_compiler_go examples/multiply.txt

# Test functions without parameters
./mixal_compiler_go examples/noparams.txt
```
