package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"log"
	"bufio"
	"time"
)

const DEFAULT_CACHE_PATH = "/.cache"

type cacheFile struct {
	name string
	file []byte
}

type cachePool struct {
	path  string
	files map[string]*cacheFile
	cachedTimes map[string]time.Time

}

type cachePool2 struct {
	path  string
	fileNames []string
}

func (c *cachePool) Put(filename string, file []byte) {

	cFile := c.files[filename]
	if(&cFile != nil){
		t := c.cachedTimes[filename]

	}

	filePath := filepath.Join(c.path, filename)

	error = saveFile(filePath, file)
	if(error == nil) {

		/*
		t1 := time.Now()
		t2 := time.Now()
		fmt.Println("t1 ", t1)
		fmt.Println("t2 ", t2)
		fmt.Println("t? ", t1.Equal(t2))
		*/
	//	t := time.
	//	c.files[filename] = &cacheFile{filename, file}
	//
	}
}

func (c *cachePool) Get(file string) ([]byte, error) {
	f, err := os.Open(file)
	if (hasError(err)) {
		return nil, err
	}
	//t := &oauth2.Token{}
	//err = json.NewDecoder(f).Decode(t)

	//fbytes := f.
	defer f.Close()
	//return fbytes, err
	return nil, err
}

func (c *cachePool) Remove(filename string)  {
	cFile := c.files[filename]
	if (&cFile == nil) {
		println("File with name ", filename, " is not exist!")
	}else {
		removeFile(filename)
		c.files[filename] = nil
		delete(c.files, filename)
	}
}

// Internal
/*func updateFile(filePath string) error {
	fmt.Printf("Update file: %s\n", filePath)
	if _, err := os.Stat("/path/to/whatever"); err == nil {
		// path/to/whatever exists
	}
}*/

func removeFile(filePath string) error {
	error := os.Remove(filePath)
	hasError(error)
	return error
}

func saveFile(filePath string, b []byte) error {
	fmt.Printf("Saving file to: %s\n", filePath)
	f, err := os.Create(filePath)
	f.Sync()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	re, err := w.Write(b) //WriteString("buffered\n")
	hasError(err)
	fmt.Printf("wrote %d bytes\n", re)

	w.Flush()

	return err
	//json.NewEncoder(f).Encode(token)
}

func NewCachePool(path string) *cachePool {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if (err != nil) {
		log.Fatal(err)
	}
	fmt.Println("dir", dir)
	path = dir + path
	err = os.MkdirAll(path, 0700)
	hasError(err)
	return &cachePool{path, make(map[string]*cacheFile), make(map[string]time.Time)}
}

func hasError(e error) bool {
	if e != nil {
		log.Fatal(e)
		//panic(e)
		return true
	}
	return false
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
