package main

import (
	"gitlab.snapp.ir/platform/snapp_object_store/internal/cmd"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd.Execute()
}
