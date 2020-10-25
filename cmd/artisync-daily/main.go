package main

import (
	"bufio"
	"os"

	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/version"
)

func main() {
	log.SetLevel("INFO")
	log.SetWriters(log.GetConsoleWriter())

	log.Debug(version.FullInfo)

	log.Info("artisync-daily is running...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
