package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpGet(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "", nil)
	req.RequestURI = "/sodoku"
	req.Header.Set("Accept", "application/json")

	sh := SodokuHandler{
		Timeout:      1,
		MinSolutions: 1,
	}
	sh.ServeHTTP(recorder, req)

	//  verify response
	assert.Equal(t, http.StatusOK, recorder.Code)

	// decode json
	dec := json.NewDecoder(recorder.Body)
	var resp response
	err := dec.Decode(&resp)
	if err != nil {
		t.Errorf("Expected nil, actual: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("Expected nil, actual: %v", resp.Error)
	}

	if len(resp.Solutions) != 1 {
		t.Errorf("Expected 1, actual: %d", len(resp.Solutions))
	}

	if len(resp.Solutions[0].Steps) != 81 {
		t.Errorf("Expected 81, actual: %d", len(resp.Solutions[0].Steps))
	}
}
