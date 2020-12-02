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

func HandleForumRoutes(m *mux.Router, use *units.UseCase) {
	api := ApiForumHadler{Forum: use.Forum}

	m.HandleFunc("/create", api.ForumCreateHandler).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/{slug}/create", api.ForumThreadCreateHandler).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/{slug}/details", api.ForumDetailsHandler).Methods(http.MethodGet)
	m.HandleFunc("/{slug}/threads", api.ForumGetThreadsHandler).Methods(http.MethodGet)
	m.HandleFunc("/{slug}/users", api.ForumGetUsersHandler).Methods(http.MethodGet)


}

func (api *ApiForumHadler) ForumCreateHandler(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)

	w.Header().Set("Content-Type", "application/json")

	in := new(models.Forum)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		delivery.NewError(w, http.StatusBadRequest, err.Error())
	}

	forum, code := api.Forum.CreateForum(in)

	switch code {
	case 201:
		delivery.ResponseJson(w, http.StatusCreated, forum)
		return
	case 500:
		delivery.NewError(w, http.StatusInternalServerError, "Internal error")
		return
	case 404:
		delivery.NewError(w, http.StatusNotFound, "Can't find user with nickname: " + forum.User)
		return
	case 409:
		delivery.ResponseJson(w, http.StatusConflict, forum)
		return
	}
}

func (api *ApiForumHadler) ForumThreadCreateHandler(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	in := new(models.Thread)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		delivery.NewError(w, http.StatusBadRequest, err.Error())
	}

	slug := vars["slug"]
	newThread, code := api.Forum.CreateForumThread(in, slug)

	switch code {
	case 201:
		delivery.ResponseJson(w, http.StatusCreated, newThread)
		return
	case 500:
		delivery.NewError(w, http.StatusInternalServerError, "Internal error")
		return
	case 404:
		delivery.NewError(w, http.StatusNotFound, "Can't find user with nickname: " + newThread.Author)
		return
	case 409:
		delivery.ResponseJson(w, http.StatusConflict, newThread)
		return
	}
}

func (api *ApiForumHadler) ForumDetailsHandler(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug := vars["slug"]
	forum, err := api.Forum.GetForum(slug)

	if forum == nil || err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find forum with slug: " + slug)
		return
	}

	delivery.ResponseJson(w, http.StatusOK, forum)
}

func (api *ApiForumHadler) ForumGetThreadsHandler(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug := vars["slug"]

	threads, err := api.Forum.GetForumThreads(slug, r.URL.Query())
	if err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find thread with slug: " + slug)
		return
	}

	delivery.ResponseJson(w, http.StatusOK, threads)
}

func (api *ApiForumHadler) ForumGetUsersHandler(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug := vars["slug"]

	users, err := api.Forum.GetForumUsers(slug, r.URL.Query())
	if err != nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find user with slug: " + slug)
		return
	}

	delivery.ResponseJson(w, http.StatusOK, users)
}