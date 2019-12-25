package application

import (
	"net/http"

	"private-conda-repo/store/models"
)

func ListUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := db.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, user := range users {
		user.Password = ""
	}

	toJson(w, &users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	u, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := u.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err = db.AddUser(u.Name, u.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u.Password = ""

	if _, err := repo.CreateChannel(u.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	toJson(w, &u)
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	u, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.RemoveUser(u.Name, u.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repo.RemoveChannel(u.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok(w)
}

func CheckUser(w http.ResponseWriter, r *http.Request) {
	u, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actual, err := db.GetUser(u.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if actual.HasValidPassword(u.Name) {
		ok(w)
	} else {
		http.Error(w, "invalid credentials", http.StatusForbidden)
	}
}

func getUser(r *http.Request) (*models.User, error) {
	var u *models.User
	if err := readJson(r, &u); err != nil {
		return nil, err
	}
	return u, nil
}
