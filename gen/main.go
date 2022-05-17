package gen

import (
	"fmt"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/x-mod/cmd"
	"github.com/x-mod/dir"
)

var packageName = "model"
var inputSuffixes = []string{"yml", "yaml"}

func init() {
	cmd.Add(
		cmd.Path("/script/mysql"),
		cmd.Short("generate mysql scripts"),
		cmd.Main(ScriptMySQL),
	)
	cmd.Add(
		cmd.Path("/code/golang"),
		cmd.Short("generate golang codes"),
		cmd.Main(CodeGolang),
	)
	cmd.Add(
		cmd.Path("/password"),
		cmd.Short("orm password string"),
		cmd.Main(PasswordCmd),
	)
	cmd.Add(
		cmd.Path("/encode"),
		cmd.Short("orm encode string"),
		cmd.Main(EncodeCmd),
	)
	cmd.Add(
		cmd.Path("/decode"),
		cmd.Short("orm decode string"),
		cmd.Main(DecodeCmd),
	)
}

func ScriptMySQL(c *cmd.Command, args []string) error {
	root := viper.GetString("workdir")
	if root == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		root = path.Join(home, ".orm")
	}
	workdir := dir.New(
		dir.Root(root),
	)
	if err := workdir.EmbedSkip(1, "tpl"); err != nil {
		return fmt.Errorf("embed orm: %v", err)
	}
	in := dir.New(dir.Root(viper.GetString("input")))
	if err := in.Open(); err != nil {
		return err
	}
	out := dir.New(dir.Root(viper.GetString("output")))
	if err := out.Open(); err != nil {
		return err
	}
	suffix := viper.GetString("template-suffix")
	packageName = viper.GetString("go-package-name")
	return generateScripts(workdir, in.Path(), out.Path(), suffix, "scripts", "mysql")
}

func CodeGolang(c *cmd.Command, args []string) error {
	root := viper.GetString("workdir")
	if root == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		root = path.Join(home, ".orm")
	}
	workdir := dir.New(
		dir.Root(root),
	)
	if err := workdir.EmbedSkip(1, "tpl"); err != nil {
		return fmt.Errorf("embed orm: %v", err)
	}
	in := dir.New(dir.Root(viper.GetString("input")))
	if err := in.Open(); err != nil {
		return err
	}
	out := dir.New(dir.Root(viper.GetString("output")))
	if err := out.Open(); err != nil {
		return err
	}
	suffix := viper.GetString("template-suffix")
	packageName = viper.GetString("go-package-name")
	return generateCode(workdir, in.Path(), out.Path(), suffix, "code", "golang")
}
