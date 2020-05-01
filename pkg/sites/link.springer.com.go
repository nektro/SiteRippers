package sites

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
)

// link.springer.com
// https://link.springer.com

func init() {
	idata.Handlers["link.springer.com"] = func(site, doneDir string) {
		//
		// books
		mbpp.CreateJob(site+": Books", func(bar *mbpp.BarProxy) {
			for i := 1; true; i++ {
				l := iutil.FetchDocAs("https://link.springer.com/search/page/"+strconv.Itoa(i)+"?showAll=false&facet-content-type=%22Book%22", nil, "#results-list li .text .title", func(a, b int, c string, el *goquery.Selection) {
					if !strings.HasPrefix(c, "/book/") {
						return
					}
					doi := c[6:]
					up := "https://link.springer.com/content/pdf/" + doi + ".pdf"
					ue := "https://link.springer.com/download/epub/" + doi + ".epub"
					title := el.Text()
					authors := strings.ReplaceAll(iutil.RemoveAll(el.Parent().Parent().Find(".meta .authors").Text(), "\n", "    ", "…"), ",", ", ")
					year, _ := el.Parent().Parent().Find(".meta .year").Attr("title")
					fileN := fmt.Sprintf("[%s][%s] %s - %s", year, strings.ReplaceAll(doi, "/", "+"), authors, title)
					//
					dir := doneDir + "/Books/" + year
					if !util.DoesDirectoryExist(dir) {
						os.MkdirAll(dir, os.ModePerm)
					}
					bar.AddToTotal(2)
					go mbpp.CreateDownloadJob(up, dir+"/"+fileN+".pdf", bar)
					go mbpp.CreateDownloadJob(ue, dir+"/"+fileN+".epub", bar)
				})
				fmt.Fprintln(idata.Log, "books", i, l)
				if l == 0 {
					break
				}
			}
		})

		// articles
		mbpp.CreateJob(site+": Articles", func(bar *mbpp.BarProxy) {
			for i := 1; true; i++ {
				l := iutil.FetchDocAs("https://link.springer.com/search/page/"+strconv.Itoa(i)+"?showAll=false&facet-content-type=%22Article%22", nil, "#results-list li h2 .title", func(a, b int, c string, el *goquery.Selection) {
					if !strings.HasPrefix(c, "/article/") {
						return
					}
					doi := c[9:]
					authors := strings.ReplaceAll(iutil.RemoveAll(el.Parent().Parent().Find(".meta .authors").Text(), "\n", "    ", "…"), ",", ", ")
					title := el.Text()
					year, _ := el.Parent().Parent().Find(".meta .year").Attr("title")
					year = year[strings.Index(year, " ")+1:]
					fileN := fmt.Sprintf("[%s][%s] %s - %s", year, strings.ReplaceAll(doi, "/", "+"), authors, title)
					//
					dir := doneDir + "/Articles/" + year
					if !util.DoesDirectoryExist(dir) {
						os.MkdirAll(dir, os.ModePerm)
					}
					urlS := "https://link.springer.com/content/pdf/" + doi + ".pdf"
					bar.AddToTotal(1)
					go mbpp.CreateDownloadJob(urlS, dir+"/"+fileN+".pdf", bar)
				})
				fmt.Fprintln(idata.Log, "articles", i, l)
				if l <= 0 {
					break
				}
			}
		})
	}
}
