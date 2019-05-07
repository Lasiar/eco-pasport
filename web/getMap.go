package web

import (
	"EcoPasport/model"
	"EcoPasport/web/context"
	"fmt"
	"net/http"
)

func webGetMap(w http.ResponseWriter, r *http.Request) {
	req := struct {
		RegionID int `json:"region_id"`
	}{}

	if err := parseJSON(r, &req); err != nil {
		context.SetError(r, fmt.Errorf("json encode %v", err))
		return
	}

	center, points, err := model.GetDatabase().GetMap(req.RegionID)
	if err != nil {
		context.SetError(r, fmt.Errorf("get map %v", err))
		return
	}

	response := struct {
		Center *[2]float64
		Points []model.Point
	}{}

	response.Points = points
	response.Center = center

	context.SetResponse(r, response)
}
