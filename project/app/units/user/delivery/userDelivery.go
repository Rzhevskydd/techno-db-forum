package userDelivery

import (
	"encoding/json"
	"github.com/Rzhevskydd/techno-db-forum/project/app/app"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	u "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/usecase"
	"github.com/Rzhevskydd/techno-db-forum/project/app/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type ApiUserHandler struct {
	User u.UserUseCase
}

func HandleUserRoutes(m *mux.Router, use *app.UseCase) {
	api := ApiUserHandler{User: use.User}

	m.HandleFunc("/{nickname}/create", api.HandleCreateUser).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/{nickname}/profile", api.HandleGetUser ).Methods(http.MethodGet)
	m.HandleFunc("/{nickname}/profile", api.HandleUpdateUser).Methods(http.MethodPost, http.MethodOptions)

}

func (api *ApiUserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer utils.CloseBody(w, r)

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	in := new(models.User)
	in.Nickname = vars["nickname"]
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		utils.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := api.User.CreateUser(in)

	if err != nil {
		return
	}

	if len(users) > 1 {
		utils.ResponseJson(w, http.StatusConflict, users)
	} else if len(users) == 1 {
		utils.ResponseJson(w, http.StatusCreated, in)
	}

}

func (api *ApiUserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	defer utils.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user, err := api.User.GetUser(nickname)
	if user == nil || err != nil {
		utils.NewError(w, http.StatusNotFound, "Can't find user with nickname: " + nickname)
		return
	}

	utils.ResponseJson(w, http.StatusOK, user)
}

func (api *ApiUserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	defer utils.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	in := new(models.User)
	in.Nickname = vars["nickname"]
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		utils.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.User.UpdateUser(in)

	if user == nil && err != nil {
		utils.NewError(w, http.StatusNotFound, "Can't find user with nickname: " + vars["nickname"])
		return
	}

	if user != nil && err != nil {
		utils.ResponseJson(w, http.StatusConflict, user)
		return
	}

	utils.ResponseJson(w, http.StatusOK, user)
}