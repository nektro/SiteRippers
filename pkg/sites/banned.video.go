package sites

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nektro/SiteRippers/pkg/idata"
	"github.com/nektro/SiteRippers/pkg/iutil"

	"github.com/nektro/go-util/mbpp"
	"github.com/valyala/fastjson"
)

// banned.video
// https://banned.video

func init() {
	idata.Handlers["banned.video"] = func(site, doneDir string) {

		type tB struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}
		type tA struct {
			OperationName string `json:"operationName"`
			Query         string `json:"query"`
			Variables     tB     `json:"variables"`
		}

		dat := tA{
			"GetNewVideos",
			`query GetNewVideos($limit: Float, $offset: Float) {
				getNewVideos(limit: $limit, offset: $offset) {
					_id
					title
					summary
					playCount
					largeImage
					embedUrl
					published
					videoDuration
					channel {
						_id
						title
						avatar
						__typename
					}
					createdAt
					__typename
				}
			}`,
			tB{
				24,
				0,
			},
		}

		for true {
			datt := string(iutil.MarshalJSON(dat))
			datt = strings.ReplaceAll(datt, "\\t", "")
			req, _ := http.NewRequest(http.MethodPost, "https://api.infowarsmedia.com/graphql", strings.NewReader(datt))
			req.Header.Add("content-type", "application/json")
			req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64; rv:68.0) Gecko/20100101 Firefox/68.0")
			res, _ := http.DefaultClient.Do(req)
			bys, _ := ioutil.ReadAll(res.Body)
			val, _ := fastjson.Parse(string(bys))

			arr := val.GetArray("data", "getNewVideos")
			if len(arr) == 0 {
				break
			}
			for _, item := range arr {
				uploader := string(item.GetStringBytes("channel", "title"))
				uploadDate, _ := time.Parse(time.RFC3339, string(item.GetStringBytes("createdAt")))
				uploadDateS := uploadDate.UTC().Format("2006.01.02")
				title := strings.Trim(string(item.GetStringBytes("title")), " ")
				duration := strconv.Itoa(int(item.GetFloat64("videoDuration")))

				dir := doneDir + "/" + uploader
				os.MkdirAll(dir, os.ModePerm)
				prefix := uploadDateS + " - " + title + " - [" + duration + "s]"
				pathS := dir + "/" + prefix

				ioutil.WriteFile(pathS+".json", append(item.MarshalTo([]byte{}), '\n'), os.ModePerm)

				go iutil.DownloadToPrefix(string(item.GetStringBytes("channel", "avatar")), dir, "avatar_", nil)
				go iutil.DownloadToPrefix(string(item.GetStringBytes("largeImage")), dir, prefix+"_", nil)

				doc, _ := iutil.FetchDoc(string(item.GetStringBytes("embedUrl")), nil)
				urlS, _ := doc.Find("div#embedData").Attr("downloadurl")
				go mbpp.CreateDownloadJob(urlS, dir+"/"+prefix+".mp4", nil)
			}
			dat.Variables.Offset += 24
		}
	}
}
