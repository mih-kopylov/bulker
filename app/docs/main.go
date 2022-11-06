package main

import (
	"github.com/mih-kopylov/bulker/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra/doc"
	"os"
	"path/filepath"
)

func main() {
	rootCommand := cmd.CreateRootCommand("")
	rootCommand.DisableAutoGenTag = true
	dir, err := filepath.Abs("./dist/docs")
	if err != nil {
		logrus.Fatal(err)
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logrus.Fatal(err)
	}

	err = doc.GenMarkdownTree(rootCommand, dir)
	if err != nil {
		logrus.Fatal(err)
	}
}
