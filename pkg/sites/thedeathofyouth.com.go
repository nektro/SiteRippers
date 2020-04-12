package sites

import (
	"strconv"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/PuerkitoBio/goquery"
)

// thedeathofyouth.com
// http://thedeathofyouth.com

func init() {
	idata.Handlers["thedeathofyouth.com"] = func(site, doneDir string) {

		for i := 1; true; i++ {
			n := strconv.Itoa(i)
			doc, err := iutil.FetchDoc("http://"+site+"/page/"+n+"/", nil)
			if err != nil {
				break
			}
			doc.Find("div.post a img").Each(func(_ int, el *goquery.Selection) {
				urlS, _ := el.Attr("src")
				if len(urlS) == 0 {
					return
				}
				urlS = urlS[:len(urlS)-4-7-1] + ".jpg"
				go iutil.DownloadTo(urlS, doneDir, nil)
			})
		}
	}
}
