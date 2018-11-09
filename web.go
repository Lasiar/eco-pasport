package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func run() {

	api := http.NewServeMux()

	api.HandleFunc("/get-tree", webGetTree)
	api.HandleFunc("/get-regions", webGetRegions)
	api.HandleFunc("/get-table", webGetTable)

	webServer := &http.Server{
		Addr:           GetConfig().Port,
		Handler:        http.StripPrefix("/api", middlewareCORS(api)),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := webServer.ListenAndServe(); err != nil {
		GetConfig().Err.Fatalf("Ошибка запуска сервера %v", err)
	}
}

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		GetConfig().Info.Println(r.RemoteAddr)

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == http.MethodOptions {
			return
		}

		if r.Method != http.MethodPost {
			printWarnLog(r, w, "method not allowed")
			return
		}

		printInfoLog(r)
		next.ServeHTTP(w, r)
	})
}

func webGetTree(w http.ResponseWriter, r *http.Request) {
	t := GetEpTree()

	meta := GetTablesMeta()

	ChangeName(t.TreeItem, meta)

	fmt.Println(t)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(t); err != nil {
		printWarnLog(r, w, fmt.Sprint("[WEB] json ecode", err))
		return
	}
}

func webGetRegions(w http.ResponseWriter, r *http.Request) {
	regions := new(Regions)

	if err := regions.Fetch(); err != nil {
		printWarnLog(r, w, fmt.Sprint("[WEB]", err))
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(regions); err != nil {
		printWarnLog(r, w, fmt.Sprint("[WEB] json encode", err))
		return
	}
}

func webGetTable(w http.ResponseWriter, r *http.Request) {

	tblInfo := new(RequestTableInfo)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tblInfo); err != nil {
		printWarnLog(r, w, fmt.Sprint("[WEB] error deode json: %v", err))
		return
	}

	t := new(Table)
	if err := t.Fetch(tblInfo); err != nil {
		printWarnLog(r, w, fmt.Sprint("[WEB] json encode", err))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(t)

}

func ChangeName(t []*nodeEpTree, table *TablesMeta) error {
	var sumError []string
	for _, node := range t {
		if node.Name == "" {
			id, err := strconv.Atoi(node.TableID)
			if err != nil {
				if err != nil {
					sumError = append(sumError, fmt.Sprint(err))
				}
			}

			if table, ok := (*table)[id]; ok {
				node.Name = table.VisName
			}

		}

		sumError = append(sumError, fmt.Sprint(ChangeName(node.TreeItem, table)))
	}
	return fmt.Errorf("%v", strings.Join(sumError, " "))
}


func printWarnLog(r *http.Request, w http.ResponseWriter, info string) {
	http.Error(w, "some errors", http.StatusServiceUnavailable)
	GetConfig().Warn.Printf("[WEB] %v connectMSSQL from %v, %v	", r.URL.Path, r.RemoteAddr, info)
}

func printInfoLog(r *http.Request) {
	GetConfig().Info.Printf("[WEB] %v connectMSSQL from %v", r.URL.Path, r.RemoteAddr)
}

