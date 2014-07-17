package cl

import (
	"bufio"
	"fmt"
	"os"
)

func readLine(question, defaultValue string) string {

	if defaultValue != "" {
		fmt.Printf("%s: [%s] ", question, defaultValue)
	} else {
		fmt.Printf("%s: ", question)
	}

	bio := bufio.NewReader(os.Stdin)
	line, _, _ := bio.ReadLine()

	lineStr := string(line)

	if lineStr == "" {
		return defaultValue
	} else {
		return lineStr
	}

}
