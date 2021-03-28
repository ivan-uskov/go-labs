package transport

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"orderservice/pkg/orderservice/application/data"
	"testing"
)

type mocOrderQueryService struct{}

func (m mocOrderQueryService) GetOrders() (*data.OrdersList, error) {
	return &data.OrdersList{
		Orders: []data.OrderInfo{
			{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", MenuItems: []data.MenuItem{{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Quantity: 0}}},
		},
	}, nil
}

func (m mocOrderQueryService) GetOrderInfo(id string) (*data.OrderInfo, error) {
	panic("implement me")
}

func TestOrdersList(t *testing.T) {
	srv := server{orderQueryService: mocOrderQueryService{}}
	w := httptest.NewRecorder()
	srv.getOrdersList(w, nil)
	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have: %d, want: %d", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	items := data.OrdersList{}
	if err = json.Unmarshal(jsonString, &items); err != nil {
		t.Errorf("Can't parse json: %s response with error %v", jsonString, err)
	}
}
