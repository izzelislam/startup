package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService user.Service
	AuthService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *UserHandler {
	return &UserHandler{
		UserService: userService,
		AuthService: authService,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)

	if err != nil {

		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("failed register user", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.UserService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("failed register user", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.AuthService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("failed register user", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	usserFormatter := user.FormatUser(newUser, token)

	reeponse := helper.APIResponse("User created successfully", http.StatusOK, "success", usserFormatter)
	c.JSON(http.StatusOK, reeponse)
}

func (h *UserHandler) Login(c *gin.Context) {
	var credential user.LoginInput

	err := c.ShouldBindJSON(&credential)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(" login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.UserService.Login(credential)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.AuthService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.APIResponse("login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, token)

	respose := helper.APIResponse("login successfully", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, respose)
}

func (h *UserHandler) CheckEmailAvailablelity(c *gin.Context) {
	var email user.CheckEmailInput
	err := c.ShouldBindJSON(&email)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("check email failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.UserService.IsEmailAvailable(email)
	if err != nil {
		errorMessage := gin.H{"error": "internal server error"}
		response := helper.APIResponse("check email failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	image, err := c.FormFile("avatar")

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := "images/" + image.Filename
	err = c.SaveUploadedFile(image, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	fileExt := filepath.Ext(path)
	uniqueCode := strconv.FormatInt(time.Now().UnixNano(), 10)

	newName := "images/" + uniqueCode + fileExt
	err = os.Rename(path, newName)

	if err != nil {
		os.Remove(path)
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	_, err = h.UserService.SaveAvatar(userID, newName)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"isuploaded": true}
	response := helper.APIResponse("avatar succesfuly uploaded", http.StatusOK, "error", data)
	c.JSON(http.StatusBadRequest, response)

}
