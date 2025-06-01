# MIXAL Compiler

A compiler written in Go that translates a simple procedural programming language to MIXAL assembly code for Donald Knuth's MIX computer architecture.

### Build

```bash
#simply go build . in source folder
go build .
```

### Running Tests

The `examples/` directory contains test cases both successful and ment to produce error:

```bash
# successful tests
./mixal_compiler examples/success/ 0 | 1 | 2 | 3 .txt

# error induncing tests
./mixal_compiler examples/error/ 0 | 1 .txt
```