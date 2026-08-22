Now the terminal can read the byes , handle ansi chars too

# Next Step:
2. make cursor top and down arrow work (history)
3. Make the defer raw mode and stuff work properly (working now but need testing)
4. When Text is wraped the backspace is not working
5. Implement Fuzz test
6. Fix renderer.go

# NOTE:
- before moving to the parser and AST, lets make sure that this shell is compatable with the Oh my posh and Starship

- Mare sure pipe and & , and >> , > kinda stuff works before AST and parsing

