package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PeterAkaJohn/SightSeeing/src/converters"
	"github.com/PeterAkaJohn/SightSeeing/src/model"
	"github.com/gorilla/sessions"
)

type userController struct {
}

func (uc *userController) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = model.Register(user.Username, user.Password, user.FirstName, user.LastName, user.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
}

func (uc *userController) Login(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "loginSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userLog model.User
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&userLog)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	user, err := model.Login(userLog.Username, userLog.Password)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 400)
		return
	}
	//userVM will only be used in profile page, using it now only for debugging purposes
	userVM := converters.ConvertUserToUserVM(*user)
	session.Values["userID"] = user.ID
	session.Save(r, w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userVM)
}

func (uc *userController) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "loginSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Options = &sessions.Options{
		MaxAge: -1,
	}

	session.Save(r, w)
}
