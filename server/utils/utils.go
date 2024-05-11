package utils

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "[dumgo]", log.Ldate|log.Ltime)
