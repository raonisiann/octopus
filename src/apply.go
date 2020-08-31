package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
)

type Lexer struct {
	buffer   []byte
	index    int
	size     int
	fd       *os.File
	fileName string
}

type Token int

var lexer Lexer

var reservedWords map[string]Token = map[string]Token{
	"file":    tkResourceFile,
	"package": tkResourcePackage,
}

const (
	tkUndefined       Token = -1
	tkEOF             Token = 0
	tkString          Token = 1
	tkInt             Token = 2
	tkIdentifier      Token = 3
	tkResourceFile    Token = 30
	tkResourcePackage Token = 31
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (l *Lexer) readFile(fileName string) {
	fmt.Printf("Opening file %s\n", fileName)
	f, err := os.Open(fileName)
	check(err)
	buffer := make([]byte, 10)

	for {
		count, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		check(err)
		fmt.Printf("Bytes read: %d\n", count)
		fmt.Printf("%v\n", buffer)
	}

}

func (l *Lexer) init(fileName string) {
	fmt.Printf("Opening file %s\n", fileName)
	fd, err := os.Open(fileName)
	check(err)

	// initialize lexer datastructure
	l.buffer = make([]byte, 1024)
	l.index = -1
	l.size = -1
	l.fd = fd
	l.fileName = fileName

	l.readToBuffer()
	fmt.Printf("Finished first read %d bytes\n", l.size)
}

func (l *Lexer) readToBuffer() {
	count, err := l.fd.Read(l.buffer)

	if err == io.EOF {
		fmt.Println("<<EOF")
		os.Exit(-1)
	}

	check(err)

	l.size = count
	l.index = 0
}

func (l *Lexer) getChar() string {
	fmt.Printf("%v\n", l.buffer)
	return string(l.buffer[l.index])
}

func (l *Lexer) nextChar() {
	l.index++

	if l.index >= l.size {
		l.readToBuffer()
	}
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

func captureInt() string {
	return ""
}

func captureString() string {
	return ""
}

func captureIdentifier() string {
	retIdentifier := ""

	for {
		char := lexer.getChar()
		if isIdentifier(char) {
			retIdentifier += char
		} else {
			break
		}
	}

	return retIdentifier
}

func (l *Lexer) getToken() Token {
	char := l.getChar()

	if isInt(char) {
		fmt.Println(captureInt())
		return tkInt
	}

	if isIdentifier(char) {
		fmt.Println(captureIdentifier())
		return tkIdentifier
	}

	switch char {
	case "\"":
		fmt.Println("String")
		return tkString
	}

	return tkUndefined
}

func parse() {

	for {
		tk := lexer.getToken()
		if tk == tkEOF {
			break
		}
	}
}

// ApplyCmd get the arguments from apply command
func ApplyCmd(args []string) {

	if len(args) == 0 {
		fmt.Println("File required")
		os.Exit(-1)
	}

	lexer.init(args[0])
	fmt.Printf("First char %s\n", lexer.getChar())
	//parse()
}
