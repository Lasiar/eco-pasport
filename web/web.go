package web

import (
	"EcoPasport/base"
	"EcoPasport/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type webError struct {
	Error   error
	Message string
}

type webHandler func(http.ResponseWriter, *http.Request) *webError

func (wh webHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := wh(w, r); e != nil {

		request := struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{}
		encoder := json.NewEncoder(w)

		switch e.Error {
		case sql.ErrNoRows:
			request.Message = "Нет данных по данному запросу"
			request.Code = 100
		default:
			request.Message = e.Message
		}

		log.Printf("[WEB] %v %v [METНOD] %v [URL] %v [USER AGENT] %v", e.Message, e.Error, r.Method, r.URL, r.UserAgent())

		w.WriteHeader(http.StatusInternalServerError)

		if err := encoder.Encode(request); err != nil {
			log.Printf("[WEB] %v", err)
		}

	}
}

func Run() {

	apiMux := http.NewServeMux()

	apiMux.Handle("/get-tree", webHandler(webGetTree))
	apiMux.Handle("/get-regions", webHandler(webGetRegions))
	apiMux.Handle("/get-table", webHandler(webGetTable))
	apiMux.Handle("/get-region-info", webHandler(webRegionInfo))
	apiMux.Handle("/get-region-map", webHandler(webGetMap))

	//	api := middlewareCORS(apiMux)

	logger := log.New(os.Stdout, "[connect] ", log.Flags())

	api := middlewareCORS(middlewareLogging(logger)(http.StripPrefix("/api", apiMux)))

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

func middlewareSetCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		if r.Method != http.MethodPost {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func webGetTree(w http.ResponseWriter, r *http.Request) *webError {

	res, err := model.GetTree()
	if err != nil {
		return &webError{err, fmt.Sprintf("ошибка получении дерева таблиц")}
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
		return &webError{err, fmt.Sprintf("json encode %v", err)}
	}
	return nil
}

func webRegionInfo(w http.ResponseWriter, r *http.Request) *webError {
	response := struct {
		RegionID int `json:"region_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&response); err != nil {
		return &webError{err, fmt.Sprintf("json encode %v", err)}
	}

	regionInfo, isEmpty, err := model.NewDatabase().GetRegionInfo(response.RegionID)
	if err != nil {
		return &webError{err, fmt.Sprintf("get region %v", err)}
	}
	encoder := json.NewEncoder(w)

	if !isEmpty {
		response := struct {
			Empty bool
		}{}

		response.Empty = !isEmpty

		if err := encoder.Encode(response); err != nil {
			return &webError{err, fmt.Sprintf("json encode %v", err)}
		}
		return nil
	}

	if err := encoder.Encode(regionInfo); err != nil {
		return &webError{err, fmt.Sprintf("json encode %v", err)}
	}
	return nil
}

func webGetRegions(w http.ResponseWriter, r *http.Request) *webError {
	regions, err := model.NewDatabase().GetRegions()
	if err != nil {
		return &webError{err, fmt.Sprintf("get region %v", err)}
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(regions); err != nil {
		return &webError{err, fmt.Sprintf("json encode %v", err)}
	}
	return nil
}

func webGetTable(w http.ResponseWriter, r *http.Request) *webError {

	tblInfo := &struct {
		User     string `json:"user"`
		RegionID int    `json:"region_id"`
		TableID  int    `json:"table_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tblInfo); err != nil {
		return &webError{err, fmt.Sprintf("json decode %v", err)}
	}

	t, err := model.NewDatabase().GetTable(tblInfo.User, tblInfo.RegionID, tblInfo.TableID)
	if err != nil {
		return &webError{err, fmt.Sprintf("json encode %v", err)}
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(t); err != nil {
		return &webError{err, fmt.Sprintf("json encode %v", err)}
	}
	return nil
}

func webGetMap(w http.ResponseWriter, r *http.Request) *webError {
	req := struct {
		RegionID int `json:"region_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		return &webError{err, fmt.Sprintf("json encode %v", err)}
	}

	center, points, err := model.NewDatabase().GetMap(req.RegionID)
	if err != nil {
		return &webError{err, fmt.Sprintf("get map %v", err)}
	}

	response := struct {
		Center *[]float64
		Points []model.Point
	}{}

	response.Points = points
	response.Center = center

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		return &webError{err, fmt.Sprintf("json encode %v", err)}
	}

	//fmt.Fprint(w, `{"Center":[56.26358,90.49446],"Points":[{"Name":"МКУ \"Центр бухучета\"","Address":"662150, Красноярский край, г. Ачинск, 1-й микрорайон, 27, пом. 1","WasteGenerationForTheYear":5.675,"Latitude":56.26278,"Longitude":90.48457},{"Name":"ОСП Ачинский почтамт УФПС Красноярского края-филиала ФГУП \"Почта России\"","Address":"662150, Красноярский край, г. Ачинск, 1-й микрорайон, 43","WasteGenerationForTheYear":50.139,"Latitude":56.26433,"Longitude":90.49235}]}`)

	return nil
}
