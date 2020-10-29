package main

import (
	"fmt"
	"github.com/JoinVerse/obs"
)

func main() {
	conf := obs.Config{
		NOGCloudEnabled: true,
	}

	observer := obs.New(conf)
	defer observer.Close()

	err := fmt.Errorf("main: ups, that was an error")
	observer.ErrorTags("ouch", map[string]string{"key": "value"}, err)
}
