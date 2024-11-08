package main

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
)

func main() {

	for range 30 {
		first := gofakeit.FirstName()
		last := gofakeit.LastName()
		fmt.Printf("%v,%v\n", first, last)
	}

}
