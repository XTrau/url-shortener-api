package main

import (
	"fmt"
	"urlshortener/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}
