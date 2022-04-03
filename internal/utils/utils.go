package utils

import (
	"encoding/json"
	"fmt"
)

// MapToStruct converts a map to a struct format
func MapToStruct(mapData interface{}, structData interface{}) {
	// Convert map to json string
	jsonStr, err := json.Marshal(mapData)
	if err != nil {
		fmt.Println(err)
	}

	// Convert json string to struct
	if err := json.Unmarshal(jsonStr, &structData); err != nil {
		fmt.Println(err)
	}
}
