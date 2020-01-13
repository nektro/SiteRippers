package sites

import (
	"os"
	"strconv"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/nektro/go-util/mbpp"
)

// 4chan.org

func init() {
	idata.Handlers["4chan.org"] = func(site, doneDir string) {

		grabThread := func(board, id string) {
			dir := doneDir + "/" + board + "/" + id
			m := false
			//
			mbpp.CreateJob("/"+board+"/"+id+"/", func(bar *mbpp.BarProxy) {
				val, _ := iutil.FetchJson("https://p.4chan.org/4chan/board/"+board+"/thread/"+id, nil)
				//
				ar := val.GetArray("body", "posts")
				bar.AddToTotal(int64(len(ar)))
				for _, item := range ar {
					t := strconv.Itoa(item.GetInt("tim"))
					f := string(item.GetStringBytes("filename"))
					e := string(item.GetStringBytes("ext"))
					u := "https://i.4cdn.org/" + board + "/" + t + e
					//
					if len(e) == 0 {
						bar.Increment(1)
						continue
					}
					if !m {
						os.MkdirAll(dir, os.ModePerm)
						m = true
					}
					//
					go mbpp.CreateDownloadJob(u, dir+"/"+t+"_"+f+e, bar)
				}
				bar.Wait()
			})
		}

		grabBoard := func(board string) {
			mbpp.CreateJob("/"+board+"/", func(bar *mbpp.BarProxy) {
				val, _ := iutil.FetchJson("https://p.4chan.org/4chan/board/"+board+"/catalog", nil)
				//
				ar1 := val.GetArray("body")
				ids := []string{}
				for _, item := range ar1 {
					ar2 := item.GetArray("threads")
					for _, jtem := range ar2 {
						ids = append(ids, strconv.Itoa(jtem.GetInt("no")))
					}
				}
				bar.AddToTotal(int64(len(ids)))
				for _, item := range ids {
					grabThread(board, item)
					bar.Increment(1)
				}
			})
		}

		mbpp.CreateJob(site, func(bar *mbpp.BarProxy) {
			val, _ := iutil.FetchJson("https://p.4chan.org/4chan/boards", nil)
			//
			ar := val.GetArray("body", "boards")
			bar.AddToTotal(int64(len(ar)))
			for _, item := range ar {
				grabBoard(string(item.GetStringBytes("board")))
				bar.Increment(1)
			}
		})
	}
}
