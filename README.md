# MIXAL Compiler

Μεταγλωττιστής που μεταφράζει custom γλώσσα σε MIXAL assembly.

## Χρήση

```bash
go run . examples/simple.txt
```

## Φάσεις Compilation

1. **Lexical Analysis**: Κείμενο → Tokens
2. **Syntax Analysis**: Tokens → AST
3. **Semantic Analysis**: Έλεγχος τύπων/μεταβλητών
4. **Code Generation**: AST → MIXAL

## Αρχεία

- `main.go`: Entry point
- `compiler.go`: Κύρια λογική
- `lexer.go`: Λεξικός αναλυτής
- `parser.go`: Συντακτικός αναλυτής
- `semantic.go`: Σημασιολογικός αναλυτής
- `codegen.go`: Παραγωγή MIXAL