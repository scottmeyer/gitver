package log

import (
	"fmt"
	"log"
)

var DebugEnabled = true

func debugPrefix() string {
	return "DEBUG: "
}

func Debugf(format string, v ...interface{}) {
	if DebugEnabled {
		Println(fmt.Sprintf("%s%s", debugPrefix(), fmt.Sprintf(format, v...)))
	}
}

func Debugln(v ...interface{}) {
	if DebugEnabled {
		Println(append([]interface{}{debugPrefix()}, v)...)
	}
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Println(v ...interface{}) {
	log.Println(v...)
}

type Formatter interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
