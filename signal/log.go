package signal

import (
	"log"
	"os"
)

var (
	errLog = log.New(os.Stderr, "[Error] ", 0)
)
