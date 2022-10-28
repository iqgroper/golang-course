package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Logger   *logrus.Entry
	UserRepo user.UserRepo
	Sessions *session.SessionsManager
}

type NewUser struct {
	Username string
	Password string
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Println("error reading body in Login")
	}

	newUser := &NewUser{}
	errorUnmarshal := json.Unmarshal(body, newUser)
	if errorUnmarshal != nil {
		h.Logger.Println("error unmarshling in Login: ", errorUnmarshal.Error())
		http.Error(w, fmt.Sprintf(`Bad login data format: %s`, errorUnmarshal.Error()), http.StatusBadRequest)
		return
	}

	userToLogIn, err := h.UserRepo.Authorize(newUser.Username, newUser.Password)
	if err == user.ErrNoUser {
		http.Error(w, `no user`, http.StatusBadRequest)
		return
	}
	if err == user.ErrBadPass {
		http.Error(w, `bad password`, http.StatusBadRequest)
		return
	}

	sess, _ := h.Sessions.Create(w, userToLogIn)
	h.Logger.WithFields(logrus.Fields{
		"user":   userToLogIn.Login,
		"SessID": sess.ID,
	}).Info("Login")

	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, sess.ID)))
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Println("error reading body in Register")
	}
	h.Logger.WithFields(logrus.Fields{
		"method": r.Method,
		"body":   string(body),
	}).Info("Register")

	newUser := &NewUser{}
	errorUnmarshal := json.Unmarshal(body, newUser)
	if errorUnmarshal != nil {
		h.Logger.Println("error unmarshling in Register: ", errorUnmarshal.Error())
		http.Error(w, fmt.Sprintf(`Bad register data format: %s`, errorUnmarshal.Error()), http.StatusBadRequest)
		return
	}

	createdUser, errRegister := h.UserRepo.Register(newUser.Username, newUser.Password)
	if errRegister != nil {
		h.Logger.Println("error registring user: ", errRegister)
		http.Error(w, fmt.Sprintf(`error registring user: %s`, errRegister.Error()), http.StatusBadRequest)
		return
	}

	newSession, SessErr := h.Sessions.Create(w, createdUser)
	if SessErr != nil {
		h.Logger.Println("error creating session: ", SessErr.Error())
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, newSession.ID)))
}

// func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
// 	h.Sessions.DestroyCurrent(w, r)
// 	http.Redirect(w, r, "/", 302)
// }
