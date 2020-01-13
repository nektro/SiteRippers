package sites

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
	"golang.org/x/sync/semaphore"
)

// https://thepiratebay.org

func init() {
	idata.Handlers["thepiratebay.org"] = func(site, doneDir string) {

		saveMagnetToFile := func(pathS string, urlO *url.URL, nfo string) {
			if util.DoesFileExist(pathS) {
				return
			}
			f, _ := os.Create(pathS)
			fmt.Fprintln(f, urlO.String())
			fmt.Fprintln(f, string(bytes.Repeat([]byte("-"), len(urlO.String()))))
			fmt.Fprintln(f, nfo)
		}

		doc, err := iutil.FetchDoc("https://"+site+"/recent", nil)
		util.DieOnError(err)
		mRec, _ := doc.Find("tr td div.detName a.detLink").Eq(0).Attr("href")
		max, err := strconv.Atoi(strings.Split(mRec, "/")[2])
		util.DieOnError(err)

		guard := semaphore.NewWeighted(int64(idata.Concurrency))
		ctx := context.Background()

		mbpp.CreateJob(site, func(bar *mbpp.BarProxy) {
			bar.AddToTotal(int64(max))
			for i := max; i >= 1; i-- {
				guard.Acquire(ctx, 1)
				j := i
				js := strconv.Itoa(j)
				go mbpp.CreateJob("/torrent/"+js, func(*mbpp.BarProxy) {
					defer bar.Increment(1)
					defer guard.Release(1)
					doc2, err := iutil.FetchDoc("https://"+site+"/torrent/"+js+"/", nil)
					if err != nil {
						return
					}
					//
					magS, _ := doc2.Find(".download a:nth-child(1)").Eq(0).Attr("href")
					urlO, _ := url.Parse(magS)
					nfo := doc2.Find(".nfo").Children().Eq(0).Text()
					p := doneDir + "/" + fmt.Sprintf("[%d] %s.text", i, urlO.Query().Get("dn"))
					saveMagnetToFile(p, urlO, nfo)
				})
			}
		})
	}
}
