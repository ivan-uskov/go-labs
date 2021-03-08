package transport

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"orderservice/pkg/orderservice/service"
	"testing"
)

func TestOrdersList(t *testing.T) {
	w := httptest.NewRecorder()
	getOrdersList(w, nil)
	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have: %d, want: %d", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	items := service.OrdersList{}
	if err = json.Unmarshal(jsonString, &items); err != nil {
		t.Errorf("Can't parse json: %s response with error %v", jsonString, err)
	}
}
