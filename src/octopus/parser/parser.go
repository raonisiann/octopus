package parser

import (
	"fmt"
	"octopus/lexer"
	"os"
)

var identLevel int = 0
var identSize int = 4

// Parse is the main parser function
func Parse(fileName string) {

	lexer.Init(fileName)

	identLevel = 0

	topLevel()
}

func error(e lexer.TkClassType) {
	tk := lexer.GetToken()
	fmt.Printf(
		"Expected '%s', but get '%s' on position %d line %d\n",
		lexer.GetTokenText(e),
		lexer.GetTokenText(tk.Class),
		lexer.GetTokenCurrentPosition(),
		lexer.GetTokenCurrentLine(),
	)
	os.Exit(-1)
}

func accept(c lexer.TkClassType) bool {
	if c == lexer.GetToken().Class {
		return true
	}
	return false
}

func expect(c lexer.TkClassType) bool {
	if accept(c) {
		return true
	}
	error(c)
	return false
}

func getIdentLevel() int {
	fmt.Printf("-->> Ident: %d", identLevel)
	return identLevel
}

func ident() {
	identLevel += identSize
}

func dedent() {
	identLevel -= identSize
}

func topLevel() {

	for {
		lexer.NextToken()
		tk := lexer.GetToken()

		switch tk := lexer.GetToken(); tk.Class {
		case lexer.TkEOF:
			fmt.Println("End of file")
			os.Exit(-1)
		case lexer.TkNewLine:
			fmt.Println("NEW_LINE")
			continue
		case lexer.TkClassDef:
			classHeader()
		default:
			fmt.Printf("Unexpected token %s at top level", lexer.GetTokenText(tk.Class))
		}

		fmt.Printf("%s => %s\n", lexer.GetTokenText(tk.Class), tk.Value)
	}
}

func classHeader() {
	expect(lexer.TkClassDef)

	lexer.NextToken()
	tkClassIdentifier := lexer.GetToken()
	fmt.Printf(" class name = %s\n", tkClassIdentifier.Value)

	lexer.NextToken()
	expect(lexer.TkColon)
	lexer.NextToken()
	expect(lexer.TkNewLine)
}

func classBody() {

}
