package main

import (
	"fmt"
	"log"
	"net/http"
	"project/db"
	"project/info"
	"text/template"

	"github.com/nats-io/stan.go"
)

func main() {
	db_connect, err := db.ConnectDB(db.Config{"localhost", "5432", "postgres", "qwerty", "wb_info", "disable"})
	if err != nil {
		log.Fatal(err)
	}
	cache, err := db.GetDic(db_connect)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", main_page)
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {

		infoId := r.FormValue("id")
		if _, ok := cache[infoId]; !ok {
			fmt.Fprintf(w, "Id %s не найден", infoId)
		} else {
			tmpl, _ := template.ParseFiles("template/info.html")
			tmpl.Execute(w, info.GetJson(cache[infoId]))

		}
	})
	sc, err := stan.Connect("test-cluster", "aaa")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(1)
	sub, err := sc.Subscribe("foo2", func(m *stan.Msg) {
		id, data, err := info.ValidInfo(string(m.Data))
		if err == nil {
			if _, ok := cache[id]; !ok {
				db.InsertDB(db_connect, id, data)
				cache[id] = data
			}
		}
	}, stan.DeliverAllAvailable())
	fmt.Println(2)
	defer sub.Unsubscribe()
	if err != nil {
		log.Fatal(err)
	}
	http.ListenAndServe(":8080", nil)
}
func main_page(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/index.html")
}
