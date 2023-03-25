package main

import (
	"encoding/json"
	"fmt"
)

func PrintAsJSON(in any) {
	b, _ := json.MarshalIndent(in, "\t", "")
	fmt.Println(string(b))
}
