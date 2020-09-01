package lexer

import (
	"fmt"
	"io"
	"os"
	"regexp"
)

// Lexer datastructure
type Lexer struct {
	buffer   []byte
	index    int
	size     int
	pos      int
	line     int
	fd       *os.File
	fileName string
	fdEnd    bool
	tokens   []Token
}

// Token stores the information
// about the tokens captured.
type Token struct {
	Class  TkClassType
	Value  string
	Length int
}

// Value **might** hold the values for tokens
// not sure about this one yet.
type Value struct {
	t int
	v struct {
		str    string
		number int
	}
}

// TkClassType holds the values for
// every type of token defined.
type TkClassType int

var lexer Lexer

var reservedWords map[string]TkClassType = map[string]TkClassType{
	"class":   TkClassDef,
	"file":    TkResourceFile,
	"package": TkResourcePackage,
	"service": TkResourceService,
}

// Holds the text representation of tokens.
// It's useful to provide error messages for users.
var tokenText map[TkClassType]string = map[TkClassType]string{
	TkUndefined:       "undefined",
	TkEOF:             "eof",
	TkString:          "string",
	TkInt:             "integer",
	TkIdentifier:      "identifier",
	TkNewLine:         "new line",
	TkColon:           ":",
	TkIdent:           "ident",
	TkDedent:          "dedent",
	TkEqual:           "=",
	TkPlus:            "+",
	TkClassDef:        "class",
	TkComma:           ",",
	TkPoint:           ".",
	TkLeftParentenses: "(",
	TkRightParenteses: ")",
	TkHashMark:        "#",
	TkLeftBrackets:    "[",
	TkRightBrackets:   "]",
	TkLeftBraces:      "{",
	TkRightBraces:     "}",
	TkResourceFile:    "file",
	TkResourcePackage: "package",
	TkResourceService: "service",
}

