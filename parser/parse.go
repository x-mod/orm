package parser

import (
	"github.com/x-mod/orm/object"
	"github.com/x-mod/orm/yaml"
	yamlv2 "gopkg.in/yaml.v2"
)

func Parse(filename string) ([]*object.Object, error) {
	bs, err := yaml.Split(filename)
	if err != nil {
		return nil, err
	}

	objs := []*object.Object{}
	for _, b := range bs {
		obj := object.Object{}
		err := yamlv2.Unmarshal(b.Bytes(), &obj)
		if err != nil {
			return nil, err
		}
		obj.Adjust()
		objs = append(objs, &obj)
	}
	return objs, nil
}
