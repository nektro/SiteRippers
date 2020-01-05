package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/nektro/SiteRippers/pkg/idata"

	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
	"github.com/spf13/pflag"

	_ "github.com/nektro/SiteRippers/pkg/sites"
)

func main() {
	flagSaveDir := pflag.String("save-dir", "", "")
	flagConcurr := pflag.Int("concurrency", 10, "")
	flagSites := pflag.StringArray("site", []string{}, "")
	pflag.Parse()

	//

	doneDir := "./data/"
	if len(*flagSaveDir) > 0 {
		doneDir = *flagSaveDir
	}
	doneDir, _ = filepath.Abs(doneDir)
	util.Assert(util.DoesDirectoryExist(doneDir), "--done-dir must point to a valid directory!")

	util.RunOnClose(onClose)
	mbpp.Init(*flagConcurr)

	//

	for _, item := range *flagSites {
		dd := doneDir + "/" + item
		os.MkdirAll(dd, os.ModePerm)
		idata.Handlers[item](item, dd)
	}

	//

	time.Sleep(time.Second / 2)
	mbpp.Wait()
	onClose()
}

func onClose() {
	util.Log(mbpp.GetCompletionMessage())
}
