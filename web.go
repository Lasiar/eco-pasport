package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func run() {

	api := http.NewServeMux()

	api.HandleFunc("/api/get-tree", webGetTree)

	webServer := &http.Server{
		Addr:           GetConfig().Port,
		Handler:        middlewareCORS(api),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := webServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func webGetTree(w http.ResponseWriter, r *http.Request) {
	t := GetEpTree()
	info := new(TablesInfo)

	if err := info.GetTables(); err != nil {
		log.Fatal(err)
	}

	ChangeName(t.TreeItem, info)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(t); err != nil {
		GetConfig().Err.Printf("[WEB] error encod json %v", err)
	}

}

func ChangeName(t []*nodeEpTree, table *TablesInfo) {
	for _, node := range t {

		if node.Name == "" {
			id, err := strconv.Atoi(node.TableID)
			if err != nil {
				log.Println(err)
			}

			if table, ok := (*table)[id]; ok {
				node.Name = table.VisName
			}

		}

		ChangeName(node.TreeItem, table)
	}

}
