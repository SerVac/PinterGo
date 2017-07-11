package main

import (
	"net/http"
	"log"
	"golang.org/x/oauth2"
	"fmt"
	"net/url"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"./utils/cache"
)

const htmlIndex = `<html><body>
Hello!
<a href="http://localhost/boards" class="button">Go to boards</a>
</body></html>
`

const fPath = "PinterGo/"
const serverCrt = fPath + "server.crt"
const serverKey = fPath + "server.key"

/*
	https://api.pinterest.com/oauth/?
	response_type=code&
	redirect_uri=https://mywebsite.com/connect/pinterest/&
	client_id=12345&
	scope=read_public,write_public&
	state=768uyFys
*/
const client_id = ""
const client_secret = ""
const redirect_url = "https://localhost/returnPage"

var (
	conf *oauth2.Config
	tok *oauth2.Token
	client *http.Client
)

const url_mi = "/me"
const url_pins = "/pins"
const url_boards = "/boards"
const api_host_url = "https://api.pinterest.com/v1"

const api_url_get_pins = api_host_url + url_mi + url_pins

const api_url_get_boards = api_host_url + url_mi + url_boards
const api_url_get_board_fmt_pins = api_host_url + url_boards+"/%d"+url_pins

type Image struct {
	URL   string `json:"url"`
}
type Count struct {
	name   string
	num   int32
}
//map<string,i32>

type PinEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	LINK  string `json:"link"`
	COLOR  string `json:"color"`
	//COUNTS  Count `json:"counts"`
	IMAGE Image `json:"image"`
	//IMAGES []Image `json:"image"`
}

type PinsData struct {
	Data []PinEntity `json:"data"`
}

type Board struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type BoardsData struct {
	Data []Board `json:"data"`
}

