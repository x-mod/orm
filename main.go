package main

import (
	"embed"

	"github.com/x-mod/build"
	"github.com/x-mod/cmd"
	"github.com/x-mod/dir"
	_ "github.com/x-mod/orm/gen"
)

//go:embed tpl/*
var embeddir embed.FS

func main() {
	dir.EmbedFS = embeddir
	cmd.Version(build.String())
	cmd.PersistentFlags().StringP("workdir", "c", "", "orm workdir (default: $HOME/.orm)")
	cmd.PersistentFlags().StringP("input", "i", ".", "input directory")
	cmd.PersistentFlags().StringP("output", "o", ".", "output directory")
	cmd.PersistentFlags().StringP("template-suffix", "t", ".gogo", "template suffix")
	cmd.PersistentFlags().StringP("go-package-name", "p", "model", "go package name")
	cmd.Execute()
}
