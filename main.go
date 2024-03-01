package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/lestrrat-go/httprc"
)

const httpTimeout = 5 * time.Second
const hostUrl = "0.0.0.0:41234"

var bytesResponse = []byte("response")

func main() {
	fmt.Println("SETUP: spinning up server")
	go runServer()

	time.Sleep(1 * time.Second)

	cache := newCache(3 * time.Second)
	client := http.Client{}
	client.Timeout = httpTimeout

	fmt.Println("SETUP: registering servers with cache")
	err := cache.Register("http://"+hostUrl, httprc.WithHTTPClient(&client),
		httprc.WithRefreshInterval(5*time.Second),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		count := 1
		for {
			fmt.Printf("CLIENT: cache get at   %s\n", time.Now().Format(time.StampMilli))
			val, err := cache.Get(context.Background(), "http://"+hostUrl)
			if err != nil {
				fmt.Println("CLIENT: cache get err", err)
				time.Sleep(1000 * time.Millisecond)
			}
			valPrint := "nil"
			if valBytes, ok := val.([]byte); ok && valBytes != nil {
				valPrint = string(valBytes)
			}
			fmt.Printf("CLIENT: cache value at %s: %s\n", time.Now().Format(time.StampMilli), valPrint)

			time.Sleep(1000 * time.Millisecond)
			count++
		}
	}()
	fmt.Println("SETUP: sleeping for 2 hours to watch issue in terminal")
	time.Sleep(2 * time.Hour)
}

func newCache(window time.Duration) *httprc.Cache {
	return httprc.NewCache(
		context.Background(),
		httprc.WithRefreshWindow(window),
	)
}

func runServer() {
	mux := http.NewServeMux()
	counter := 0
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("SERVER: GET received")
		// 1 in 3 chance server fails to respond in a timely fashion
		if rand.N(3) == 0 {
			fmt.Println("SERVER: will time out")
			time.Sleep(httpTimeout + 3*time.Second)
		}
		w.Write([]byte(fmt.Sprintf("count %d", counter)))
		counter++
	})

	err := http.ListenAndServe(hostUrl, mux)
	fmt.Printf("SERVER: exited: %s \n", err.Error())
	os.Exit(1)
}
