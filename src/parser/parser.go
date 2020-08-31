package parser

import (
	"fmt"
	"lexer"
	"os"
)

var identLevel int
var identSize int = 4

// Parse is the main parser function
func Parse(fileName string) {

	lexer.Init(fileName)
	identLevel = 0

	topLevel()
}

func error(expected lexer.Token) {
	fmt.Printf("Expected '%s', but get '%s'\n", expected.Class, lexer.GetTokenText())
	os.Exit(-1)
}

func accept(token lexer.Token) bool {
	if token == lexer.GetToken() {
		return true
	}
	return false
}

func expect(token lexer.Token) bool {
	if accept(token) {
		return true
	}
	error(token)
	return false
}

func topLevel() {

	for {
		switch tk := lexer.GetToken(); tk.Class {
		case lexer.TkEOF:
			fmt.Println("End of file")
			return
		case lexer.TkNewLine:
			fmt.Printf("")
		case lexer.TkClassDef:
			fmt.Printf("Class")
		}

	}
}
