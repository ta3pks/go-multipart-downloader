# go-multipart-downloader
A multipart downloader for teaching purposes

## USAGE

```go
go run main.go -u <url>
go run main.go -p -u <url>
go run main.go --progress=true -u <url>
```
## SUPPORTED FLAGS

`-u,--url` url to download from

`-o,--outfile` output filename

`-t,--threads` number of parallel threads 

`-h` help message

`-p`its progress bar(its default)

If output filename does not exits the program will try to determine the name from the URL
If you won't work with progressbar you can use "go run main.go --progress=false -u <url>" this command without progressbar
# TODO

- [x] progress bar
- [ ] stats
