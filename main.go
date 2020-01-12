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
	flagSD := pflag.String("save-dir", "./data/", "Path to folder to save downloaded data to.")
	flagCC := pflag.Int("concurrency", 10, "Maximum number of tasks to run at once. Exactly how tasks are used varies slightly.")
	flagSN := pflag.String("site", "", "Domain of site to rip. None passed means rip all.")
	pflag.Parse()

	//

	doneDir := *flagSD
	doneDir, _ = filepath.Abs(doneDir)
	util.Assert(util.DoesDirectoryExist(doneDir), "--done-dir must point to a valid existing directory!")

	idata.Concurrency = *flagCC

	util.RunOnClose(onClose)
	mbpp.Init(*flagCC)

	//

	dd := doneDir + "/" + *flagSN
	os.MkdirAll(dd, os.ModePerm)
	idata.Handlers[*flagSN](*flagSN, dd)

	//

	time.Sleep(time.Second / 2)
	mbpp.Wait()
	onClose()
}

func onClose() {
	util.Log(mbpp.GetCompletionMessage())
}
