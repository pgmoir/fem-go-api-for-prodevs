package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/pgmoir/femGoProject/internal/store"
	"github.com/pgmoir/femGoProject/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(r *registerUserRequest) error {
	if r.Username == "" {
		return errors.New("username is required")
	}

	if len(r.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}

	if r.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
// 	userID, err := utils.ReadIDParam(r)
// 	if err != nil {
// 		uh.logger.Printf("ERROR: readIDParam: %v", err)
// 		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user id"})
// 		return
// 	}

// 	user, err := uh.userStore.GetUserByID(userID)
// 	if err != nil {
// 		uh.logger.Printf("ERROR: getUserByID: %v", err)
// 		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
// }

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: decodingRegisterUser: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	err = uh.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hashingPassword: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: creatingUser: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create user"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}

// func (uh *UserHandler) HandleUpdateUserByID(w http.ResponseWriter, r *http.Request) {
// 	userID, err := utils.ReadIDParam(r)
// 	if err != nil {
// 		uh.logger.Printf("ERROR: readIDParam: %v", err)
// 		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user update id"})
// 		return
// 	}

// 	existingUser, err := uh.userStore.GetUserByID(userID)
// 	if err != nil {
// 		uh.logger.Printf("ERROR: getUserByID: %v", err)
// 		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to fetch user"})
// 		return
// 	}

// 	if existingUser == nil {
// 		uh.logger.Printf("ERROR: getUserByID: %v", err)
// 		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
// 		return
// 	}

// 	var updateUserRequest struct {
// 		Username       *string           `json:"username"`
// 		Email          *string           `json:"email"`
// 		Password       *int              `json:"password"`
// 		CaloriesBurned *int              `json:"calories_burned"`
// 		Bio            []store.UserEntry `json:"bio"`
// 	}

// 	err = json.NewDecoder(r.Body).Decode(&updateUserRequest)
// 	if err != nil {
// 		uh.logger.Printf("ERROR: decodingUpdate: %v", err)
// 		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request"})
// 		return
// 	}

// 	if updateUserRequest.Username != nil {
// 		existingUser.Username = *updateUserRequest.Username
// 	}
// 	if updateUserRequest.Email != nil {
// 		existingUser.Email = *updateUserRequest.Email
// 	}
// 	if updateUserRequest.Password != nil {
// 		existingUser.Password = *updateUserRequest.Password
// 	}
// 	if updateUserRequest.CaloriesBurned != nil {
// 		existingUser.CaloriesBurned = *updateUserRequest.CaloriesBurned
// 	}
// 	if updateUserRequest.Bio != nil {
// 		existingUser.Bio = updateUserRequest.Bio
// 	}

// 	err = uh.userStore.UpdateUser(existingUser)
// 	if err != nil {
// 		uh.logger.Printf("ERROR: updatingUser: %v", err)
// 		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update the user"})
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": existingUser})
// }

// func (uh *UserHandler) HandleDeleteUserByID(w http.ResponseWriter, r *http.Request) {
// 	userID, err := utils.ReadIDParam(r)
// 	if err != nil {
// 		uh.logger.Printf("ERROR: readIDParam: %v", err)
// 		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user update id"})
// 		return
// 	}

// 	err = uh.userStore.DeleteUser(userID)
// 	if err == sql.ErrNoRows {
// 		uh.logger.Printf("ERROR: deletingUser: %v", err)
// 		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user does not exist"})
// 		return
// 	}

// 	if err != nil {
// 		uh.logger.Printf("ERROR: deletingUser: %v", err)
// 		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)

// 	fmt.Fprintf(w, "deleted user id %d\n", userID)
// }
