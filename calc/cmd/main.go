package main

import (
	"log"

	calc "github.com/wonderf00l/go-park-vk/calc/pkg"
)

func main() {
	config := calc.NewConfig()
	err := calc.Calculator(config.InputSteam, config.OutputStream, config)
	if err != nil {
		log.Fatalln(err)
	}
}
