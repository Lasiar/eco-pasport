package web

import (
	"eco-passport-back/model"
	"eco-passport-back/web/context"
	"fmt"
	"net/http"
)

func webRegionInfo(w http.ResponseWriter, r *http.Request) {
	response := struct {
		RegionID int `json:"region_id"`
	}{}

	if err := parseJSON(r, &response); err != nil {
		context.SetError(r, fmt.Errorf("json encode: %v", err))
		return
	}

	regionInfo, isEmpty, err := model.GetDatabase().GetRegionInfo(response.RegionID)
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
