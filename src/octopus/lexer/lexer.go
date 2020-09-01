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
	Class  tkClass
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

type tkClass int

var lexer Lexer

var reservedWords map[string]tkClass = map[string]tkClass{
	"file":    TkResourceFile,
	"package": TkResourcePackage,
	"service": TkResourceService,
}

// Holds the text representation of tokens.
// It's useful to provide error messages for users.
var tokenText map[tkClass]string = map[tkClass]string{
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
	TkUndefined       tkClass = -1
	TkEOF             tkClass = 0
	TkString          tkClass = 1
	TkInt             tkClass = 2
	TkIdentifier      tkClass = 3
	TkNewLine         tkClass = 4
	TkColon           tkClass = 5
	TkIdent           tkClass = 6
	TkDedent          tkClass = 7
	TkEqual           tkClass = 8
	TkPlus            tkClass = 9
	TkClassDef        tkClass = 10
	TkPoint           tkClass = 20
	TkLeftParentenses tkClass = 21
	TkRightParenteses tkClass = 22
	TkHashMark        tkClass = 23
	TkLeftBrackets    tkClass = 24
	TkRightBrackets   tkClass = 25
	TkLeftBraces      tkClass = 26
	TkRightBraces     tkClass = 27
	TkResourceFile    tkClass = 50
	TkResourcePackage tkClass = 51
	TkResourceService tkClass = 52
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (l *Lexer) createToken(class tkClass, value string) Token {
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

func (l *Lexer) pushToken(class tkClass, value string) {
	tk := l.createToken(class, value)
	l.tokens = append(l.tokens, tk)
}

// GetTokenText gets text value of tokens
func GetTokenText(class tkClass) string {
	return tokenText[class]
}

// GetToken is invoked by parser to see the last
// captured token
func GetToken() Token {
	return lexer.tokens[len(lexer.tokens)-1]
}

// NextToken is invoked by parser.Parse()
// to capture the next token (if available)
func NextToken() {
	// check for EOF flag
	if lexer.fdEnd {
		lexer.pushToken(TkEOF, "")
		return
	}

	lexer.ignoreWhiteSpaces()
	char := lexer.getChar()

	if isInt(char) {
		lexer.pushToken(TkInt, lexer.captureInt())
		return
	}

	if isIdentifier(char) {
		lexer.pushToken(TkIdentifier, lexer.captureIdentifier())
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
