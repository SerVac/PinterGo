package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"log"
	"io/ioutil"
)

const DEFAULT_CACHE_PATH = "/.cache"

type cacheFile struct {
	name string
	file []byte
}

type cachePool struct {
	path        string
	files       map[string]*os.File
	//cachedTimes map[string]time.Time
}


func NewCachePool(pathname string) *cachePool {
	pathPoolName := filepath.Join(DEFAULT_CACHE_PATH, pathname)
	err := os.MkdirAll(pathPoolName, os.ModePerm)
	if(err != nil){
		log.Fatal("CachePool.newCachePool.makeDir error: ",err)
		return nil
	}
	dir, err := filepath.Abs(filepath.Dir(pathPoolName))
	println("abs dir = "+dir)

	return &cachePool{path:pathPoolName, files: map[string]*os.File{}}
}

func (c *cachePool) Put(filename string, file []byte) {

	cFile := c.files[filename]
	if (&cFile != nil) {
		//c.cachedTimes[filename] = time.Now()

	}

	filePath := filepath.Join(c.path, filename)

	err := saveFile(filePath, file)
	hasError(err)

}

func (c *cachePool) Get(filename string) (*os.File, error) {
	fullfilename := getFullFilePath(c, filename)
	f, err := os.Open(fullfilename)
	println("fullFile ", fullfilename)
	if (hasError(err)) {
		return nil, err
	} else {
		b, errr := ioutil.ReadFile(fullfilename)
		if(errr != nil){
			fmt.Println(errr)
		}
		fmt.Println(string(b))
		return f, err
	}


	//t := &oauth2.Token{}
	//err = json.NewDecoder(f).Decode(t)

	//fbytes := f.
	defer f.Close()
	//return fbytes, err
	return nil, err
}

func getFullFilePath(c *cachePool, filename string) string {
	join := filepath.Join(c.path, filename)
	/*println("join = ", join)
	dir_join := filepath.Dir(join)
	println("dir_join = ", dir_join)
	abs_dir, err :=filepath.Abs(dir_join)
	println("abs_dir = ", abs_dir," err=",err)*/

	s, err := filepath.Abs(join)
	hasError(err)
	return s
}

func (c *cachePool) Remove(filename string) {
	cFile := c.files[filename]
	if (&cFile == nil) {
		println("File with name ", filename, " is not exist!")
	} else {
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
	hasError(err)
	defer f.Close()

	nb,err := f.Write(b)
	hasError(err)
	log.Println("Writed ",nb," bytes to path:",filePath)
	//f.Sync()
	//w := bufio.NewWriter(f)
	//re, err := w.Write(b) //WriteString("buffered\n")
	//hasError(err)
	//fmt.Printf("wrote %d bytes\n", re)
	//w.Flush()

	return err
	//json.NewEncoder(f).Encode(token)
}


func hasError(e error) bool {
	if e != nil {
		//panic(e)
		log.Fatalf("Unable to cache oauth token: %v", e)
		return true
	}
	return false
}
