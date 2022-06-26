package main

import (
	"github.com/mih-kopylov/bulker/cmd"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/afero"
)

func main() {
	utils.ConfigureFS(afero.NewOsFs())
	cmd.Execute()
}
