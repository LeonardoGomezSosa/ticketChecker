package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func x() {
	str := "^(tr)-\\d{4}"
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingresar entrada: ")
	entrada, _ := reader.ReadString('\n')
	fmt.Println("Ingresar REGEX: ")
	str, _ = reader.ReadString('\n')
	matched, err := regexp.MatchString(str, entrada)
	fmt.Println(matched, err)
}
