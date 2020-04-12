package sites

import (
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
)

// www.webtoons.com
// https://www.webtoons.com

func init() {
	idata.Handlers["www.webtoons.com"] = func(site, doneDir string) {

		getCanonicalPg := func(d *goquery.Document) string {
			h, _ := d.Find(`link[rel="canonical"]`).Attr("href")
			u, _ := url.Parse(h)
			p := u.Query().Get("page")
			return p
		}
		mbpp.CreateJob(site, func(bar1 *mbpp.BarProxy) {
			t := 0
			iutil.FetchDocAs("https://"+site+"/en/dailySchedule", nil, ".daily_lst ul li a", func(total, _ int, hf string, el1 *goquery.Selection) {
				if t == 0 {
					bar1.AddToTotal(int64(total))
					t = total
				}
				author := el1.Find(".info .author").Text()
				author = strings.ReplaceAll(author, "/", "∕")
				title := el1.Find(".info .subj").Text()
				title = strings.ReplaceAll(title, "/", "∕")
				//
				mbpp.CreateJob(author+" - "+title, func(bar2 *mbpp.BarProxy) {
					defer bar1.Increment(1)
					b := false
					for i := 1; true; i++ {
						if b {
							break
						}
						n := strconv.Itoa(i)
						bar2.AddToTotal(1)
						//
						mbpp.CreateJob(author+" - "+title+", pg:"+n, func(bar3 *mbpp.BarProxy) {
							defer bar2.Increment(1)
							doc, _ := iutil.FetchDoc(hf+"&page="+n, nil)
							if n != getCanonicalPg(doc) {
								b = true
								return
							}
							dir1 := doneDir + "/" + author + "/" + title
							os.MkdirAll(dir1, os.ModePerm)
							//
							arr := doc.Find(".detail_lst ul li[data-episode-no] a")
							bar3.AddToTotal(int64(arr.Length()))
							arr.Each(func(_ int, el *goquery.Selection) {
								etitle := el.Find(".subj span").Text()
								etitle = strings.ReplaceAll(etitle, "/", "∕")
								eiss, _ := el.Parent().Attr("data-episode-no")
								dir2 := dir1 + "/" + "[" + eiss + "]" + " " + etitle
								//
								mbpp.CreateJob(eiss+"-"+etitle, func(bar4 *mbpp.BarProxy) {
									defer bar3.Increment(1)
									if util.DoesFileExist(dir2 + ".cbz") {
										return
									}
									os.MkdirAll(dir2, os.ModePerm)
									//
									hf2, _ := el.Attr("href")
									doc2, _ := iutil.FetchDoc(hf2, nil)
									lst := doc2.Find("#_imageList img[data-url]")
									bar4.AddToTotal(int64(lst.Length()))
									c := iutil.MakeCounter(dir2, lst.Length())
									lst.Each(func(i int, el *goquery.Selection) {
										urlS, _ := el.Attr("data-url")
										j := iutil.PadLeft(strconv.Itoa(i+1), 3, "0")
										e := iutil.UrlExt(urlS)
										go func() {
											mbpp.CreateDownloadJobOps(urlS, dir2+"/"+j+e, bar4, &mbpp.DownloadJobOptions{
												Headers: map[string]string{
													"Referer": hf2,
												},
											})
											c.Increment()
										}()
									})
								}) // bar 4
							})
						}) // bar 3
					}
				}) // bar 2
			})
		}) // bar 1
	}
}
