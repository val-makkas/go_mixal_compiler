package main

type CodeGenerator struct{}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (c *CodeGenerator) Generate(ast *AST, symbolTables map[string]interface{}) (string, error) {
	return "mixal code", nil
}
