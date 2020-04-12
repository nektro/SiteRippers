package iutil

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
	"github.com/valyala/fastjson"
)

func CreateTarball(tarballFilePath string, topDir string, filePaths []string) error {
	file, _ := os.Create(tarballFilePath)
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	mbpp.CreateJob("Packing: "+tarballFilePath, func(bar *mbpp.BarProxy) {
		bar.AddToTotal(int64(len(filePaths)))
		for _, filePath := range filePaths {
			addFileToTarWriter(topDir, filePath, tarWriter)
			bar.Increment(1)
		}
	})

	return nil
}

func addFileToTarWriter(topDir string, filePath string, tarWriter *tar.Writer) {
	file2, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file2.Close()

	stat, _ := file2.Stat()
	header := &tar.Header{
		Name:    "/" + topDir + "/" + GetPathFile(filePath),
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}
	tarWriter.WriteHeader(header)
	io.Copy(tarWriter, file2)
	return
}

func GetPathFile(p string) string {
	q := strings.Split(p, "/")
	return q[len(q)-1]
}

func GetUrlPathFile(urlS string) string {
	urlO, _ := url.Parse(urlS)
	return GetPathFile(urlO.Path)
}

func Fetch(urlS string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, urlS, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36")
	req.Header.Add("connection", "close")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	return res, nil
}

func FetchDoc(urlS string, headers map[string]string) (*goquery.Document, error) {
	res, err := Fetch(urlS, headers)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func FetchBin(urlS string, headers map[string]string) ([]byte, error) {
	res, err := Fetch(urlS, headers)
	if err != nil {
		return nil, err
	}
	bys, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bys, nil
}

func FetchJson(urlS string, headers map[string]string) (*fastjson.Value, error) {
	bys, err := FetchBin(urlS, headers)
	if err != nil {
		return nil, err
	}
	val, err := fastjson.ParseBytes(bys)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// FetchDocAs means Fetch Doc A's
// params are total, index, string
// returns if total > 0
func FetchDocAs(urlS string, headers map[string]string, sel string, f func(int, int, string, *goquery.Selection)) bool {
	doc, err := FetchDoc(urlS, headers)
	if err != nil {
		return false
	}
	arr := doc.Find(sel)
	n := arr.Length()
	arr.Each(func(i int, el *goquery.Selection) {
		hf, ok := el.Attr("href")
		if !ok {
			return
		}
		f(n, i, hf, el)
	})
	return arr.Length() > 0
}

func RemoveAll(s string, subs ...string) string {
	for _, item := range subs {
		s = strings.ReplaceAll(s, item, "")
	}
	return s
}

func PadLeft(s string, leng int, pre string) string {
	if len(s) >= leng {
		return s
	}
	return PadLeft(pre+s, leng, pre)
}

func DownloadTo(urlS string, dir string, bar *mbpp.BarProxy) {
	mbpp.CreateDownloadJob(urlS, dir+"/"+GetUrlPathFile(urlS), bar)
}

// from SO
func SplitLen(str string, size int) []string {
	strLength := len(str)
	splitedLength := int(math.Ceil(float64(strLength) / float64(size)))
	splited := make([]string, splitedLength)
	var start, stop int
	for i := 0; i < splitedLength; i++ {
		start = i * size
		stop = start + size
		if stop > strLength {
			stop = strLength
		}
		splited[i] = str[start:stop]
	}
	return splited
}

func UrlExt(urlS string) string {
	urlO, _ := url.Parse(urlS)
	return filepath.Ext(urlO.Path)
}

func MakeCounter(dirIn string, after int) *Counter {
	return NewCounter(after, func() {
		go PackCbzArchive(dirIn)
	})
}

func PackCbzArchive(dirIn string) {
	mbpp.CreateJob("Packing: "+dirIn, func(b *mbpp.BarProxy) {
		outf, _ := os.Create(dirIn + ".cbz")
		defer outf.Close()
		outz := zip.NewWriter(outf)
		defer outz.Close()
		files, _ := ioutil.ReadDir(dirIn)
		b.AddToTotal(int64(len(files) + 1))
		for _, item := range files {
			zw, _ := outz.Create(item.Name())
			bs, _ := ioutil.ReadFile(dirIn + "/" + item.Name())
			zw.Write(bs)
			b.Increment(1)
		}
		os.RemoveAll(dirIn)
		b.Increment(1)
	})
}

func MarshalJSON(v interface{}) []byte {
	res, _ := json.Marshal(&v)
	return res
}

func DownloadToPrefix(urlS string, dir string, prefix string, bar *mbpp.BarProxy) {
	mbpp.CreateDownloadJob(urlS, dir+"/"+prefix+GetUrlPathFile(urlS), bar)
}
