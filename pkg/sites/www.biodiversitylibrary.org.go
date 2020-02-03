package sites

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
)

// www.biodiversitylibrary.org
// https://www.biodiversitylibrary.org

// https://www.biodiversitylibrary.org/item/10
// https://www.biodiversitylibrary.org/itempdf/10

func init() {
	idata.Handlers["www.biodiversitylibrary.org"] = func(site, doneDir string) {

		getMax := func() int {
			doc, err := iutil.FetchDoc("https://www.biodiversitylibrary.org/recent", nil)
			util.DieOnError(err)
			id, _ := doc.Find("section.recentfeed li a.booktitle").First().Attr("href")
			id = strings.Split(id, "/")[2]
			max, _ := strconv.Atoi(id)
			return max
		}

		max := getMax()
		util.Log("max:", max)
		mbpp.CreateJob(site, func(b *mbpp.BarProxy) {
			b.AddToTotal(int64(max))
			for i := 1; i <= max; i++ {
				idata.Guard.Add()
				x := i
				go func() {
					defer b.Increment(1)
					defer idata.Guard.Done()
					//
					n := strconv.Itoa(x)
					iurlS := "https://www.biodiversitylibrary.org/item/" + n
					res, err := http.Head(iurlS)
					if err != nil {
						fmt.Fprintln(idata.Log, 1, n, err)
						return
					}
					if res.StatusCode != http.StatusOK {
						if res.StatusCode == http.StatusNotFound {
							return
						}
						fmt.Fprintln(idata.Log, 2, n, res.StatusCode)
						return
					}
					doc, err := iutil.FetchDoc(iurlS, nil)
					if err != nil {
						fmt.Fprintln(idata.Log, 3, n, err)
						return
					}
					t, _ := doc.Find(`meta[name="citation_title"]`).Attr("content")
					t = strings.ReplaceAll(t, "/", "+")
					f := doneDir + "/" + "[" + n + "]" + " " + t + ".pdf"
					urlS := "https://www.biodiversitylibrary.org/itempdf/" + n
					go mbpp.CreateDownloadJob(urlS, f, nil)
				}()
			}
		})
	}
}
