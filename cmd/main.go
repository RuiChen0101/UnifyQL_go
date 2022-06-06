package main

import (
	"fmt"
	"regexp"
)

func main() {

	reg := regexp.MustCompile(`\s*,\s*`)
	split := reg.Split("tableB ,  tableC ,  tableD", -1)
	for _, n := range split {
		fmt.Println(n)
	}
}
