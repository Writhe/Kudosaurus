package utils

import (
	"encoding/json"
	"fmt"
)

// PrettyPrint - prints prettily
func PrettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(s))
}
