package main

import (
	"github.com/bannovdaniil/zeronews"
	"log"
)

func main() {
	server := new(zeronews.Server)
	if err := server.Run("8000"); err != nil {
		log.Fatalf("server error: %s", err.Error())
	}
}
