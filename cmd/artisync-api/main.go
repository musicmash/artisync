package main

import (
	"bufio"
	"os"

	"github.com/musicmash/artisync/internal/log"
)

func main() {
	log.SetLevel("INFO")
	log.SetWriters(log.GetConsoleWriter())

	log.Info("hello, world")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
