package web

import (
	"eco-passport-back/model"
	"eco-passport-back/web/context"
	"fmt"
	"net/http"
)

func webGetRegions(w http.ResponseWriter, r *http.Request) {
	regions, err := model.GetDatabase().SelectRegions()
	if err != nil {
		context.SetError(r, fmt.Errorf("get regions: %v", err))
		return
	}
	context.SetResponse(r, regions)
}
