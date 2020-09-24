package main

import (
	cmd "github.com/lichao-mobanche/go-extractor-server/cmd/extractor-server/command"
	"github.com/cfhamlet/os-rq-pod/pkg/command"
)

func main() {
	command.Execute(cmd.Root)
}