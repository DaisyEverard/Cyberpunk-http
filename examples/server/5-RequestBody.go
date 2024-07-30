// To test an http request with a body, run a curl command with the syntax:
// curl -X POST -d 'This is the body' 'http://localhost:3333?first=1&second='
// For the form, its:
// curl -X POST -F 'myName=Sammy' 'http://localhost:3333/hello'
// For a form with no data to get an error message: 
// curl -v -X POST 'http://localhost:3333/hello'

// -F send form data
// -v verbose
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"io/ioutil"
)

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	body, err := ioutil.ReadAll(r.Body)
	// ioutil.ReadAll reads data from an io.Reader. r.Body is an io.Reader
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	fmt.Printf("%s: got / request. first(%t)=%s. second(%t)=%s, body:\n%s\n", ctx.Value(keyServerAddr), hasFirst, first, hasSecond, second, body)
	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))

	myName := r.PostFormValue("myName")
	// r.PostFormValue on includes values posted in a form in the body of a request
	// r.FormValue includes form body and any values in the query string. http...?myName=aName
	// use PostFormValue by default unless you want the flexibility to be able to put in a query string too
	if myName == "" {
		// You can just set the name to a default with myName = "HTTP"

		// OR
		// send the client an error message
		w.Header().Set("x-missing-field", "myName")
		w.WriteHeader(http.StatusBadRequest) // WriteHeader page status is send with all headers set on w.
		// The body can stil be written to after
		return 
	
	}
	io.WriteString(w, fmt.Sprintf("Hello, %s!\n", myName))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
 }