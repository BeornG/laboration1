# interpreter

Kärnan i tolken

## interpreter filer

### lexer.go
delar upp text i tokens

TODO:
- Stöd för kommentarer
- Position/linjeinfo per token (för felmeddelanden)
- Stöd för attribut i noder, t.ex. (color=red)
- Fler nyckelord och diagramtyper



### parser.go
bygger up AST/datastruktur av tokens

### evaluator.go
transformerar AST till intern representation

### types.go
Token, Node, Edge, AST-strukturer

### interpreter.go
sammanordnar allt för paketet

