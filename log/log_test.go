package log_test

import (
	"anto.pt/x/log"
)

func Example() {
	l := log.Module("test")
	l.Info("this is info")
	l.Warn("this is warn")
	l.Error("this is error")

	l2 := log.Module("module_2")
	l2.Error("this is error")
}
