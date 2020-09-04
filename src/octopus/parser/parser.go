package parser

import (
	"fmt"
	"octopus/lexer"
	"os"
)

const firstIdentLevel int = 0

// Parse is the main parser function
func Parse(fileName string) {

	lexer.Init(fileName)
	lexer.NextToken()
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
	tk := lexer.GetToken()

	if c == tk.Class {
		fmt.Printf("Accepted '%s'\n", lexer.GetTokenText(tk.Class))
		lexer.NextToken()
		return true
	}
	return false
}

// acceptAny will check if the current token class
// matches any of the token classes in the list provided
// as argument.
func acceptAny(classes ...lexer.TkClassType) bool {

	for _, v := range classes {
		if accept(v) {
			return true
		}
	}
	return false
}

// acceptLookAhead can lookup and match 2 tokens ahead.
// If lookahead fails, we rollback to the previous token
func acceptLookAhead(c lexer.TkClassType, lookahead lexer.TkClassType) bool {
	if accept(c) {
		if accept(lookahead) {
			return true
		}
		lexer.PopToken()
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

func expectOne(list ...lexer.TkClassType) bool {

	strClass := ""

	for _, v := range list {
		if accept(v) {
			return true
		}
		strClass = strClass + "'" + lexer.GetTokenText(v) + "', "
	}

	fmt.Printf(
		"Expected at least one of %s but get '%s'\n",
		strClass,
		lexer.GetTokenText(lexer.GetToken().Class),
	)

	return false
}

func ignoreEmptyNewLines() {

	for accept(lexer.TkNewLine) {
	}
}

func topLevel() {

	for {
		tk := lexer.GetToken()

		switch tk := lexer.GetToken(); tk.Class {
		case lexer.TkEOF:
			return
		case lexer.TkNewLine:
			expect(lexer.TkNewLine)
			continue
		case lexer.TkClassDef:
			expect(lexer.TkClassDef)
			class(firstIdentLevel)
			continue
		default:
			fmt.Printf("Unexpected token '%s' at top level\n", lexer.GetTokenText(tk.Class))
			os.Exit(-1)
		}

		fmt.Printf("%s => %s\n", lexer.GetTokenText(tk.Class), tk.Value)
		lexer.NextToken()
	}
}

func class(expectedIdent int) {

	tkClassIdentifier := lexer.GetToken()
	fmt.Printf(" class name = %s\n", tkClassIdentifier.Value)
	expect(lexer.TkIdentifier)
	expect(lexer.TkColon)
	expect(lexer.TkNewLine)

	statement(expectedIdent + 1)
}

func statement(expectedIdent int) {

	ignoreEmptyNewLines()

	for lexer.GetIdentLevel() == expectedIdent {
		ignoreEmptyNewLines()
		fmt.Printf("IDENTITY_LEVEL_HERE=%d\n", lexer.GetIdentLevel())
		switch tk := lexer.GetToken(); tk.Class {

		case lexer.TkResourceStmt:
			resource(tk.Value, expectedIdent)
		case lexer.TkIfStmt:
			fmt.Println("IF")
		case lexer.TkForStmt:
			fmt.Println("FOR")
		case lexer.TkSwitchStmt:
			fmt.Println("SWITCH")
		case lexer.TkPoint:
			expect(lexer.TkPoint)
			fmt.Printf("Resource Attribute --->>>> %s\n", lexer.GetToken().Value)
			expect(lexer.TkIdentifier)
			expect(lexer.TkEqual)
			expression()
		case lexer.TkIdentifier:
			expect(lexer.TkIdentifier)
			expect(lexer.TkEqual)
			expression()
		default:
			expression()
		}
	}

}

func resource(name string, expectedIdent int) {

	expect(lexer.TkResourceStmt)
	tkResourceName := lexer.GetToken()
	acceptAny(lexer.TkString, lexer.TkIdentifier)

	fmt.Printf("Resource : %s=%s\n", name, tkResourceName.Value)
	expect(lexer.TkColon)
	expect(lexer.TkNewLine)

	statementBlock(expectedIdent + 1)
}

func statementBlock(expectedIdent int) {

	ignoreEmptyNewLines()

	for lexer.GetIdentLevel() == expectedIdent {
		fmt.Printf("IDENTITY_LEVEL_HERE=%d\n", lexer.GetIdentLevel())
		statement(expectedIdent)
	}
}

func expression() {

	term()
	for acceptAny(
		lexer.TkPlus,
		lexer.TkMinus,
		lexer.TkAndOper,
		lexer.TkOrOper,
	) {
		term()
	}

}

func term() {

	factor()
	for acceptAny(
		lexer.TkEqual,
		lexer.TkNotEqual,
		lexer.TkGt,
		lexer.TkGte,
		lexer.TkLt,
		lexer.TkLte,
	) {
		factor()
	}
}

func factor() {

	switch tk := lexer.GetToken(); tk.Class {
	case lexer.TkIdentifier:
		expect(lexer.TkIdentifier)
	case lexer.TkString:
		expect(lexer.TkString)
	case lexer.TkInt:
		expect(lexer.TkInt)
	case lexer.TkBool:
		expect(lexer.TkBool)
	default:
		fmt.Printf("Unexpected token %s at factor\n", lexer.GetTokenText(tk.Class))
		os.Exit(-1)
	}
}
