package http

import (
	"encoding/json"
	"errors"
	"github.com/NadiaBat/permissionsChecker/rbac"
	"net/http"
	"strconv"
	"sync"
	"context"
)

type Server struct{}
type handler struct{ wg *sync.WaitGroup }

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/check/":
		result, _ := routeCheck(r)
		w.Write(result)
	default:
	}
}

func (h Server) Serve() {
	handler := handler{}
	server := http.Server{Addr: ":9999", Handler: handler}

	defer server.Shutdown(context.Background())
	go server.ListenAndServe()
}

func routeCheck(r *http.Request) ([]byte, error) {
	userId, actions, err := getParams(r)
	if err != nil {
		return nil, err
	}

	permissions := rbac.BulkCheck(userId, actions)
	result, err := json.Marshal(permissions)

	return result, err
}

func getParams(r *http.Request) (int, []string, error) {
	params := r.URL.Query()

	actions := params["actions[]"]
	if actions == nil || len(actions) == 0 {
		return nil, nil, errors.New("Обязательный параметр actions должен быть не пустым массивом.")
	}

	userId, err := strconv.Atoi(params["userId"][0])
	if err != nil {
		return nil, actions, err
	}

	if userId == 0 {
		return nil, actions, errors.New("Обязательный параметр userId должен быть больше нуля.")
	}

	return userId, actions, nil
}