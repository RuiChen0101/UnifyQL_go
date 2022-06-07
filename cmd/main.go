package main

import (
	"fmt"
)

func main() {

	table := map[string]bool{
		"aaaa": true,
	}
	_, o := table["bbbbb"]
	fmt.Println(o)
	fmt.Println(table["bbbbb"])
}
