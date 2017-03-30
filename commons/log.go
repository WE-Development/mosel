package commons

import "log"

func LogFatal(err error) {
	if err == nil {
		return
	}

	log.Fatal(err)
}
