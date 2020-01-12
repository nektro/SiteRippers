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
	flagSaveDir := pflag.String("save-dir", "./data/", "Path to folder to save downloaded data to.")
	flagConcurr := pflag.Int("concurrency", 10, "Maximum number of tasks to run at once. Exactly how tasks are used varies slightly.")
	flagSite := pflag.String("site", "", "Domain of site to rip. None passed means rip all.")
	pflag.Parse()

	//

	doneDir := *flagSaveDir
	doneDir, _ = filepath.Abs(doneDir)
	util.Assert(util.DoesDirectoryExist(doneDir), "--done-dir must point to a valid existing directory!")

	idata.Concurrency = *flagConcurr

	util.RunOnClose(onClose)
	mbpp.Init(*flagConcurr)

	//

	dd := doneDir + "/" + *flagSite
	os.MkdirAll(dd, os.ModePerm)
	idata.Handlers[*flagSite](*flagSite, dd)

	//

	time.Sleep(time.Second / 2)
	mbpp.Wait()
	onClose()
}

func onClose() {
	util.Log(mbpp.GetCompletionMessage())
}
