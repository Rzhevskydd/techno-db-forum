package forumDelivery

import (
	"encoding/json"
	"github.com/Rzhevskydd/techno-db-forum/project/app/app"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	"github.com/Rzhevskydd/techno-db-forum/project/app/units/forum/usecase"
	"github.com/Rzhevskydd/techno-db-forum/project/app/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type ApiForumHadler struct {
	ForumUseCase forumUsecase.ForumUseCase
}

func HandleForumRoutes(m *mux.Router, u *app.UseCase) {
	handler := ApiForumHadler{ForumUseCase: forumUsecase.ForumUseCase{}}

	m.HandleFunc("/create", handler.ForumCreateHandler).Methods(http.MethodPost, http.MethodOptions)

}

func (h *ApiForumHadler) ForumCreateHandler(w http.ResponseWriter, r *http.Request) {
	defer utils.CloseBody(w, r)

	in := new(models.Forum)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		utils.NewError(w, http.StatusBadRequest, err.Error())
	}

	h.ForumUseCase.CreateForum(in)

}