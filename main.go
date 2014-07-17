// prosit project main.go
package main

import (
	"log"
	"os"
	//"path/filepath"
	"fmt"
	"prosit/cl"
	"prosit/web"
	"strings"
)

func main() {

	//thisFile, _ := filepath.Abs(os.Args[0])

	if len(os.Args) == 1 {
		log.Printf("Starting as daemon process\n")
		web.StartWeb(9999)
		return
	}

	var err error

	switch strings.ToLower(os.Args[1]) {
	case "add-process":
		err = cl.AddProcessCL()
	case "list-processes":
		err = cl.ListProcessesCL()
	case "stop-process":
		err = cl.StopProcessCL()
	}

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
