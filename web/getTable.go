package web

import (
	"EcoPasport/model"
	"EcoPasport/web/context"
	"encoding/base64"
	"fmt"
	"net/http"
)

func webGetTable(w http.ResponseWriter, r *http.Request) {
	tblInfo := &struct {
		Key      string `json:"key"`
		User     string `json:"user"`
		RegionID int    `json:"region_id"`
		TableID  int    `json:"table_id"`
	}{}
	if err := parseJSON(r, tblInfo); err != nil {
		context.SetError(r, fmt.Errorf("json decode %v", err))
		return
	}
	userToken, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(tblInfo.Key)
	if err != nil {
		context.SetError(r, fmt.Errorf("key %v", err))
		return
	}
	userEmail, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(tblInfo.User)
	if err != nil {
		context.SetError(r, fmt.Errorf("user %v", err))
		return
	}
	if ac, err := model.GetDatabase().GetPrivilege(string(userEmail), string(userToken), tblInfo.TableID); err != nil || !ac {
		context.SetError(r, err)
		return
	}
	t, err := model.GetDatabase().GetTable(tblInfo.User, tblInfo.RegionID, tblInfo.TableID)
	if err != nil {
		context.SetError(r, fmt.Errorf("json encode %v", err))
		return
	}
	context.SetResponse(r, t)
}
