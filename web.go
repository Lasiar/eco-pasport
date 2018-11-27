package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func run() {

	apiMux := http.NewServeMux()

	apiMux.HandleFunc("/get-tree", webGetTree)
	apiMux.HandleFunc("/get-regions", webGetRegions)
	apiMux.HandleFunc("/get-table", webGetTable)

	api := middlewareCORS(apiMux)

	staticMux := http.NewServeMux()

	staticMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "data/index.html")
	})

	staticMux.Handle("/api/", http.StripPrefix("/api", api))

	staticMux.Handle("/data/", http.StripPrefix("/data/", middlewareSetCacheControl(http.FileServer(http.Dir("./data")))))

	webServer := &http.Server{
		Addr:           GetConfig().Port,
		Handler:        staticMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := webServer.ListenAndServe(); err != nil {
		GetConfig().Err.Fatalf("Ошибка запуска сервера %v", err)
	}
}

func middlewareLogging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func middlewareSetCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("fuck")
		w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
		w.Header().Set("Cache-Control", "max-age:290304000, public")
		w.Header().Set("Expires", time.Now().AddDate(60, 0, 0).Format(http.TimeFormat))
		next.ServeHTTP(w, r)
	})
}

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		next.ServeHTTP(w, r)
	})
}

func webGetTree(w http.ResponseWriter, r *http.Request) {
	res := struct {
		*TablesMeta
		*epTree
	}{}

	res.epTree = GetEpTree()

	res.TablesMeta = GetTablesMeta()

	changeName(res.TreeItem, res.TablesMeta)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
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
		printWarnLog(r, w, fmt.Sprintf("[WEB] error deode json: %v", err))
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

func changeName(t []*nodeEpTree, table *TablesMeta) error {
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

		sumError = append(sumError, fmt.Sprint(changeName(node.TreeItem, table)))
	}
	return fmt.Errorf("%v", strings.Join(sumError, " "))
}

func printWarnLog(r *http.Request, w http.ResponseWriter, info string) {
	http.Error(w, "some errors", http.StatusServiceUnavailable)
	GetConfig().Warn.Printf("[WEB] %v connect from %v, %v	", r.URL.Path, r.RemoteAddr, info)
}

func printInfoLog(r *http.Request) {
	GetConfig().Info.Printf("[WEB] %v connectMSSQL from %v", r.URL.Path, r.RemoteAddr)
}
