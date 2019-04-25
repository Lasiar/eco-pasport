package web

import (
	"EcoPasport/model"
	"EcoPasport/web/context"
	"encoding/json"
	"fmt"
	"net/http"
)

func webGetMap(w http.ResponseWriter, r *http.Request) {
	req := struct {
		RegionID int `json:"region_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		context.SetError(r, fmt.Errorf("json encode %v", err))
		return
	}

	center, points, err := model.NewDatabase().GetMap(req.RegionID)
	if err != nil {
		context.SetError(r, fmt.Errorf("get map %v", err))
		return
	}

	response := struct {
		Center *[]float64
		Points []model.Point
	}{}

	response.Points = points
	response.Center = center

	context.SetResponse(r, response)
}
