package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"net/http"

	"github.com/spf13/pflag"
)

var file = pflag.StringP("url", "u", "", "url to download")
var outname = pflag.StringP("outfile", "o", "", "output filename")

const CHUNK_SIZE = 10000

type hdrs = map[string][]string

var wg sync.WaitGroup

func main() {
	pflag.Parse()
	if *file == "" {
		log.Fatalln("no file")
	}
	var headers = ReadHeaders(*file)
	Download(headers, *file)
}
func ReadHeaders(url string) map[string][]string { //{{{

	resp, err := http.Head(url)
	if nil != err {
		log.Fatalln("cannot continue: ", err)
	}
	return resp.Header
} //}}}
func IsRangeSupported(headers map[string][]string) bool { //{{{
	_, ok := headers["Accept-Ranges"]
	return ok
} //}}}
func Download(h hdrs, url string) { //{{{
	if !IsRangeSupported(h) {
		//download single part
		DownloadSinglepart(url)
		return
	}
	fmt.Println("multipart")
	//download multipart
	DownloadMultipart(h, url)
} //}}}
func DownloadSinglepart(url string) { //{{{
	resp, err := http.Get(url)
	if nil != err {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	f, err := os.Create(GetFilename(url))
	if nil != err {
		log.Fatalln(err)
	}
	defer f.Close()
	io.Copy(f, resp.Body)
} //}}}
func GetFilename(s string) string { //{{{
	if *outname != "" {
		return *outname
	}
	_filename := strings.Split(s, "/")
	return _filename[len(_filename)-1]
} //}}}
func DownloadMultipart(h hdrs, url string) { //{{{
	_ln, ok := h["Content-Length"]
	if !ok && len(_ln) > 0 {
		fmt.Println("fallback to single part")
		DownloadSinglepart(url)
	}
	ln, _ := strconv.Atoi(_ln[0])
	var lower_bound = 0
	var last bool
	var part int
	for {
		if last {
			break
		}
		lower_bound = ln - CHUNK_SIZE
		if lower_bound <= 0 {
			last = true
			lower_bound = 0
		}
		part_str := fmt.Sprintf("%d-%d", lower_bound, ln)
		fmt.Println("downloading", part_str, "part:", part)
		wg.Add(1)
		go download_part(url, part_str, part)
		ln = ln - CHUNK_SIZE - 1
		part++
	}
	wg.Wait()
	glue_parts(GetFilename(url), part)
	fmt.Println("finished")
} //}}}
func download_part(url, rng string, part int) { //{{{
	defer wg.Done()
	file, err := os.Create(fmt.Sprintf("%d.%s.part", part, GetFilename(url)))
	if nil != err {
		log.Fatalln(err)
	}
	defer file.Close()
	req, err := http.NewRequest("GET", url, nil)
	if nil != err {
		log.Fatalln(err)
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%s", rng))
	res, err := http.DefaultClient.Do(req)
	if nil != err {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	io.Copy(file, res.Body)
	fmt.Println("finished", rng, "part:", part)
} //}}}
func glue_parts(fname string, last_part int) { //{{{
	file, err := os.Create(fname)
	if nil != err {
		log.Fatalln(err)
	}
	defer file.Close()
	var partname string
	for last_part >= 0 {
		partname = fmt.Sprintf("%d.%s.part", last_part, fname)
		fmt.Println("gluing", partname)
		f, _ := os.Open(partname)
		defer f.Close()
		io.Copy(file, f)
		last_part--
	}
} //}}}
