# Simple Bank Manager Setup

Make sure `go` and `sqlite` are installed. Run the following command in the terminal:

1. `go mod tidy`
2. `go run ./main.go`

**Note:** If you get the following error on step 2 in Linux, then execute `sudo apt-get install build-essential libsqlite3-dev` in the terminal and re-execute step 2:

```
[error] failed to initialize database, got error Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work
```
