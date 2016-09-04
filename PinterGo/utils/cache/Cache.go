package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"log"
	"bufio"
)

const DEFAULT_CACHE_PATH = "/.cache"

type CacheFile struct {
	name string
	file []byte
}

type CachePool struct {
	path string
	files map[string]*CacheFile
}

func (c *CachePool) put(filename string, file []byte) error {
	fmt.Println("test")
	filePath := filepath.Join(c.path, filename)
	saveFile(filePath)


	cacheFile := CacheFile{"test", ""}
	c.files[filename] = *cacheFile

	// os.MkdirAll(tokenCacheDir, 0700)
	// filepath.Join(tokenCacheDir, url.QueryEscape("drive-go-quickstart.json")), err

	return nil
}

func saveFile(filePath string, b []byte) {
	fmt.Printf("Saving file to: %s\n", filePath)
	f, err := os.Create(filePath)
	f.Sync()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	re, err := w.Write(b) //WriteString("buffered\n")
	fmt.Printf("wrote %d bytes\n", re)

	w.Flush()

	//json.NewEncoder(f).Encode(token)
}

func NewCachePool(path string, file *[]byte) *CachePool {
	os.MkdirAll(path, 0700)
	return CachePool{path, new(map[string]*CacheFile)}

}

/*
func NewFile(fd int, name string) *CachePool {
	if fd < 0 {
		return nil
	}
	pool :=CachePool{}
	f := File{fd, name, nil, 0}
	return &f
}*/
