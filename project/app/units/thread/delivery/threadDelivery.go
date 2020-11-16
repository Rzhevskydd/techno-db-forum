package threadDelivery

import (
	"encoding/json"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	"github.com/Rzhevskydd/techno-db-forum/project/app/units"
	u "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/usecase"
	"github.com/Rzhevskydd/techno-db-forum/project/app/utils/delivery"
	"github.com/Rzhevskydd/techno-db-forum/project/app/utils/validator"
	"github.com/gorilla/mux"
	"net/http"
)

type ApiUserHandler struct {
	User u.UserUseCase
}

func HandleUserRoutes(m *mux.Router, use *units.UseCase) {
	api := ApiUserHandler{User: use.User}

	m.HandleFunc("/{nickname}/create", api.HandleCreateUser).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/{nickname}/profile", api.HandleGetUser ).Methods(http.MethodGet)
	m.HandleFunc("/{nickname}/profile", api.HandleUpdateUser).Methods(http.MethodPost, http.MethodOptions)

}

func (api *ApiUserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	in := new(models.User)
	in.Nickname = vars["nickname"]
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		delivery.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := api.User.CreateUser(in)

	if err != nil {
		delivery.NewError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if users != nil && len(users) > 0 {
		delivery.ResponseJson(w, http.StatusConflict, users)
		return
	}

	delivery.ResponseJson(w, http.StatusCreated, in)
}

func (api *ApiUserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user, err := api.User.GetUser(nickname)
	if err != nil {
		delivery.NewError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find user with nickname: " + nickname)
		return
	}

	delivery.ResponseJson(w, http.StatusOK, user)
}

func (api *ApiUserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	defer delivery.CloseBody(w, r)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	in := new(models.User)
	in.Nickname = vars["nickname"]
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		delivery.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.User.UpdateUser(in)

	if err != nil {
		if user == nil && !validator.IsEmpty(user.Nickname) {
			delivery.NewError(w, http.StatusInternalServerError, err.Error())
			return
		}
		delivery.ResponseJson(w, http.StatusConflict, user)
		return
	}

	if user == nil {
		delivery.NewError(w, http.StatusNotFound, "Can't find user with nickname: " + vars["nickname"])
		return
	}

	delivery.ResponseJson(w, http.StatusOK, user)
}