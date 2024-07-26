package main

import ("errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
"io/ioutil"
"bytes")
const serverPort = 3333

func main() {

	// server
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("server: %s /\n", r.Method)
			fmt.Printf("server: query id: %s\n", r.URL.Query().Get("id"))
			fmt.Printf("server: content-type: %s\n", r.Header.Get("content-type")) 
			fmt.Printf("server: headers:\n")

			for headerName, headerValue := range r.Header {
				fmt.Printf("\t%s = %s\n", headerName, strings.Join(headerValue, ", "))
			}

			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("server: could not read request body: %s\n", err)
			}
			fmt.Printf("server: request body: %s\n", reqBody)

			fmt.Fprintf(w, `{"message": "hello!"}`)
		})
		server := http.Server {
			Addr: fmt.Sprintf(":%d", serverPort),
			Handler: mux,
		}
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("error running http server: %s\n", err)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)


	// client

	jsonBody := []byte(`{"client_message": "hello, server!"}`) // supply a json object as a string, which will be sent as an array of bytes
	// encoding/json package returns a []byte, not a string
	bodyReader := bytes.NewReader(jsonBody)
	// body needs to be of type io.Reader (interface) to satisfy reqirements of http.Request

	requestURL := fmt.Sprintf("http://localhost:%d?id=1234", serverPort)
	// add an id query to the requestURL

	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	// changed method to Post and added a body reader
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	// ReadAll reads io.Reader until error or data ends
	// datatype returned is []byte
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: response body: %s\n", resBody)
}
