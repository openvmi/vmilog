package main

import "github.com/openvmi/vmilog"

func main() {
	vmilog.EnableConsole(true)
	vmilog.EnableFile(true)
	vmilog.SetLevel("error")
	p := 10.123
	vmilog.Info("main", "hello main", "second ", p)
	vmilog.Error("main", "this is error msg")
	vmilog.Warn("main", "this is warn msg")
}
