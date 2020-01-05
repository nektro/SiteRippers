package sites

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
)

func init() {
	idata.Handlers["rule34.xxx"] = func(site, doneDir string) {

		grabPost := func(id string, b *mbpp.BarProxy, w *sync.WaitGroup) string {
			defer b.Increment(1)

			req, _ := http.NewRequest(http.MethodGet, "https://"+site+"/index.php?page=post&s=view&id="+id, nil)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				w.Done()
				return ""
			}
			if res.StatusCode != http.StatusOK {
				w.Done()
				return ""
			}
			doc, _ := goquery.NewDocumentFromResponse(res)
			urlS, _ := doc.Find(`div.sidebar a[href^="https"]`).Eq(0).Attr("href")
			pth := iutil.GetPathFile(urlS)
			if len(pth) == 0 {
				w.Done()
				return ""
			}
			p := doneDir + "/" + id + "_" + pth
			go func() {
				mbpp.CreateDownloadJob(urlS, p, nil)
				w.Done()
			}()
			return p
		}

		//
		//

		req, _ := http.NewRequest(http.MethodGet, "https://"+site+"/index.php?page=post&s=list", nil)
		res, _ := http.DefaultClient.Do(req)
		doc, _ := goquery.NewDocumentFromResponse(res)
		c, _ := doc.Find("div span.thumb").Eq(0).Attr("id")

		max, _ := strconv.Atoi(c[1:])
		fmt.Println(max)

		wg := new(sync.WaitGroup)
		paths := []string{}

		mbpp.CreateJob(site, func(bar *mbpp.BarProxy) {
			bar.AddToTotal(int64(max))
			for i := 0; i <= max; i++ {
				if i%1000 == 0 {
					if util.DoesFileExist(doneDir + "/" + strconv.Itoa(i) + ".tar") {
						i += 999
						bar.Increment(1000)
						continue
					}
				}
				if i%1000 == 0 && i != 0 {
					wg.Wait()
					iutil.CreateTarball(doneDir+"/"+strconv.Itoa(i-1000)+".tar", site, paths)
					for _, item := range paths {
						os.Remove(item)
					}
					paths = []string{}
				}
				wg.Add(1)
				pt := grabPost(strconv.Itoa(i), bar, wg)
				if len(pt) > 0 {
					paths = append(paths, pt)
				}
			}
		})
	}
}
