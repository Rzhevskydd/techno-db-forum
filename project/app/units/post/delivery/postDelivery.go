package postDelivery

import (
	"encoding/json"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	"github.com/Rzhevskydd/techno-db-forum/project/app/units"
	p "github.com/Rzhevskydd/techno-db-forum/project/app/units/post/usecase"
	"github.com/Rzhevskydd/techno-db-forum/project/app/utils/delivery"
	"github.com/gorilla/mux"
	"net/http"
)

type ApiPostHadler struct {
	Post p.PostUseCase
}

func HandlePostRoutes(m *mux.Router, use *units.UseCase) {
	api := ApiPostHadler{Post: use.Post}

	m.HandleFunc("/{id}/details", api.PostGetDetailsHandler).Methods(http.MethodGet)
	m.HandleFunc("/{id}/details", api.PostPostDetailsHandler).Methods(http.MethodPost, http.MethodOptions)

}

func (api *ApiPostHadler) PostGetDetailsHandler(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	postFull, err := api.Post.GetPostDetails(id , r.URL.Query())

	if postFull == nil || err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find post with id: " + id)
		return
	}

	delivery.ResponseJson(w, http.StatusOK, postFull)
}

func (api *ApiPostHadler) PostPostDetailsHandler(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	in := new(models.Post)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		delivery.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := api.Post.UpdatePostDetails(id, in.Message)
	if post == nil || err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find post with id: " + id)
		return
	}

	delivery.ResponseJson(w, http.StatusOK, post)
}