package utils

import (
	"encoding/json"
	"fmt"
)

func Print(i interface{}) {
	marshal, err := json.Marshal(i)
	if err != nil {
		return
	}

	fmt.Println(string(marshal))
}
