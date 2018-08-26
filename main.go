package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
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
	router.OPTIONS("/route", routeAllowed)
	router.POST("/route", routeCors)
	router.GET("/route", routeCors)
	router.PUT("/route", routeCors)
	router.DELETE("/route", routeCors)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func routeAllowed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	keys := reflect.ValueOf(r.Header).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	headres := strings.Join(strkeys, ",") + ", X-Naver-Client-Id,X-Naver-Client-Secret,X-TARGET-URL"

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", headres)
}

func routeCors(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	keys := reflect.ValueOf(r.Header).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	headres := strings.Join(strkeys, ",") + ", X-Naver-Client-Id,X-Naver-Client-Secret,X-TARGET-URL"

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", headres)

	hostname := r.Header.Get("X-TARGET-URL")

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, hostname, r.Body)
	if err != nil {
		fmt.Fprint(w, err)
		log.Fatalln("Fail CORS from `" +
			r.Host +
			"` to `" +
			hostname +
			"`: New Request - " + err.Error())
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
			"`: Client Do - " + err.Error())
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
			"`: RealAll - " + err.Error())
	}

	fmt.Fprint(w, string(data))
	log.Println("Route CORS from `" + r.Host + "` to `" + hostname + "`")
}
