package core

import (
	"os"
	"path/filepath"
)

type Yaml struct {
	Replies []Reply
}

var ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

var Config Yaml
