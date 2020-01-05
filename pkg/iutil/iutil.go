package iutil

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/nektro/go-util/mbpp"
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
