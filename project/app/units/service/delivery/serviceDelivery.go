package serviceDelivery

import (
	"github.com/Rzhevskydd/techno-db-forum/project/app/units"
	s "github.com/Rzhevskydd/techno-db-forum/project/app/units/service/usecase"
	"github.com/Rzhevskydd/techno-db-forum/project/app/utils/delivery"
	"github.com/gorilla/mux"
	"net/http"
)

type ApiServiceHandler struct {
	Service s.ServiceUseCase
}

func HandleServiceRoutes(m *mux.Router, use *units.UseCase) {
	api := ApiServiceHandler{Service: use.Service}

	m.HandleFunc("/clear", api.ServiceClearHandler).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/status", api.ServiceGetStatusHandler ).Methods(http.MethodGet)

}

func (api *ApiServiceHandler) ServiceClearHandler(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)

	w.Header().Set("Content-Type", "application/json")

	err := api.Service.Clear()
	if err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find user with nickname: ")
		return
	}

	delivery.ResponseJson(w, http.StatusOK, nil)
}

func (api *ApiServiceHandler) ServiceGetStatusHandler(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")


	status, err := api.Service.GetStatus()
	if status == nil || err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find user with nickname: ")
		return
	}

	delivery.ResponseJson(w, http.StatusOK, status)
}