func createURL(url_link string, params map[string]string) *url.URL {
	u, err := url.Parse(url_link)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	fmt.Println("url_link = ", url_link)
	fmt.Println(" +params = ", params)
	for k, v := range params {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	fmt.Println(u)
	fmt.Println("str = " + u.String())
	return u
	/*
	u, err := url.Parse(api_host_url + api_url_boards)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Add("access_token", tok.AccessToken)
	q.Add("fields", "id,name,url")
	//q.Add("limit", "10")
	u.RawQuery = q.Encode()
	fmt.Println(u)
	fmt.Println("str = " + u.String())
	*/
}


// boards
func handleBoards(w http.ResponseWriter, r *http.Request) {

	url_link := createURL(api_url_get_boards, map[string]string{"access_token": tok.AccessToken, "fields": "id,name,url" })

	resp, err := client.Get(url_link.String())
	if err != nil {
		log.Println(err)
	}  else {
		//readFromRespToBuffer(buf, resp)
		//var boards []int
		responseBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}else {

			var f interface{}
			err := json.Unmarshal(responseBytes, &f)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Parse response JSON")
			fmt.Println("f =", f)
			fmt.Println("&f =", &f)
			//fmt.Println("*f =", *f) //invalid indirect of f (type interface {})
			parseUnknownJSON(f)

			pinData := BoardsData{}
			err = json.Unmarshal(responseBytes, &pinData)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("pinData =", pinData)

			pinsSlice := []string{}
			for _, board := range pinData.Data {
				pinNumber, err := strconv.ParseInt(board.ID, 10, 64)
				if err != nil {
					fmt.Println(err)
				}else {
					pins_url := getPinsUrlForBoard(pinNumber)
					pinsSlice = append(pinsSlice, pins_url)
					//url_link := createURL(, map[string]interface{}{"access_token":, "fields":  })
				}
			}

			// get boards pins
			for _, pins_url := range pinsSlice {
				fmt.Println("-----------\n")
				fmt.Println("pins_url =", pins_url)
				// &fields=id,creator,note
				pins_url = createURL(pins_url, map[string]string{"fields": "id,name,url,color,counts,image" }).String()

				resp, err := client.Get(pins_url)
				if err != nil {
					log.Println(err)
				}  else {
					pinsResponseBytes, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Println(err)
					}else {
						pinData := PinsData{}
						fmt.Println("resp.Body =", resp.Body)
						err = json.Unmarshal(pinsResponseBytes, &pinData)
						if err != nil {
							fmt.Println(err)
						}else{
							for _, entity := range pinData.Data {
								fmt.Println("Pins url = ", entity.URL)
								fmt.Println("Pins color = ", entity.COLOR)
								//fmt.Println("Pins counts = ", entity.COUNTS)
								fmt.Println("Pins img = ", entity.IMAGE)
								//fmt.Println("Pins imges = ", entity.IMAGES)


								//for _, img := range entity.IMAGES {
								//	fmt.Println(" img = ", img.URL)
								//}

							}
						}
					}
				}
			}

		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body>`))
		//w.Write([]byte(boards))
		w.Write(responseBytes)
		//w.Write([]byte(`<a href="http://localhost/boards">Go to Boards</a>`))
		w.Write([]byte(`</body></html>`))

	}

}

func getPinsUrlForBoard(boardNumber int64) string{
	return fmt.Sprintf(api_url_get_board_fmt_pins, boardNumber)
}

func parseUnknownJSON(unknownInterface interface{}) {
	unknMap := unknownInterface.(map[string]interface{})
	for k, v := range unknMap {
		switch t := v.(type) {
		case string:
			fmt.Println(k, "is string", "t = ", t, " v=", v)
		case int:
			fmt.Println(k, "is int", "t = ", t, " v=", v)
		case []interface{}:
			fmt.Println(k, "is an array:")

			for i, av := range t {
				fmt.Println("")
				fmt.Println(i, av)
				fmt.Println("with addres", i, &av)

				parseUnknownJSON(av)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}

func readFromRespToBuffer(bufReader *bytes.Buffer, resp *http.Response) {
	defer resp.Body.Close()
	_, err := bufReader.ReadFrom(resp.Body)
	//_, err := io.Copy(os.Stdout, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//return bufReader
}

// /returnPage
func handleReturn(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	fmt.Println("state = ", state)
	code := r.FormValue("code")
	fmt.Println("code = ", code)

	var err error
	tok, err = conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	//api_url+"/me/pins/?"+"access_token=<YOUR-ACCESS-TOKEN>
	//&fields=id,creator,note
	//&limit=1
	u, err := url.Parse(api_url_get_pins)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Add("access_token", tok.AccessToken)
	q.Add("fields", "id,creator,note")
	q.Add("limit", "1")
	u.RawQuery = q.Encode()
	fmt.Println(u)
	fmt.Println("str = " + u.String())

	client = conf.Client(oauth2.NoContext, tok)

	buf := new(bytes.Buffer)
	resp, err := client.Get(u.String())
	if err != nil {
		log.Println(err)
	}  else {
		readFromRespToBuffer(buf, resp)
		fmt.Println("buf = " + buf.String())
		//defer resp.Body.Close()
		//buf.ReadFrom(resp.Body)
		//_, err := io.Copy(os.Stdout, resp.Body)
		/*if err != nil {
			log.Fatal(err)
		}*/
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if (buf != nil) {
		w.Write([]byte(`<html><body>`))
		w.Write([]byte(buf.String()))
		w.Write([]byte(`<a href="https://localhost/boards">Go to Boards</a>`))
		w.Write([]byte(`</body></html>`))
	}
	// http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)

}


// /
func handleMain(w http.ResponseWriter, r *http.Request) {
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	authURL := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	//fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Use the authorization code that is pushed to the redirect URL.
	// NewTransportWithCode will do the handshake to retrieve
	// an access token and initiate a Transport that is
	// authorized and authenticated by the retrieved token.

	/*var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}*/

	/*
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}*/

	//client := conf.Client(oauth2.NoContext, tok)
	//client.Get("...")

	/*
		resp, err := http.Get(url_info)
		if (err != nil) {
			fmt.Print("Some error: ", err)
		}

		var reader io.ReadCloser
		typeEnc := resp.Header.Get("Content-Encoding")
		switch  typeEnc{
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			defer reader.Close()
		default:
			reader = resp.Body
		}

		io.Copy(os.Stdout, reader)
		fmt.Println("Body =  ", reader)
		fmt.Println("No erroro ", resp)*/

}

func main() {
	//fmt.Println("Started running on http://127.0.0.1:8080")
	//fmt.Println("Started running on https://localhost:8080")
	//fmt.Println(http.ListenAndServe("localhost:8080", nil))
	//var t utils.CacheFile = utils.CacheFile{}

	//CachePool{path, make(map[string]*CacheFile)}

	pool := cache.NewCachePool("/.test")
	pool.Put("test_filename.txt", []byte("tetdt_12391291239"))

	pool.Get("test_filename.txt")

	//dat, err := ioutil.ReadFile(f.Read(f.))
	//check(err)
	//fmt.Print(string(dat))

	/*conf = &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Scopes:       []string{"read_public", "write_public"},
		RedirectURL: redirect_url,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.pinterest.com/oauth/",
			TokenURL: "https://api.pinterest.com/v1/oauth/token",
		},
	}

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/returnPage", handleReturn)
	http.HandleFunc("/boards", handleBoards)

	//err := http.ListenAndServeTLS(":10443", serverCrt, serverKey, nil)
	err := http.ListenAndServeTLS("localhost:443", serverCrt, serverKey, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		os.Exit(1)
	}*/

}