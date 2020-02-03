package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nektro/SiteRippers/pkg/idata"

	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/types"
	"github.com/nektro/go-util/util"
	"github.com/spf13/pflag"

	_ "github.com/nektro/SiteRippers/pkg/sites"
)

func main() {
	flagSD := pflag.String("save-dir", "./data/", "Path to folder to save downloaded data to.")
	flagCC := pflag.Int("concurrency", 10, "Maximum number of tasks to run at once. Exactly how tasks are used varies slightly.")
	flagSN := pflag.String("site", "", "Domain of site to rip. None passed means rip all.")
	flagLS := pflag.Bool("list", false, "Pass this to list all supported domains.")
	flagBC := pflag.Int("job-workers", 5, "Maximum number of tasks to initialize in parallel the the background.")
	pflag.Parse()

	//

	if *flagLS {
		for k := range idata.Handlers {
			fmt.Println(k)
		}
		os.Exit(0)
	}

	//

	doneDir := *flagSD
	doneDir, _ = filepath.Abs(doneDir)
	util.Assert(util.DoesDirectoryExist(doneDir), "--done-dir must point to a valid existing directory!")

	idata.Concurrency = *flagCC
	idata.Guard = types.NewSemaphore(*flagBC)

	util.RunOnClose(onClose)
	mbpp.Init(*flagCC)

	//

	dd := doneDir + "/" + *flagSN
	fn, ok := idata.Handlers[*flagSN]
	util.DieOnError(util.Assert(ok, "SiteRipper does not support that domain!"))
	os.MkdirAll(dd, os.ModePerm)
	fn(*flagSN, dd)

	//

	mbpp.Wait()
	onClose()
}

func onClose() {
	util.Log(mbpp.GetCompletionMessage())
}
