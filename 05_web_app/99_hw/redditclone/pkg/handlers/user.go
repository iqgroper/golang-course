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

func getNameAndPass(r *http.Request) (*user.NewUser, error) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading body")
		return nil, err
	}

	newUser := &user.NewUser{}
	errorUnmarshal := json.Unmarshal(body, newUser)
	if errorUnmarshal != nil {
		fmt.Println("error unmarshling: ", errorUnmarshal.Error())
		return nil, errorUnmarshal
	}
	return newUser, nil
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	newUser, errGetting := getNameAndPass(r)
	if errGetting != nil {
		h.Logger.Println("error in Login:", errGetting.Error())
		http.Error(w, fmt.Sprintf(`Bad Login data format: %s`, errGetting.Error()), http.StatusBadRequest)
		return
	}

	userToLogIn, err := h.UserRepo.Authorize(newUser.Username, newUser.Password)
	if err == user.ErrNoUser {
		h.Logger.Println("error in Login:", err.Error())
		http.Error(w, `{"message":"user not found"}`, http.StatusUnauthorized)
		return
	}
	if err == user.ErrBadPass {
		h.Logger.Println("error in Login:", err.Error())
		http.Error(w, `{"message": "invalid password"}`, http.StatusUnauthorized)
		return
	}

	sess, _ := h.Sessions.Create(w, userToLogIn)
	h.Logger.WithFields(logrus.Fields{
		"user": userToLogIn.Login,
		"task": "logged in",
	}).Info("Login")

	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, sess.ID)))
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	newUser, errGetting := getNameAndPass(r)
	if errGetting != nil {
		h.Logger.Println("error in Login:", errGetting.Error())
		http.Error(w, fmt.Sprintf(`Bad Login data format: %s`, errGetting.Error()), http.StatusBadRequest)
		return
	}

	h.Logger.WithFields(logrus.Fields{
		"method": r.Method,
		"body":   fmt.Sprintf(`username: %s, password: %s`, newUser.Username, newUser.Password),
	}).Info("Register")

	createdUser, errRegister := h.UserRepo.Register(newUser.Username, newUser.Password)
	if errRegister == user.ErrUserExists {
		h.Logger.Println("error registring user: ", errRegister.Error())
		http.Error(w, fmt.Sprintf(`{"errors":[{"location":"body","param":"username","value":"%s","msg":"%s"}]}`, newUser.Username, errRegister.Error()), http.StatusUnprocessableEntity)
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
