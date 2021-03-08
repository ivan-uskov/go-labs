package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"orderservice/pkg/orderservice/repository"
	"orderservice/pkg/orderservice/service"
	"time"
)

type server struct {
	orderService service.OrderService
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, "hello-world")
	if err != nil {
		log.Error(err)
	}
}

func getOrdersList(w http.ResponseWriter, _ *http.Request) {
	orders := service.OrdersList{
		Orders: []service.Order{
			{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", MenuItems: []service.MenuItem{{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Quantity: 0}}},
		},
	}

	renderJson(w, orders)
}

func getOrderInfo(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["ID"]
	if found {
		if id == "3fa85f64-5717-4562-b3fc-2c963f66afa6" {
			order := service.OrderInfo{
				Order: service.Order{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", MenuItems: []service.MenuItem{{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Quantity: 0}}},
				Cost:  1,
				Time:  1,
			}

			renderJson(w, order)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprint(w, "Not found")
	if err != nil {
		log.Error(err)
	}
}

func (s *server) addOrder(w http.ResponseWriter, r *http.Request) {
	orderRequest := service.AddOrderRequest{}
	err := jsonFromRequest(r, &orderRequest)
	if err != nil {
		log.Debugf("Can't parse request: %s", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = s.orderService.AddOrder(orderRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func renderJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func jsonFromRequest(r *http.Request, output interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer func() {
		log.Error(r.Body.Close())
	}()

	err = json.Unmarshal(b, &output)
	if err != nil {
		err = fmt.Errorf("can't parse %s to json", b)
	}

	return err
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
			"duration":   time.Since(startTime).String(),
			"at":         startTime,
		}).Info("got request")
	})
}

func Router(db *sql.DB) http.Handler {
	srv := makeServer(db)

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/orders", getOrdersList).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z]+}", getOrderInfo).Methods(http.MethodGet)
	s.HandleFunc("/order", srv.addOrder).Methods(http.MethodPost)

	return logMiddleware(r)
}

func makeServer(db *sql.DB) *server {
	return &server{orderService: service.NewOrderService(repository.NewOrderRepository(db))}
}
