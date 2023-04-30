package helpers

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "", log.Lshortfile)
