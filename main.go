package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("could not read body: %s", err)
	}

	fmt.Printf("%s: got / request, first(%t): %s, second(%t): %s, body: %s\n", ctx.Value(keyServerAddr), hasFirst, first, hasSecond, second, body)
	io.WriteString(w, "This is my website 🚀")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Printf(" %s : got /hello request\n", ctx.Value(keyServerAddr))

	myName := r.PostFormValue("myName")
	if myName == "" {
		w.Header().Set("x-missing-field", "myName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(w, fmt.Sprintf("Hello, %s\n", myName))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	ctx, cancelCtx := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("SERVER CLOSED! 📉")
		} else if err != nil {
			fmt.Printf("error starting the server: %s\n", err)
			cancelCtx()
		}
	}()
	<-ctx.Done()
}
