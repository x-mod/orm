package gen

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"github.com/x-mod/dir"
	"github.com/x-mod/errors"
)

func getTemplate(root *dir.Dir, elems ...string) (*template.Template, error) {
	t := template.New("orm")

	if err := filepath.Walk(root.Path(elems...), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			t, err = t.New(path).Parse(string(data))
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	// for _, name := range tpl.AssetNames() {
	// 	if strings.HasPrefix(name, prefix) {
	// 		data, err := tpl.Asset(name)
	// 		if err != nil {
	// 			return nil, errors.Annotatef(err, "asset read %s", name)
	// 		}
	// 		_, err = t.Parse(string(data))
	// 		if err != nil {
	// 			return nil, errors.Annotatef(err, "asset parse %s", name)
	// 		}
	// 	}
	// }
	return t, nil
}

func copyFilesExcludeSuffix(root *dir.Dir, prefix string, suffix string, destDir string, force bool) error {
	osfs := afero.NewOsFs()
	return filepath.Walk(root.Path(prefix), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			dest := filepath.Join(destDir, strings.TrimPrefix(path, root.Path(prefix)))
			if !strings.HasSuffix(dest, suffix) {
				dir, _ := filepath.Split(dest)
				if err := osfs.MkdirAll(dir, 0777); err != nil {
					return errors.Annotatef(err, "mkdir %s", dir)
				}
				exist, err := afero.Exists(osfs, dest)
				if err != nil {
					return errors.Annotatef(err, "exist %s", dest)
				}
				if exist && !force {
					return nil
				}
				if exist {
					if err := osfs.Remove(dest); err != nil {
						return errors.Annotatef(err, "rm %s", dest)
					}
				}
				fd, err := osfs.Create(dest)
				if err != nil {
					return errors.Annotatef(err, "create %s", dest)
				}
				defer fd.Close()

				data, err := ioutil.ReadFile(path)
				if err != nil {
					return errors.Annotatef(err, "read %s", path)
				}
				if _, err := fd.Write(data); err != nil {
					return errors.Annotatef(err, "write %s", dest)
				}
			}
		}
		return nil
	})

	// fs := afero.NewOsFs()
	// for _, name := range tpl.AssetNames() {
	// 	if strings.HasPrefix(name, prefix) {
	// 		dest := filepath.Join(destDir, strings.TrimPrefix(name, prefix))
	// 		if !strings.HasSuffix(dest, suffix) {
	// 			dir, _ := filepath.Split(dest)
	// 			if err := fs.MkdirAll(dir, 0777); err != nil {
	// 				return errors.Annotatef(err, "mkdir %s", dir)
	// 			}
	// 			exist, err := afero.Exists(fs, dest)
	// 			if err != nil {
	// 				return errors.Annotatef(err, "exist %s", dest)
	// 			}
	// 			if exist && !force {
	// 				continue
	// 			}
	// 			if exist {
	// 				if err := fs.Remove(dest); err != nil {
	// 					return errors.Annotatef(err, "rm %s", dest)
	// 				}
	// 			}
	// 			fd, err := fs.Create(dest)
	// 			if err != nil {
	// 				return errors.Annotatef(err, "create %s", dest)
	// 			}
	// 			defer fd.Close()

	// 			data, err := tpl.Asset(name)
	// 			if err != nil {
	// 				return errors.Annotatef(err, "asset read %s", name)
	// 			}
	// 			if _, err := fd.Write(data); err != nil {
	// 				return errors.Annotatef(err, "write %s", dest)
	// 			}
	// 		}
	// 	}
	// }
	// return nil
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
