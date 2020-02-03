package sites

import (
	"os"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
)

// ipsw.me
// https://ipsw.me

func init() {
	idata.Handlers["ipsw.me"] = func(site, doneDir string) {

		v, err := iutil.FetchJson("https://api.ipsw.me/v4/devices", nil)
		util.DieOnError(err)
		mbpp.CreateJob("ipsw.me", func(bar *mbpp.BarProxy) {
			l := v.GetArray()
			for _, item := range l {
				iden := string(item.GetStringBytes("identifier"))
				dir := doneDir + "/" + iden
				os.Mkdir(dir, os.ModePerm)
				w, err := iutil.FetchJson("https://api.ipsw.me/v4/device/"+iden+"?type=ipsw", nil)
				util.DieOnError(err)
				m := w.GetArray("firmwares")
				bar.AddToTotal(int64(len(m)))
				for _, jtem := range m {
					urlS := string(jtem.GetStringBytes("url"))
					go iutil.DownloadTo(urlS, dir, bar)
				}
			}
		})
	}
}
