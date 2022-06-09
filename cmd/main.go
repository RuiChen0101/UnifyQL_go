package main

import (
	"fmt"
	"regexp"
)

func main() {

	c, err := regexp.MatchString(`=|!=|<|<=|>|>=|LIKE`, "=")
	fmt.Println(c)
	fmt.Println(err)
}
