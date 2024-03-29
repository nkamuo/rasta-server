package controller

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
)

// GET /users
// Get all users
func FindUsers(c *gin.Context) {
	var users []model.User
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	query := model.DB

	if search := c.Query("search"); search != "" {
		like := fmt.Sprintf("%%%s%%", search)
		query = query.Where("users.first_name LIKE ? OR users.last_name LIKE ?", like, like)
	}

	if err := query.Scopes(pagination.Paginate(users, &page, query)).Find(&users).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = users
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateUser(c *gin.Context) {

	userService := service.GetUserService()
	// Validate input
	var input dto.UserCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	// Create user
	user := model.User{
		Email: input.Email,
		// PERSONAL INFOR
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		IsAdmin:   &input.IsAdmin,
	}

	userService.HashUserPassword(&user, input.Password)
	if err := userService.Save(&user); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func FindUser(c *gin.Context) {
	userService := service.GetUserService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	user, err := userService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find user with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdateUserAvatar(c *gin.Context) {
	userService := service.GetUserService()

	config, err := initializers.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	user, err := userService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find user with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Uathentication error: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}

	if !*rUser.IsAdmin && rUser.ID != user.ID {
		message := fmt.Sprintf("You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	// var err error;

	// single file
	file, err := c.FormFile("file")

	if nil != err {
		message := fmt.Sprintf("An error occurred: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{"message": message, "status": "error"})
		return
	}

	uploadDir := config.UPLOAD_DIR
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	if config.USER_AVATAR_UPLOAD_DIR != "" {
		uploadDir = config.USER_AVATAR_UPLOAD_DIR
	}

	// log.Println(file.Filename)
	ext := filepath.Ext(file.Filename)

	uploadPath := fmt.Sprintf("%s/users/%s/avatar%s", uploadDir, user.ID, ext)
	dst := fmt.Sprintf("%s/%s", config.ASSET_DIR, uploadPath)

	// Upload the file to specific dst.
	err = c.SaveUploadedFile(file, dst)

	if err != nil {
		message := fmt.Sprintf("An error occurred uploading file: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{"message": message, "status": "error"})
		return
	}

	user.AvatarPath = &uploadPath

	if err = userService.Save(user); err != nil {
		message := fmt.Sprintf("An error occurred uploading file: %s", err.Error())
		c.JSON(http.StatusOK, gin.H{"message": message, "status": "error"})
		return
	}

	// if err != nil
	{
		// message := fmt.Sprintf("'%s' uploaded!", file.Filename)
		c.JSON(http.StatusOK, gin.H{"data": user, "status": "success"})
	}
	// c.JSON(http.StatusOK, gin.H{"data": user, "status": "success"})
}

func UpdateUser(c *gin.Context) {
	userService := service.GetUserService()

	// Validate input
	var input dto.UserUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	user, err := userService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find user with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.FirstName {
		user.FirstName = *input.FirstName
	}

	if nil != input.LastName {
		user.LastName = *input.LastName
	}

	if nil != input.Email {
		user.Email = *input.Email
	}

	if nil != input.Phone {
		user.Phone = *input.Phone
	}

	if nil != input.IsAdmin {
		user.IsAdmin = input.IsAdmin
	}

	if nil != input.Password {
		userService.HashUserPassword(user, *input.Password)
	}

	userService.UpdateStripeCustomer(user, false)
	if err := userService.Save(user); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated user \"%s\"", user.FullName())
	c.JSON(http.StatusOK, gin.H{"data": user, "status": "success", "message": message})
}

func DeleteUser(c *gin.Context) {
	userService := service.GetUserService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("You may not access this resource ")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}
	user, err := userService.GetById(id)

	if !*rUser.IsAdmin {
		message := fmt.Sprintf("You may not access this resource ")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	if nil != err {
		message := fmt.Sprintf("Could not find user with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if err := userService.Delete(user); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted user \"%s\"", user.FullName())
	c.JSON(http.StatusOK, gin.H{"data": user, "status": "success", "message": message})
}
