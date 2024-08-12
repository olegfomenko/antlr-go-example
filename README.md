# ANTLR Go Example

Example of calculator program using ANTLR 4 Visitor on Go

To re-build grammar use:
```shell
antlr -visitor -no-listener -Dlanguage=Go -o internal/antlr Expr.g4
```