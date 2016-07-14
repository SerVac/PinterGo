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
	"os"
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

const api_host_url = "https://api.pinterest.com/v1/"
const api_url_pins = "me/pins/"
const api_url_boards = "me/boards/"

// boards
func handleBoards(w http.ResponseWriter, r *http.Request) {
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

	resp, err := client.Get(u.String())
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
			print("f =",f)
			print("&f =",&f)
			parseUnknownJSON(f)
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

func parseUnknownJSON(unknownInterface interface{}) {

	unknMap := unknownInterface.(map[string]interface{})
	for k, v := range unknMap {
		switch t := v.(type) {
		case string:
			fmt.Println(k, "is string", "t = ", t, " v=", v)
		case int:
			fmt.Println(k, "is int", "t = " , t, " v=", v)
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
	u, err := url.Parse(api_host_url + api_url_pins)
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
	conf = &oauth2.Config{
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

	err := http.ListenAndServeTLS("localhost:443", serverCrt, serverKey, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		os.Exit(1)
	}

}