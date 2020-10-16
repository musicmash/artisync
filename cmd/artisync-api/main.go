package main

import (
	"bufio"
	"os"

	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/version"
)

func main() {
	log.SetLevel("DEBUG")
	log.SetWriters(log.GetConsoleWriter())

	log.Debug(version.FullInfo)
	log.Info("hello, world")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
