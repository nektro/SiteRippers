package sites

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"
	"github.com/nektro/go-util/mbpp"
)

// unsplash.com
// https://unsplash.com

func init() {
	idata.Handlers["unsplash.com"] = func(site, doneDir string) {
		mbpp.CreateJob("unsplash.com", func(bar *mbpp.BarProxy) {
			for i := 1; true; i++ {
				n := strconv.Itoa(i)
				val, _ := iutil.FetchJson("https://unsplash.com/napi/photos?page="+n+"&per_page=30", nil)
				arr := val.GetArray()
				if len(arr) == 0 {
					break
				}
				bar.AddToTotal(int64(len(arr)))
				for _, item := range arr {
					ui := string(item.GetStringBytes("user", "id"))
					un := string(item.GetStringBytes("user", "username"))
					dir := doneDir + "/" + ui + " - " + un
					os.MkdirAll(dir, os.ModePerm)
					id := string(item.GetStringBytes("id"))
					urlS := string(item.GetStringBytes("urls", "raw"))
					pathS := dir + "/" + id + ".jpg"
					metaF, _ := os.Create(pathS + ".json")
					fmt.Fprintln(metaF, string(item.MarshalTo([]byte{})))
					go mbpp.CreateDownloadJob(urlS, pathS, bar)
				}
			}
		})
	}
}
