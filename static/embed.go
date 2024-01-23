package static

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed all:files
var staticFilesFS embed.FS

func StaticFilesFS() http.FileSystem {
	files, err := fs.Sub(staticFilesFS, "files")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(files)
}

func GetFilePath(path string) string {
	returnPath := fmt.Sprintf("/static/%s", path)

	file, err := staticFilesFS.Open(fmt.Sprintf("files/%s", path))
	if err != nil {
		return returnPath
	}
	stat, err := file.Stat()
	if err != nil {
		return returnPath
	}
	if stat.IsDir() {
		return returnPath
	}
	return fmt.Sprintf("%s?ts=%d", returnPath, stat.Size())
}
