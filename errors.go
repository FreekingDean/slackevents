package slackevents

import (
	"log"
)

type ErrorHandler func(error) (int, string)

func DefaultErrorHandler(err error) (int, string) {
	log.Println(err.Error())
	return 500, ""
}
