package gen

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"github.com/x-mod/errors"
	"github.com/x-mod/orm/tpl"
)

func getTemplate(prefix string) (*template.Template, error) {
	t := template.New("orm")
	for _, name := range tpl.AssetNames() {
		if strings.HasPrefix(name, prefix) {
			data, err := tpl.Asset(name)
			if err != nil {
				return nil, errors.Annotatef(err, "asset read %s", name)
			}
			_, err = t.Parse(string(data))
			if err != nil {
				return nil, errors.Annotatef(err, "asset parse %s", name)
			}
		}
	}
	return t, nil
}

func copyFilesExcludeSuffix(prefix string, suffix string, destDir string, force bool) error {
	fs := afero.NewOsFs()
	for _, name := range tpl.AssetNames() {
		if strings.HasPrefix(name, prefix) {
			dest := filepath.Join(destDir, strings.TrimPrefix(name, prefix))
			if !strings.HasSuffix(dest, suffix) {
				dir, _ := filepath.Split(dest)
				if err := fs.MkdirAll(dir, 0777); err != nil {
					return errors.Annotatef(err, "mkdir %s", dir)
				}
				exist, err := afero.Exists(fs, dest)
				if err != nil {
					return errors.Annotatef(err, "exist %s", dest)
				}
				if exist && !force {
					continue
				}
				if exist {
					if err := fs.Remove(dest); err != nil {
						return errors.Annotatef(err, "rm %s", dest)
					}
				}
				fd, err := fs.Create(dest)
				if err != nil {
					return errors.Annotatef(err, "create %s", dest)
				}
				defer fd.Close()

				data, err := tpl.Asset(name)
				if err != nil {
					return errors.Annotatef(err, "asset read %s", name)
				}
				if _, err := fd.Write(data); err != nil {
					return errors.Annotatef(err, "write %s", dest)
				}
			}
		}
	}
	return nil
}

func getInputFilesBySuffix(dir string, suffix string) ([]string, error) {
	stat, err := os.Stat(dir)

	if err != nil && os.IsNotExist(err) {
		return nil, err
	}

	files := []string{}
	if !stat.IsDir() && strings.HasSuffix(dir, suffix) {
		files = append(files, dir)
		return files, nil
	}

	if err := filepath.Walk(dir, func(src string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			files = append(files, src)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}
