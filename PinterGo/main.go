package main

import (
	"net/http"
	"log"
	"golang.org/x/oauth2"
	"fmt"
)

const htmlIndex = `<html><body>
Hello!
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
const client_id = "your_client_id"
const client_secret = "your_secret"
const redirect_url = "https://localhost/returnPage"
//redirect_url:="http://www.test.com"

var (
	conf *oauth2.Config
)


// /returnPage
func handleReturn(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	fmt.Println("state = ", state)
	code := r.FormValue("code")
	fmt.Println("code = ", code)

	//http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`test`))

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

	err := http.ListenAndServeTLS("localhost:443", serverCrt, serverKey, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}