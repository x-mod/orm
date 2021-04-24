package main

import (
	"github.com/x-mod/build"
	"github.com/x-mod/cmd"
	_ "github.com/x-mod/orm/gen"
)

//go:generate go-bindata -prefix tpl -nometadata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl/...
func main() {
	cmd.Version(build.String())
	cmd.PersistentFlags().StringP("input", "i", ".", "input directory")
	cmd.PersistentFlags().StringP("output", "o", ".", "output directory")
	cmd.PersistentFlags().StringP("template-suffix", "t", ".gogo", "template suffix")
	cmd.PersistentFlags().StringP("go-package-name", "p", "model", "go package name")
	cmd.Execute()
}
