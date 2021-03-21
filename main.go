package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/riadafridishibly/go-monkey/repl"
)

func main() {
	currentUser, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hello %s!\n", currentUser.Username)

	repl.Start(os.Stdin, os.Stdout)
}
