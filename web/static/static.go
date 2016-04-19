package static

//go:generate go-bindata -pkg static -o static_gen.go files/...

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

func FileSystem() http.FileSystem {
	fs := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "files"}
	return &binaryFileSystem{
		fs,
	}
}

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name[1:])
}
