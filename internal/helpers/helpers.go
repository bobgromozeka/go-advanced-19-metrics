package helpers

import (
	"log"
	"strconv"
)

func StrToInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Println("Error when converting string to int: ", err)
		return 0 //ignore error and return 0
	}

	return v
}
