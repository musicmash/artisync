package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello, world")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
