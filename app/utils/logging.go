package utils

import "fmt"

func LogBytesInHex(data []byte) {
	for i, b := range data {
		hex := "0" + fmt.Sprintf("%x", b)
		hex = hex[len(hex)-2:]

		fmt.Printf("%s ", hex)
		if i%8 == 7 {
			fmt.Println()
		}
	}

	if len(data)%8 != 0 {
		fmt.Println()
	}
}
