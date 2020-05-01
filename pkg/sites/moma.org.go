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

// www.moma.org
// https://www.moma.org/collection/

func init() {
	idata.Handlers["www.moma.org"] = func(site, doneDir string) {

		for i := 1; true; i++ {
			doc, _ := iutil.FetchDoc("https://"+site+"/collection/works?with_images=1&page="+strconv.Itoa(i), map[string]string{"X-Requested-With": "XMLHttpRequest"})
			l := doc.Find("ul.grid li")
			if l.Size() == 0 {
				return
			}
			mbpp.CreateJob("moma.org - page: "+strconv.Itoa(i), func(bA *mbpp.BarProxy) {
				bA.AddToTotal(int64(l.Size()))
				l.Each(func(_ int, el *goquery.Selection) {
					defer bA.Increment(1)
					//
					ah, _ := el.Find("a").Attr("href")
					id := strings.Split(strings.Split(ah, "/")[3], "?")[0]
					doc2, err := iutil.FetchDoc("https://"+site+"/collection/works/"+id, nil)
					if err != nil {
						return
					}
					//
					cap := doc2.Find(`div.work__short-caption h1 span`)
					artist := strings.ReplaceAll(iutil.RemoveAll(cap.Eq(0).Text(), "\t", "\n"), "/", "+")
					title := strings.ReplaceAll(iutil.RemoveAll(cap.Eq(1).Text(), "\t", "\n"), "/", "+")
					year := strings.ReplaceAll(iutil.RemoveAll(cap.Eq(2).Text(), "\t", "\n"), "/", "+")
					dir := doneDir + "/" + artist + "/" + year + "/" + title
					os.MkdirAll(dir, os.ModePerm)
					//
					descFP := dir + "/description.txt"
					if !util.DoesFileExist(descFP) {
						dF, _ := os.Create(descFP)
						fmt.Fprintln(dF, artist)
						fmt.Fprintln(dF, title)
						fmt.Fprintln(dF, year)
						fmt.Fprintln(dF)
						fmt.Fprintln(dF, iutil.RemoveAll(doc2.Find("div.uneven-columns__second blockquote").Text(), "\t", "\n", "            "))
						fmt.Fprintln(dF, iutil.RemoveAll(doc2.Find("div.uneven-columns__second p").Text(), "\t", "\n"))
						doc2.Find("div.uneven-columns__second span").Each(func(j int, el *goquery.Selection) {
							if j%2 == 0 {
								fmt.Fprintln(dF)
							}
							fmt.Fprintln(dF, iutil.RemoveAll(el.Text(), "    ", "\n", "  "))
						})
					}
					//
					doc2.Find(`meta[property="og:image"]`).Each(func(j int, el *goquery.Selection) {
						urlS, _ := el.Attr("content")
						urlS = strings.ReplaceAll(urlS, "http:", "https:")
						go mbpp.CreateDownloadJob(urlS, dir+"/"+strconv.Itoa(j)+"."+iutil.GetUrlPathFile(urlS), nil)
					})
					doc2.Find(`section aside[data-audio-url]`).Each(func(j int, el *goquery.Selection) {
						c, _ := el.Attr("data-audio-url")
						urlS := "https://" + site + c
						go mbpp.CreateDownloadJob(urlS, dir+"/"+strconv.Itoa(j)+"."+iutil.GetUrlPathFile(urlS), nil)
					})
				})
			})
		}
	}
}
