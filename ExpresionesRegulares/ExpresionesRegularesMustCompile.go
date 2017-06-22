package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	str := "^(tr)-\\d{4}"
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingresar entrada: ")
	entrada, _ := reader.ReadString('\n')
	fmt.Println("Ingresar REGEX: ")
	str, _ = reader.ReadString('\n')
	re := regexp.MustCompile(str)
	matched := re.MatchString(entrada)
	fmt.Println(matched)
}
