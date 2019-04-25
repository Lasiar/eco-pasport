package web

import (
	"EcoPasport/base"
	context "EcoPasport/web/context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Run settings and run web server on specified port in config
func Run() {

	apiMux := http.NewServeMux()

	apiMux.HandleFunc("/get-tree", webGetTree)
	apiMux.HandleFunc("/get-regions", webGetRegions)
	apiMux.HandleFunc("/get-table", webGetTable)
	apiMux.HandleFunc("/get-region-info", webRegionInfo)
	apiMux.HandleFunc("/get-region-map", webGetMap)

	logger := log.New(os.Stdout, "[connect] ", log.Flags())

	api := middlewareCORS(JSONWriteHandler(middlewareLogging(logger)(http.StripPrefix("/api", apiMux))))

	webServer := &http.Server{
		Addr:           base.GetConfig().Port,
		Handler:        api,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := webServer.ListenAndServe(); err != nil {
		base.GetConfig().Err.Fatalf("Ошибка запуска сервера %v", err)
	}
}

func middlewareLogging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), "time: ", time.Since(start))
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST, OPTIONS")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// JSONWriteHandler хандлер для ответа в виде json
func JSONWriteHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if next != nil {
			next.ServeHTTP(w, r)
		}

		if err := r.Context().Err(); err != nil {

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")

			switch err {
			case sql.ErrNoRows:
				http.Error(w, "Нет данных по данному запросы", 404)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusOK)

		data := r.Context().Value(context.ResponseDataKey)
		if data == nil {
			return
		}
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println(err)
		}
	})
}
