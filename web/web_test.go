package web

import (
	"EcoPasport/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodPost, "/wef/wefwef/sdafsadf/asdfsdf", nil)
	if err != nil {
		t.Fatal(err)
	}

	database := new(model.Database)
	db, mock, err := sqlmock.New()
	database.SetDB(db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dataSet := []model.Region{
		{0, 0, "first", false},
		{1, 1, "second", true},
	}
	rows := sqlmock.NewRows([]string{
		"id",
		"num_region",
		"name",
		"isTown",
	})
	for _, d := range dataSet {
		rows.AddRow(d.ID, d.NumRegion, d.Name, d.IsTown)
	}

	mock.ExpectQuery("[a-z]*").WillReturnRows(rows)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webGetRegions)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	JSONWriteHandler(handler).ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
