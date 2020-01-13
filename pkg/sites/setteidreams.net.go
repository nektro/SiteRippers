package sites

import (
	"fmt"
	"os"
	"strings"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
)

// http://setteidreams.net/
// http://setteidreams.net/artbooks/
// http://setteidreams.net/settei/

// need to fix bug in mbpp, may need to be run multiple times to obtain all data

func init() {
	idata.Handlers["setteidreams.net"] = func(site, doneDir string) {
		// grab artbooks
		docA, _ := iutil.FetchDoc("http://"+site+"/artbooks/", nil)
		os.MkdirAll(doneDir+"/artbooks", os.ModePerm)

		lA := docA.Find("table#dtBasicExample tr td a")
		lA.Each(func(i int, el *goquery.Selection) {
			d := el.Parent().Parent().Children()
			a := d.Eq(0).Children().Eq(0)
			name := strings.ReplaceAll(a.Text(), ":", "")
			year := d.Eq(2).Text()
			series := d.Eq(1).Text()
			//
			lnk, _ := a.Attr("href")
			doc2, _ := iutil.FetchDoc("http://"+site+lnk+"/", nil)
			rar, _ := doc2.Find(".content-wrap").Children().Eq(1).Children().Eq(0).Attr("href")
			//
			title := strings.ReplaceAll(fmt.Sprintf("[Artbook] %s (%s)(%s).cbr", name, year, series), "/", "+")
			go mbpp.CreateDownloadJob(rar, doneDir+"/artbooks/"+title, nil)
		})

		// grab settei
		docS, _ := iutil.FetchDoc("http://"+site+"/settei/", nil)
		os.MkdirAll(doneDir+"/settei", os.ModePerm)

		lS := docS.Find("table#dtBasicExample tr td a")
		lS.Each(func(i int, el *goquery.Selection) {
			d := el.Parent().Parent().Children()
			a := d.Eq(0).Children().Eq(0)
			name := strings.ReplaceAll(a.Text(), ":", "")
			year := d.Eq(1).Text()
			studio := d.Eq(4).Text()
			//
			lnk, _ := a.Attr("href")
			doc2, _ := iutil.FetchDoc("http://"+site+lnk+"/", nil)
			rar, _ := doc2.Find(".content-wrap").Children().Eq(1).Children().Eq(0).Attr("href")
			//
			title := strings.ReplaceAll(fmt.Sprintf("[Settei] %s (%s)(%s).cbr", name, year, studio), "/", "+")
			go mbpp.CreateDownloadJob(rar, doneDir+"/settei/"+title, nil)
		})
	}
}
