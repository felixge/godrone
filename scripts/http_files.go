package main

import (
	"github.com/felixge/makefs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// Get name/dir of this source file
var (
	_, __filename, _, _ = runtime.Caller(0)
	__dirname           = filepath.Dir(__filename)
)

func main() {
	httpDir := filepath.Join(__dirname, "../apis/http")
	file, err := os.OpenFile(filepath.Join(httpDir, "fs.go"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fs := makefs.NewFs(http.Dir(filepath.Join(httpDir, "files")))
	if err := fs.Fprint(file, "http", "files"); err != nil {
		panic(err)
	}
}
