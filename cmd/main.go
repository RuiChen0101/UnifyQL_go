package main

import (
	"fmt"
)

func main() {

	m := map[string]interface{}{
		"a": 123,
		"b": "456",
	}
	a, ok := m["b"].(string)
	fmt.Println(a)
	fmt.Println(ok)
}
