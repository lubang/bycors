package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("`$PORT` must be set")
	}

	router := httprouter.New()

	router.NotFound = http.FileServer(http.Dir("public"))
	router.POST("/route", routeCors)
	router.GET("/route", routeCors)
	router.PUT("/route", routeCors)
	router.DELETE("/route", routeCors)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func routeCors(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	hostname := r.Header.Get("X-TARGET-URL")

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, hostname, r.Body)
	if err != nil {
		fmt.Fprint(w, err)
		log.Fatalln("Fail to route CORS from `" +
			r.Host +
			"` to `" +
			hostname +
			"`: " + err.Error())
		return
	}

	for k, v := range r.Header {
		if !strings.EqualFold(k, "X-TARGET-URL") {
			req.Header.Set(k, v[0])
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, err)
		log.Fatalln("Fail to route CORS from `" +
			r.Host +
			"` to `" +
			hostname +
			"`: " + err.Error())
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(w, err)
		log.Fatalln("Fail to route CORS from `" +
			r.Host +
			"` to `" +
			hostname +
			"`: " + err.Error())
		return
	}

	fmt.Fprint(w, string(data))
	log.Println("Route CORS from `" + r.Host + "` to `" + hostname + "`")
}
