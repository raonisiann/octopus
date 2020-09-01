package parser

import (
	"fmt"
	"octopus/lexer"
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
	tk := lexer.GetToken()
	fmt.Printf(
		"Expected '%s', but get '%s'\n",
		lexer.GetTokenText(expected.Class),
		lexer.GetTokenText(tk.Class),
	)
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
		lexer.NextToken()
		tk := lexer.GetToken()

		if tk.Class == lexer.TkEOF {
			fmt.Println("End of file")
			break
		}

		fmt.Printf("%s(%s)\n", lexer.GetTokenText(tk.Class), tk.Value)
	}
}
