# orm

## 快速开始

- YAML 定义数据结构

```yaml
version: v1
object: table
name: user
comment: ""
attributes:
  table: user
  engine: innodb
  charset: utf8mb4
tags:
  - json
  - yaml
fields:
  - name: id
    type: int32
    comment: name of Hello
    attributes:
      autoIncrement: true
  - name: name
    type: string
    comment: name of Hello
    attributes:
      primary: true
      index: true
  - name: Age
    type: int32
    comment: age of Hello
    attributes:
      nullable: true
  - name: Sex
    type: bool
    attributes:
      dbtype: TINYINT
  - name: FooBar
    type: int32
    comment: fooBar
    tags:
      json: foo
      yaml: bar
    attributes:
      nullable: true
  - name: CreateAt
    type: datetime
    comment: create_at
  - name: UpdateAt
    type: datetime
    comment: update_at
    attributes:
      nullable: true
uniques:
  - fields:
      - name
indexes:
  - fields:
      - age
      - sex
primary:
  - id
---
version: v1
object: query
name: summary
tags:
  - json
  - yaml
fields:
  - name: id
    type: int32
    comment: name of Hello
    attributes:
      autoIncrement: true
  - name: name
    type: string
    comment: name of Hello
    attributes:
      primary: true
      index: true
  - name: Age
    type: int32
    comment: age of Hello
    attributes:
      nullable: true
  - name: Sex
    type: bool
    attributes:
      dbtype: TINYINT
```

- 命令生成代码与脚本

```sh
$: orm -h
Usage:
  orm [command]

Available Commands:
  code        code command
  decode      orm decode string
  encode      orm encode string
  help        Help about any command
  password    orm password string
  script      script command

Flags:
  -p, --go-package-name string   go package name (default "model")
  -h, --help                     help for orm
  -i, --input string             input directory (default ".")
  -o, --output string            output directory (default ".")
  -t, --template-suffix string   template suffix (default ".gogo")
      --version                  version for orm
  -c, --workdir string           orm workdir (default: $HOME/.orm)

Use "orm [command] --help" for more information about a command.
```

## 修改模板

安装 `orm` 命令后，完成首次代码生成后，默认代码生成模板文件会安装到 `$HOME/.orm` 文件夹。

用户可以根据自己的需求修改代码生成模板实现个性化配置。

**生成 Golang 代码**

```go
$: orm -i ./examples -o ./examples/model code golang
```

**生成 Mysql 脚本**

```go
$: orm -i ./examples -o ./examples/mysql script mysql
```
