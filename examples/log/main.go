package main

import (
	"fmt"
	"github.com/JoinVerse/obs/log"
)

func main() {

	log.Info("hello world")
	log.Infof("hello %s", "world")

	err := fmt.Errorf("be water my friend")
	log.Error("He said", err)

	log.Fatal("hello world", err)
}
