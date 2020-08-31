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
	"file":    tkResourceFile,
	"package": tkResourcePackage,
	"service": tkResourceService,
}

const (
	tkUndefined       tkClass = -1
	tkEOF             tkClass = 0
	tkString          tkClass = 1
	tkInt             tkClass = 2
	tkIdentifier      tkClass = 3
	tkNewLine         tkClass = 4
	tkColon           tkClass = 5
	tkIdent           tkClass = 6
	tkDedent          tkClass = 7
	tkEqual           tkClass = 8
	tkPlus            tkClass = 9
	tkPoint           tkClass = 10
	tkLeftParentenses tkClass = 11
	tkRightParenteses tkClass = 12
	tkHashMark        tkClass = 13
	tkLeftBrackets    tkClass = 14
	tkRightBrackets   tkClass = 15
	tkLeftBraces      tkClass = 16
	tkRightBraces     tkClass = 17
	tkResourceFile    tkClass = 50
	tkResourcePackage tkClass = 51
	tkResourceService tkClass = 52
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

// GetToken is invoked by parser.Parse()
// to request tokens from the input
func GetToken() Token {
	// check for EOF flag
	if lexer.fdEnd {
		return lexer.createToken(tkEOF, "")
	}

	lexer.ignoreWhiteSpaces()
	char := lexer.getChar()

	if isInt(char) {
		return lexer.createToken(tkInt, lexer.captureInt())
	}

	if isIdentifier(char) {
		return lexer.createToken(tkIdentifier, lexer.captureIdentifier())
	}

	switch char {
	case "\"":
		return lexer.createToken(tkString, lexer.captureDoubleQuoteString())
	case "'":
		return lexer.createToken(tkString, lexer.captureSingleQuoteString())
	case "\n":
		lexer.match("\n")
		// supporting LFCR
		lexer.lookAheadMatch("\r")
		lexer.pos = 0
		lexer.line++
		return lexer.createToken(tkNewLine, "NEW_LINE")
	case "+":
		lexer.match("+")
		return lexer.createToken(tkPlus, "+")
	case ":":
		lexer.match(":")
		return lexer.createToken(tkColon, ":")
	case "=":
		lexer.match("=")
		return lexer.createToken(tkEqual, "=")
	case "#":
		lexer.match("#")
		return lexer.createToken(tkHashMark, "#")
	case "(":
		lexer.match("(")
		return lexer.createToken(tkLeftParentenses, "(")
	case ")":
		lexer.match(")")
		return lexer.createToken(tkRightParenteses, ")")
	case "[":
		lexer.match("[")
		return lexer.createToken(tkLeftBrackets, "[")
	case "]":
		lexer.match("]")
		return lexer.createToken(tkRightBrackets, "]")
	case "{":
		lexer.match("{")
		return lexer.createToken(tkLeftBraces, "{")
	case "}":
		lexer.match("}")
		return lexer.createToken(tkRightBraces, "}")
	}

	return lexer.createToken(tkUndefined, "")
}
