package pgx

import (
	"fmt"
)

type Logger struct{}

func (l Logger) Printf(format string, v ...interface{}) {
	fmt.Printf("%+v\n", v)
}

func (l Logger) Verbose() bool {
	return false
}
