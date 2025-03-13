package main

import (
	"go-kitsx/log/logx"
)

func main() {
	logx.Info("hello", logx.String("name", "world"))
	logx.Debug("hello", logx.String("name", "world"))

}
