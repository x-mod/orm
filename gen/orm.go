package gen

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/x-mod/cmd"
	"github.com/x-mod/errors"
)

//PasswordFunc
type PasswordFunc func(string) string

func md5Password(password string) string {
	algo := md5.New()
	if _, err := algo.Write([]byte(password)); err != nil {
		return ""
	}
	return hex.EncodeToString(algo.Sum([]byte("orm")))
}

//Cipher
type Cipher interface {
	Encode(string) string
	Decode(string) string
}

type b64Cipher struct{}

func (c *b64Cipher) Encode(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

func (c *b64Cipher) Decode(src string) string {
	decoded, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return ""
	}
	return string(decoded)
}

func PasswordCmd(c *cmd.Command, args []string) error {
	if len(args) > 0 {
		fmt.Println(md5Password(args[0]))
	}
	return nil
}
func EncodeCmd(c *cmd.Command, args []string) error {
	if len(args) > 0 {
		b := &b64Cipher{}
		fmt.Println(b.Encode(args[0]))
		return nil
	}

	//先取程序的标准输入属性信息
	info, err := os.Stdin.Stat()
	if err != nil {
		return errors.Annotate(err, "stdin stat failed")
	}
	// 判断标准输入设备属性 os.ModeCharDevice 是否设置
	// 同时判断是否有数据输入
	// if (info.Mode()&os.ModeCharDevice) == os.ModeCharDevice &&
	if info.Size() > 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return errors.Annotate(err, "stdin read failed")
		}
		b := &b64Cipher{}
		fmt.Fprint(os.Stdout, b.Encode(string(bytes)))
	}

	return nil
}
func DecodeCmd(c *cmd.Command, args []string) error {
	if len(args) > 0 {
		b := &b64Cipher{}
		fmt.Println(b.Decode(args[0]))
		return nil
	}

	//先取程序的标准输入属性信息
	info, err := os.Stdin.Stat()
	if err != nil {
		return errors.Annotate(err, "stdin stat failed")
	}

	// 判断标准输入设备属性 os.ModeCharDevice 是否设置
	// 同时判断是否有数据输入
	// if (info.Mode()&os.ModeCharDevice) == os.ModeCharDevice &&
	if info.Size() > 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return errors.Annotate(err, "stdin read failed")
		}
		b := &b64Cipher{}
		fmt.Fprint(os.Stdout, b.Decode(string(bytes)))
	}
	return nil
}
