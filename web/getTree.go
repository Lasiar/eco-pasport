package web

import (
	"EcoPasport/model"
	context "EcoPasport/web/context"
	"fmt"
	"net/http"
)

func webGetTree(w http.ResponseWriter, r *http.Request) {
	data, err := model.GetTree()
	if err != nil {
		context.SetError(r, fmt.Errorf("error get tree: %v", err))
		return
	}

	context.SetResponse(r, data)
}
