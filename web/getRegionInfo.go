package web

import (
	"EcoPasport/model"
	"EcoPasport/web/context"
	"encoding/json"
	"fmt"
	"net/http"
)

func webRegionInfo(w http.ResponseWriter, r *http.Request) {
	response := struct {
		RegionID int `json:"region_id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		context.SetError(r, fmt.Errorf("json encode: %v", err))
		return
	}

	regionInfo, isEmpty, err := model.NewDatabase().GetRegionInfo(response.RegionID)
	if err != nil {
		context.SetError(r, fmt.Errorf("get region: %v", err))
		return
	}

	if !isEmpty {
		response := struct {
			Empty bool
		}{}

		response.Empty = !isEmpty

		context.SetResponse(r, response)
		return
	}

	context.SetResponse(r, regionInfo)
}
