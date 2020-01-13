package iutil

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	req.Header.Add("user-agent", "nektro/SiteRippers")
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

func FetchJson(urlS string, headers map[string]string) (*fastjson.Value, error) {
	res, err := Fetch(urlS, headers)
	if err != nil {
		return nil, err
	}
	bys, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	val, err := fastjson.ParseBytes(bys)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func RemoveAll(s string, subs ...string) string {
	for _, item := range subs {
		s = strings.ReplaceAll(s, item, "")
	}
	return s
}
