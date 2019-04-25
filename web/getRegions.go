package web

import (
	"EcoPasport/model"
	context "EcoPasport/web/context"
	"fmt"
	"net/http"
)

func webGetRegions(w http.ResponseWriter, r *http.Request) {
	regions, err := model.NewDatabase().GetRegions()
	if err != nil {
		context.SetError(r, fmt.Errorf("get regions: %v", err))
		return
	}

	context.SetResponse(r, regions)
}
