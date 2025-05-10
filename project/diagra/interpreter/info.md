# interpreter

Kärnan i tolken

TODO:
- Stöd för kommentarer
- Position/linjeinfo per token (för felmeddelanden)
- Fler nyckelord och diagramtyper

## interpreter filer

### lexer.go
delar upp text i tokens

### parser.go
bygger up AST/datastruktur av tokens

### evaluator.go
transformerar AST till intern representation

### types.go
Token, Node, Edge, AST-strukturer

### interpreter.go
sammanordnar allt för paketet

