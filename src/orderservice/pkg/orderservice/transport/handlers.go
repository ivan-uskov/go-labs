package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"orderservice/pkg/orderservice/application/data"
	query2 "orderservice/pkg/orderservice/application/query"
	"orderservice/pkg/orderservice/application/service"
	"orderservice/pkg/orderservice/infrastructure/query"
	"orderservice/pkg/orderservice/infrastructure/repository"
	"time"
)

type server struct {
	orderService      service.OrderService
	orderQueryService query2.OrderQueryService
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, "hello-world")
	if err != nil {
		log.Error(err)
	}
}

func processError(w http.ResponseWriter, e error) {
	if e == data.InternalError {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		http.Error(w, e.Error(), http.StatusBadRequest)
	}
}

func (s *server) getOrdersList(w http.ResponseWriter, _ *http.Request) {
	orders, err := s.orderQueryService.GetOrders()
	if err != nil {
		processError(w, err)
		return
	}

	renderJson(w, orders)
}

func (s *server) getOrderInfo(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["ID"]
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	info, err := s.orderQueryService.GetOrderInfo(id)
	if err != nil {
		processError(w, err)
		return
	}

	if info == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	renderJson(w, info)
}

func (s *server) deleteOrder(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["ID"]
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	err := s.orderService.Delete(id)
	if err != nil {
		log.Error(err)
		processError(w, err)
	}
}

func (s *server) updateOrder(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["ID"]
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	var orderRequest service.UpdateOrderRequest
	err := jsonFromRequest(r, &orderRequest)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = s.orderService.Update(id, orderRequest)
	if err != nil {
		processError(w, err)
	}
}

func (s *server) addOrder(w http.ResponseWriter, r *http.Request) {
	orderRequest := service.AddOrderRequest{}
	err := jsonFromRequest(r, &orderRequest)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = s.orderService.Add(orderRequest)

	if err != nil {
		processError(w, err)
	}
}

func renderJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error(err)
		http.Error(w, "server error", http.StatusInternalServerError)
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
	s.HandleFunc("/orders", srv.getOrdersList).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z-]+}", srv.getOrderInfo).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z-]+}", srv.deleteOrder).Methods(http.MethodDelete)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z-]+}", srv.updateOrder).Methods(http.MethodPut)
	s.HandleFunc("/order", srv.addOrder).Methods(http.MethodPost)

	return logMiddleware(r)
}

func makeServer(db *sql.DB) *server {
	return &server{
		orderService:      service.NewOrderService(repository.NewOrderRepository(db)),
		orderQueryService: query.NewOrderQueryService(db),
	}
}
