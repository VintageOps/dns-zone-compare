package utils

import "log"

func FatalOnErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
