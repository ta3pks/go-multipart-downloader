# go-multipart-downloader
A multipart downloader for teaching purposes

## USAGE

```go
go run main.go -u <url>
```
## SUPPORTED FLAGS

`-u,--url` url to download from

`-o,--outfile` output filename

`-t,--threads` number of parallel threads 
`-h` help message

If output filename does not exits the program will try to determine the name from the URL

# TODO

- [ ] progress bar
- [ ] stats
