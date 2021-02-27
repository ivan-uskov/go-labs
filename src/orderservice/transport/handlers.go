package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderservice/model"
	"time"
)

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "hello-world")
}

func orders(w http.ResponseWriter, _ *http.Request) {
	orders := model.OrdersList{
		Orders: []model.Order{
			{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", MenuItems: []model.MenuItem{{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Quantity: 0}}},
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

func order(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["ID"]
	if found {
		if id == "3fa85f64-5717-4562-b3fc-2c963f66afa6" {
			order := model.OrderDetails{
				Order: model.Order{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", MenuItems: []model.MenuItem{{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Quantity: 0}}},
				Cost:  1,
				Time:  1,
			}

			renderJson(w, order)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "Not found")
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
			"time":       time.Now().Format(time.RFC3339),
		}).Info("got request")
		h.ServeHTTP(w, r)
	})
}

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/orders", orders).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z]+}", order).Methods(http.MethodGet)

	return logMiddleware(r)
}
