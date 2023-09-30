package main

import (
	"fmt"

	"github.com/tim117/chatserver/cmd"
)

func main() {
	fmt.Println("Sup dude")
	db := Connect()
	SeedData(db)
	cmd.Execute(db)
}
