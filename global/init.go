package global

import (
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var RootDir string
var once = new(sync.Once)

func Init() {
	once.Do(func() {
		inferRootDir()
		initConfig()
	})
}

// inferRootDir 推斷出項目的根目錄
func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string {
		// 要確認項目根目錄下存在 template 目錄
		if exists(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	RootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
