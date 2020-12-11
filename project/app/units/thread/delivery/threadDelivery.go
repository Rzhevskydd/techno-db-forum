package threadDelivery

import (
	"encoding/json"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	"github.com/Rzhevskydd/techno-db-forum/project/app/units"
	t "github.com/Rzhevskydd/techno-db-forum/project/app/units/thread/usecase"
	"github.com/Rzhevskydd/techno-db-forum/project/app/utils/delivery"
	"github.com/gorilla/mux"
	"net/http"
)

type ApiThreadHandler struct {
	Thread t.ThreadUseCase
}

func HandleThreadRoutes(m *mux.Router, use *units.UseCase) {
	api := ApiThreadHandler{Thread: use.Thread}

	m.HandleFunc("/{slug_or_id}/create", api.HandleCreatePosts).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/{slug_or_id}/details", api.HandleGetThreadDetails ).Methods(http.MethodGet)
	m.HandleFunc("/{slug_or_id}/details", api.HandleUpdateThreadDetails).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/{slug_or_id}/posts", api.HandleGetThreadPosts).Methods(http.MethodGet)
	m.HandleFunc("/{slug_or_id}/vote", api.HandleVoteThread).Methods(http.MethodPost, http.MethodOptions)
}

func (api *ApiThreadHandler) HandleCreatePosts(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	in := new(models.Posts)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		delivery.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	slugOrId := vars["slug_or_id"]
	posts, code, _ := api.Thread.CreateThreadPosts(slugOrId, *in)

	switch code {
	case 201:
		delivery.ResponseJson(w, http.StatusCreated, posts)
		return
	case 500:
		delivery.NewError(w, http.StatusInternalServerError, "Internal error")
		return
	case 404:
		delivery.NewError(w, http.StatusNotFound, "Can't find user with slug_or_id: " + slugOrId)
		return
	case 409:
		delivery.NewError(w, http.StatusConflict, "Parent post was created in another thread")
		return
	}
}

func (api *ApiThreadHandler) HandleGetThreadDetails(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	thread, err := api.Thread.GetThread(slugOrId)

	if err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find user with slug_or_id: " + slugOrId)
		return
	}

	delivery.ResponseJson(w, http.StatusOK, thread)
}

func (api *ApiThreadHandler) HandleUpdateThreadDetails(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]

	in := new(models.ThreadUpdate)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		delivery.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	thread, err := api.Thread.UpdateThread(slugOrId, in)
	if err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find user with slug_or_id: " + slugOrId)
		return
	}

	delivery.ResponseJson(w, http.StatusOK, thread)
}

func (api *ApiThreadHandler) HandleVoteThread(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]

	in := new(models.Vote)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		delivery.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	thread, err := api.Thread.VoteThread(slugOrId, in)
	if err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find user with slug_or_id: " + slugOrId)
		return
	}

	delivery.ResponseJson(w, http.StatusOK, thread)
}

func (api *ApiThreadHandler) HandleGetThreadPosts(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]

	posts, err := api.Thread.GetThreadPosts(slugOrId, r.URL.Query())
	if err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find user with slug_or_id: " + slugOrId)
		return
	}

	delivery.ResponseJson(w, http.StatusOK, posts)
}