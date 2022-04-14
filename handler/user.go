package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{
		UserService: userService,
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

	usserFormatter := user.FormatUser(newUser, "tokentokentokentokentokentoken")

	reeponse := helper.APIResponse("User created successfully", http.StatusOK, "success", usserFormatter)
	c.JSON(http.StatusOK, reeponse)
}
