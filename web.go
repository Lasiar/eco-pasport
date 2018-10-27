package main

import (
	"encoding/json"
	"fmt"
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

	for _, node := range t.TreeItem {
		if len(node.TableID) != 4 {
			continue
		}
		id, err := strconv.Atoi(node.TableID)
		if err != nil {
			log.Println(err)
		}

		info := new(TablesInfo)

		if err := info.GetTables(); err != nil {
			log.Fatal(err)
		}

		for _, region := range *info {

			if id == region.ID {
				node.Name = region.VisName
			}

		}

		if len(t.TreeItem) != 0 {
			ChangeName(node, info)
		}
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(t)

}

func ChangeName(t *nodeEpTree, table *TablesInfo) {
	for _, node := range t.TreeItem {
		if len(node.TableID) != 4 {
			continue
		}

		id, err := strconv.Atoi(node.TableID)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, region := range *table {
			if id == region.ID {
				node.Name = region.VisName
				fmt.Println(id, region.VisName)
			}
		}

		if len(t.TreeItem) != 0 {
			ChangeName(node, table)
		}
	}

}
