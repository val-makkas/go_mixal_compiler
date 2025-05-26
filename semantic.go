package main

type SemanticAnalyzer struct{}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{}
}

func (s *SemanticAnalyzer) Analyze(ast *AST) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}
