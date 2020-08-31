package parser

import (
	"fmt"
	"lexer"
)

func Parse(fileName string) {

	lexer.Init(fileName)

	for {
		tk := lexer.GetToken()
		fmt.Printf("(%d, '%s', %d)\n", tk.Class, tk.Value, tk.Length)
		// 0 == EOF
		if tk.Class == 0 {
			break
		}

	}
}
