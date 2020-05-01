package sites

import (
	"os"
	"strconv"
	"strings"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
)

// punkrockvids.com
// https://punkrockvids.com

// https://punkrockvids.com/?mod=videos&subid=official&vidid=&page=400

func init() {
	idata.Handlers["punkrockvids.com"] = func(site, doneDir string) {
		//
		grabCategory := func(cat string) {
			mbpp.CreateJob(site+": "+cat, func(bar *mbpp.BarProxy) {
				dir := doneDir + "/" + cat
				os.MkdirAll(dir, os.ModePerm)
				fv := ""
				for i := 1; true; i++ {
					stop := false
					iutil.FetchDocAs("https://punkrockvids.com/?mod=videos&subid="+cat+"&page="+strconv.Itoa(i), nil, "body table table td:nth-child(3) a:nth-of-type(2)", func(a, b int, c string, el *goquery.Selection) {
						if stop {
							return
						}
						if b == 0 {
							if fv == c {
								stop = true
								return
							}
							fv = c
						}
						if !strings.HasPrefix(c, "videos/") {
							return
						}
						bar.AddToTotal(1)
						go func() {
							iutil.DownloadTo("https://"+site+"/"+c, dir, bar)
						}()
					})
					if stop {
						break
					}
				}
			})
		}

		//

		grabCategory("official")
		grabCategory("live")
		grabCategory("ivideos")
		// grabCategory("archive")
		// grabCategory("interviews")
	}
}
