package sites

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/util"
)

// horriblesubs.info
// https://horriblesubs.info

func init() {
	idata.Handlers["horriblesubs.info"] = func(site, doneDir string) {

		getPage := func(i int, j int) (*goquery.Document, error) {
			fmt.Println("show:", i, ",", "page:", j)
			return iutil.FetchDoc("https://"+site+"/api.php?method=getshows&type=show&showid="+strconv.Itoa(i)+"&nextid="+strconv.Itoa(j), nil)
		}
		getTitle := func(d *goquery.Selection) string {
			t := d.Find(".rls-label").Eq(0).Text()
			if len(t) == 0 {
				return ""
			}
			t = t[9:]
			a := strings.Split(t, " ")
			a = a[:len(a)-2]
			t = strings.Join(a, " ")
			return t
		}
		getMagnet := func(d *goquery.Selection) (string, string) {
			a := d.Find(".rls-link").Last().Find(".hs-magnet-link a")
			if a.Size() == 0 {
				return "", ""
			}
			mag, _ := a.Attr("href")
			ql, _ := a.Parent().Parent().Attr("id")
			ql = strings.Split(ql, "-")[1]
			return ql, mag
		}
		saveFile := func(pathS string, data string) bool {
			if util.DoesFileExist(pathS) {
				return false
			}
			f, _ := os.Create(pathS)
			util.Log("Creating:", pathS)
			defer f.Close()
			fmt.Fprintln(f, data)
			return true
		}
		processDoc := func(d *goquery.Document) (le string) {
			q := false
			l := d.Find(".rls-info-container")
			if l.Size() == 0 {
				return "01"
			}
			l.Each(func(_ int, el *goquery.Selection) {
				t := getTitle(el)
				ep, _ := el.Attr("id")
				ql, mag := getMagnet(el)
				//
				if len(mag) > 0 {
					t = strings.ReplaceAll(t, " - ", ".")
					t = strings.ReplaceAll(t, "/", "+")
					t = strings.ReplaceAll(t, " ", ".")
					dir := doneDir + "/" + t
					os.Mkdir(dir, os.ModePerm)
					p := dir + "/" + t + ".E" + ep + "." + ql + ".magnet.txt"
					// fmt.Println(p, util.TrimLen(mag, 100))
					q = !saveFile(p, mag)
				}
				le = ep
			})
			if q {
				return "01"
			}
			return le
		}

		//
		//

		bc := 0
		for i := 0; true; i++ {
			doc, _ := getPage(i, 0)
			if len(doc.Find(".rls-label").Eq(0).Text()) == 0 {
				util.LogError(doc.Text())
				bc++
				if bc == 10 {
					break
				}
				continue
			}
			bc = 0
			lastEp := processDoc(doc)
			for k := 1; lastEp != "01"; k++ {
				d, _ := getPage(i, k)
				lastEp = processDoc(d)
			}
		}
	}
}
