package parser

import (
	"fmt"
	"octopus/lexer"
	"os"
)

// Parse is the main parser function
func Parse(fileName string) {

	lexer.Init(fileName)

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

func ignoreEmptyNewLines() {

	for {
		lexer.NextToken()
		if lexer.GetToken().Class != lexer.TkNewLine {
			break
		}
	}
}

func topLevel() {

	for {
		lexer.NextToken()
		tk := lexer.GetToken()

		switch tk := lexer.GetToken(); tk.Class {
		case lexer.TkEOF:
			return
		case lexer.TkNewLine:
			fmt.Println("NEW_LINE")
			continue
		case lexer.TkClassDef:
			class(0)
		default:
			fmt.Printf("Unexpected token %s at top level", lexer.GetTokenText(tk.Class))
		}

		fmt.Printf("%s => %s\n", lexer.GetTokenText(tk.Class), tk.Value)
	}
}

func class(ident int) {
	// class header
	expect(lexer.TkClassDef)

	lexer.NextToken()
	tkClassIdentifier := lexer.GetToken()
	fmt.Printf(" class name = %s\n", tkClassIdentifier.Value)

	lexer.NextToken()
	expect(lexer.TkColon)
	lexer.NextToken()
	expect(lexer.TkNewLine)

	definitions()
	os.Exit(-1)
}

func definitions() {
	// class body
	expectedIdent := lexer.GetIdentLevel() + 1
	ignoreEmptyNewLines()
	for lexer.GetIdentLevel() == expectedIdent {

		switch tk := lexer.GetToken(); tk.Class {
		case lexer.TkIdentifier:
			fmt.Println("IDENTIFIER")
			identifier()
		case lexer.TkResourceFile:
			fmt.Println("RESOURCE")
		case lexer.TkResourcePackage:
			fmt.Println("PACKAGE")
		case lexer.TkResourceService:
			fmt.Println("SERVICE")
		default:
			fmt.Printf("Unexpected token %s at this point\n", lexer.GetTokenText(tk.Class))
			os.Exit(-1)
		}

		lexer.NextToken()
	}

	fmt.Printf("IDENTITY_LEVEL_HERE=%d\n", lexer.GetIdentLevel())
}

func identifier() {
	expect(lexer.TkIdentifier)
	lexer.NextToken()
	expect(lexer.TkEqual)
	lexer.NextToken()
	expression()
}

func expression() {

	for {
		switch tk := lexer.GetToken(); tk.Class {
		case lexer.TkString:
			fmt.Println("STRING")
		case lexer.TkInt:
			fmt.Println("INT")
		default:
			fmt.Printf("Unexpected token %s at this point", lexer.GetTokenText(tk.Class))
			os.Exit(-1)
		}

		lexer.NextToken()
	}

}

func resourceFile(ident int) {

}

func resourcePackage(ident int) {

}

func resourceService(ident int) {

}