const (
	// TkUndefined is the default token type
	TkUndefined TkClassType = -1
	// TkEOF means literary EOF
	TkEOF TkClassType = 0
	// TkString token of type string
	// Strings are any caracters delimited by
	// either " (double quote) or ' (single quote)
	TkString TkClassType = 1
	// TkInt token of type integer
	TkInt TkClassType = 2
	// TkIdentifier token of type identifier.
	// Follows the `[a-zA-Z_]{1}[a-zA-Z0-9_\-]+` regex
	TkIdentifier TkClassType = 3
	// TkNewLine match either \n or \n\r
	TkNewLine TkClassType = 4
	// TkColon matches ':' symbol
	TkColon  TkClassType = 5
	TkIdent  TkClassType = 6
	TkDedent TkClassType = 7
	// TkEqual matches the '=' sign
	TkEqual           TkClassType = 8
	TkPlus            TkClassType = 9
	TkClassDef        TkClassType = 10
	TkComma           TkClassType = 11
	TkPoint           TkClassType = 20
	TkLeftParentenses TkClassType = 21
	TkRightParenteses TkClassType = 22
	TkHashMark        TkClassType = 23
	TkLeftBrackets    TkClassType = 24
	TkRightBrackets   TkClassType = 25
	TkLeftBraces      TkClassType = 26
	TkRightBraces     TkClassType = 27
	TkResourceFile    TkClassType = 50
	TkResourcePackage TkClassType = 51
	TkResourceService TkClassType = 52
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (l *Lexer) createToken(class TkClassType, value string) Token {
	newTk := Token{
		class,
		value,
		len(value),
	}
	return newTk
}

func (l *Lexer) readToBuffer() {
	l.buffer = make([]byte, 1024)
	count, err := l.fd.Read(l.buffer)

	if err == io.EOF {
		l.fdEnd = true
		return
	}

	check(err)

	l.size = count
	l.index = -1

	fmt.Printf("--> Read %d bytes\n", l.size)
}

func (l *Lexer) error(expected string) {
	fmt.Printf("Expected char '%s' but get '%s' at position %d, line %d\n", expected, l.getChar(), l.pos, l.line)
	os.Exit(-1)
}

func (l *Lexer) getChar() string {
	return string(l.buffer[l.index])
}

func (l *Lexer) nextChar() {
	l.index++
	l.pos++
	if l.index >= l.size {
		l.readToBuffer()
		l.index = 0
		l.pos = 0
	}
}

// prevChar is the oposite of nextChar.
// Instead of advance the cursor for the next
// caracter, it get the previuos caracter
func (l *Lexer) prevChar() {
	l.index--
	l.pos--
}

func (l *Lexer) match(char string) bool {

	if char == l.getChar() {
		return true
	}

	return false
}

func (l *Lexer) matchAny(chars ...string) bool {

	for _, v := range chars {
		if v == l.getChar() {
			return true
		}
	}
	return false
}

func (l *Lexer) matchIfReservedWord(s string) (TkClassType, bool) {

	for i, v := range reservedWords {
		if i == s {
			return v, true
		}
	}

	return TkUndefined, false
}

func (l *Lexer) lookAheadMatch(char string) bool {
	lookAheadChar := string(l.buffer[l.index+1])

	if char == lookAheadChar {
		l.index++
		return true
	}

	return false
}

func isInt(char string) bool {

	matched, err := regexp.Match("[0-9]", []byte(char))
	check(err)

	return matched
}

func isIdentifier(char string) bool {

	matched, err := regexp.Match("[a-zA-Z_]", []byte(char))
	check(err)

	return matched
}

func (l *Lexer) captureInt() string {
	return ""
}

func (l *Lexer) captureDoubleQuoteString() string {
	retString := ""
	l.match("\"")
	l.nextChar()
	for {
		char := l.getChar()
		if char == "\"" {
			break
		}
		retString += char
		l.nextChar()
	}
	l.match("\"")
	return retString
}

func (l *Lexer) captureSingleQuoteString() string {
	retString := ""
	l.match("'")
	l.nextChar()
	for {
		char := l.getChar()
		if char == "'" {
			break
		}
		retString += char
		l.nextChar()
	}
	l.match("'")
	return retString
}

func (l *Lexer) captureIdentifier() string {
	retIdentifier := ""

	for {
		char := l.getChar()
		if isIdentifier(char) {
			retIdentifier += char
		} else {
			break
		}
		l.nextChar()
	}
	l.prevChar()
	return retIdentifier
}

func (l *Lexer) ignoreWhiteSpaces() {
	l.nextChar()
	for {
		char := l.getChar()
		if char != " " {
			break
		}
		l.nextChar()
	}
}

// Init initilize the global variable
// lexer.
func Init(fileName string) {
	fmt.Printf("Opening file %s\n", fileName)
	fd, err := os.Open(fileName)
	check(err)

	// initialize lexer datastructure
	lexer.buffer = nil
	lexer.index = -1
	lexer.size = -1
	lexer.line = 1
	lexer.pos = 0
	lexer.fd = fd
	lexer.fileName = fileName

	lexer.readToBuffer()
}

func (l *Lexer) pushToken(c TkClassType, v string) {
	tk := l.createToken(c, v)
	l.tokens = append(l.tokens, tk)
}

// GetTokenText gets text value of tokens
func GetTokenText(c TkClassType) string {
	return tokenText[c]
}

// GetToken is invoked by parser to see the last
// captured token
func GetToken() Token {
	return lexer.tokens[len(lexer.tokens)-1]
}

func GetTokenCurrentPosition() int {
	return lexer.pos
}

func GetTokenCurrentLine() int {
	return lexer.line
}


func 
// NextToken is invoked by parser.Parse()
// to capture the next token (if available)
func NextToken() {

	// check for EOF flag
	if lexer.fdEnd {
		lexer.pushToken(TkEOF, "")
		return
	}

	// if the last token was a new line, 
	// we a going to count the number of idents
	if lexer.getChar().Class == lexer.TkNewLine {

	}

	lexer.ignoreWhiteSpaces()
	char := lexer.getChar()

	if isInt(char) {
		lexer.pushToken(TkInt, lexer.captureInt())
		return
	}

	if isIdentifier(char) {
		identifier := lexer.captureIdentifier()
		class, found := lexer.matchIfReservedWord(identifier)

		if !found {
			class = TkIdentifier
		}

		lexer.pushToken(class, identifier)
		return
	}

	switch char {
	case "\"":
		lexer.pushToken(TkString, lexer.captureDoubleQuoteString())
	case "'":
		lexer.pushToken(TkString, lexer.captureSingleQuoteString())
	case "\n":
		lexer.match("\n")
		// supporting LFCR
		lexer.lookAheadMatch("\r")
		lexer.pos = 0
		lexer.line++
		lexer.pushToken(TkNewLine, "NEW_LINE")
	case ",":
		lexer.match(",")
		lexer.pushToken(TkComma, ",")
	case "+":
		lexer.match("+")
		lexer.pushToken(TkPlus, "+")
	case ":":
		lexer.match(":")
		lexer.pushToken(TkColon, ":")
	case "=":
		lexer.match("=")
		lexer.pushToken(TkEqual, "=")
	case "#":
		lexer.match("#")
		lexer.pushToken(TkHashMark, "#")
	case "(":
		lexer.match("(")
		lexer.pushToken(TkLeftParentenses, "(")
	case ")":
		lexer.match(")")
		lexer.pushToken(TkRightParenteses, ")")
	case "[":
		lexer.match("[")
		lexer.pushToken(TkLeftBrackets, "[")
	case "]":
		lexer.match("]")
		lexer.pushToken(TkRightBrackets, "]")
	case "{":
		lexer.match("{")
		lexer.pushToken(TkLeftBraces, "{")
	case "}":
		lexer.match("}")
		lexer.pushToken(TkRightBraces, "}")
	default:
		lexer.pushToken(TkUndefined, "")
	}
}
