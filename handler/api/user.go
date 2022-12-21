package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type UserAPI interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)
}

type userAPI struct {
	userService service.UserService
}

type successResponse struct {
	User_id int    `json:"user_id"`
	Message string `json:"message"`
}

type successResponseLogout struct {
	Message string `json:"message"`
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	// TODO: answer here
	if user.Email == "" || user.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("email or password is empty"))
		return
	}

	userEntity := &entity.User{
		Email:    user.Email,
		Password: user.Password,
	}

	userID, errService := u.userService.Login(r.Context(), userEntity)

	if errService != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse(fmt.Sprint(errService)))
		return
	}

	successResponse := successResponse{
		User_id: userID,
		Message: "login success",
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(userID),
	})

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(successResponse)

}

func (u *userAPI) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Fullname == "" || user.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("register data is empty"))
		return
	}

	userEntity := &entity.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	userReturn, err := u.userService.Register(r.Context(), userEntity)

	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse(fmt.Sprint(err)))
		return
	}

	successResponse := successResponse{
		User_id: userReturn.ID,
		Message: "register success",
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(successResponse)
}

func (u *userAPI) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	http.SetCookie(w, &http.Cookie{
		Name:  "",
		Value: "",
	})

	response := successResponseLogout{
		Message: "logout success",
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(response)
}

func (u *userAPI) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("user_id is empty"))
		return
	}

	deleteUserId, _ := strconv.Atoi(userId)

	err := u.userService.Delete(r.Context(), int(deleteUserId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "delete success"})
}
