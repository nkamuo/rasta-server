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
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func FindRespondents(c *gin.Context) {
	var respondents []model.Respondent
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.DB.Preload("User").Preload("Place").Scopes(pagination.Paginate(respondents, &page, model.DB)).Find(&respondents).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = respondents
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func FindRespondentsByCompany(c *gin.Context) {

	companyService := service.GetCompanyService()

	companyID, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid companyId provided"})
		return
	}

	if _, err := companyService.GetById(companyID); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return

	}

	var respondents []model.Respondent
	if err = model.DB.
		Joins("JOIN companies ON companies.id = respondents.company_id").
		Where("companies.id = ?", companyID).
		Preload("User").
		// Preload("Company").
		Find(&respondents).Error; nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if respondents == nil {
		respondents = make([]model.Respondent, 0)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondents})
}

func AddRespondentToCompany(c *gin.Context) {

	companyService := service.GetCompanyService()
	respondentService := service.GetRespondentService()

	companyID, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid companyId provided"})
		return
	}

	company, err := companyService.GetById(companyID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var input dto.RespondentCompanyAssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	respondant, err := respondentService.GetById(input.RespondentID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := respondentService.AssignToCompany(respondant, company); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Respondent assinged to company successfully"})
}

func RemoveRespondentFromCompany(c *gin.Context) {

	companyService := service.GetCompanyService()
	respondentService := service.GetRespondentService()

	companyID, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid companyId provided"})
		return
	}

	respondentID, err := uuid.Parse(c.Param("respondent_id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid companyId provided"})
		return
	}

	company, err := companyService.GetById(companyID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	respondant, err := respondentService.GetById(respondentID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := respondentService.RemoveFromCompany(respondant, company); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": "success", "message": "Respondent removed from company successfully"})
}

func CreateRespondent(c *gin.Context) {

	userService := service.GetUserService()
	placeService := service.GetPlaceService()
	companyService := service.GetCompanyService()
	vehicleService := service.GetVehicleService()
	respondentService := service.GetRespondentService()
	respondentRepo := repository.GetRespondentRepository()
	// Validate input
	var input dto.RespondentCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := userService.GetById(input.UserID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if _, err := respondentRepo.GetByUser(*user); nil == err {
		message := fmt.Sprintf("There is already a respondant account \"%s\" associated with this user", user.FullName())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return

	}

	// Create respondent
	respondent := model.Respondent{
		UserID: &user.ID,
	}

	if nil != input.VehicleID {
		if _, err := vehicleService.GetById(*input.VehicleID); nil != err {
			message := fmt.Sprintf("Could not find vehicle with [id:%s]: %s", input.VehicleID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.VehicleID = input.VehicleID
	}

	if nil != input.CompanyID {
		if _, err := companyService.GetById(*input.CompanyID); nil != err {
			message := fmt.Sprintf("Could not find company with [id:%s]: %s", input.CompanyID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.CompanyID = input.CompanyID
	}

	if nil != input.PlaceID {
		if _, err := placeService.GetById(*input.PlaceID); nil != err {
			message := fmt.Sprintf("Could not find place with [id:%s]: %s", input.PlaceID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.PlaceID = input.PlaceID
	}

	if nil != input.Active {
		respondent.Active = input.Active
	}

	// fmt.Printf("Input USer ID: %s\n user.ID: %s\n respondent.UserId: %s\n", input.UserId, user.ID, respondent.UserID)

	if err := respondentService.Save(&respondent); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondent})
}

func FindRespondent(c *gin.Context) {
	respondentService := service.GetRespondentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Access denied: %s", "You may not access this resource")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}
	// respondent, err := respondentService.GetById(id)
	// var respondent model.Respondent
	respondent, err := respondentService.GetById(id, "User", "AccessBalance", "AccessSubscription", "Place", "Company", "Documents")
	if nil != err {
		message := fmt.Sprintf("Could not find respondent with [id:%s]: %s", id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if !*rUser.IsAdmin && rUser.ID.String() != respondent.UserID.String() {

		message := fmt.Sprintf("Access denied: %s", "You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}
	model.ResolveDocumentSlicePublicPaths(respondent.Documents)

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondent})
}

func GetCurrentRespondent(c *gin.Context) {
	respondentRepo := repository.GetRespondentRepository()

	user, err := auth.GetCurrentUser(c)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	respondent, err := respondentRepo.GetByUser(*user, "User", "Place", "Vehicle", "Company", "AccessBalance", "AccessSubscription", "Documents")
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	model.ResolveDocumentSlicePublicPaths(respondent.Documents)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondent})

}

func UpdateRespondent(c *gin.Context) {
	placeService := service.GetPlaceService()
	companyService := service.GetCompanyService()
	vehicleService := service.GetVehicleService()
	respondentService := service.GetRespondentService()
	respondentRepo := repository.GetRespondentRepository()

	// Validate input
	var input dto.RespondentUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	respondent, err := respondentService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find respondent with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if nil != input.VehicleID {
		if _, err := vehicleService.GetById(*input.VehicleID); nil != err {
			message := fmt.Sprintf("Could not find vehicle with [id:%s]: %s", input.VehicleID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.VehicleID = input.VehicleID
	}

	if nil != input.CompanyID {
		if _, err := companyService.GetById(*input.CompanyID); nil != err {
			message := fmt.Sprintf("Could not find company with [id:%s]: %s", input.CompanyID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.CompanyID = input.CompanyID
	}

	if nil != input.PlaceID {
		if _, err := placeService.GetById(*input.PlaceID); nil != err {
			message := fmt.Sprintf("Could not find place with [id:%s]: %s", input.PlaceID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.PlaceID = input.PlaceID
	}

	if nil != input.Active {
		respondent.Active = input.Active
	}

	if err := respondentService.Save(respondent); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	nRespondent, err := respondentRepo.GetById(respondent.ID, "Place", "User", "Vehicle")
	if nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated respondent \"%s\"", respondent.ID)
	c.JSON(http.StatusOK, gin.H{"data": nRespondent, "status": "success", "message": message})
}

func DeleteRespondent(c *gin.Context) {
	respondentService := service.GetRespondentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	respondent, err := respondentService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find respondent with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := respondentService.Delete(respondent); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted respondent \"%s\"", respondent.ID)
	c.JSON(http.StatusOK, gin.H{"data": respondent, "status": "success", "message": message})
}

func FindRespondentDocuments(c *gin.Context) {
	respondentService := service.GetRespondentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	docType := c.Param("type")
	if docType == "" {
		docType = c.Query("type")
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Access denied: %s", "You may not access this resource")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}
	// respondent, err := respondentService.GetById(id)
	// var respondent model.Respondent
	respondent, err := respondentService.GetById(id, "User", "AccessBalance", "AccessSubscription", "Place", "Company")
	if nil != err {
		message := fmt.Sprintf("Could not find respondent with [id:%s]: %s", id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if !*rUser.IsAdmin && rUser.ID.String() != respondent.UserID.String() {

		message := fmt.Sprintf("Access denied: %s", "You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	var page pagination.Page
	var documents []*model.ImageDocument

	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Where("responder_id = ?", respondent.ID)
	if docType != "" {
		query = query.Where("doc_type = ?", docType)
	}

	query = query.Scopes(pagination.Paginate(documents, &page, query))

	if err = query.Find(&documents).Error; err != nil {
		message := fmt.Sprintf("An error occured: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
	}

	model.ResolveDocumentSlicePublicPaths(&documents)

	page.Rows = documents

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": page})
}

func UpdateRespondentDocuments(c *gin.Context) {
	respondentService := service.GetRespondentService()

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

	// docType := c.Param("type")
	// if docType == "" {
	// 	docType = c.Query("type")
	// }

	var input dto.RespondentDocumentVerificationInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Access denied: %s", "You may not access this resource")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}
	// respondent, err := respondentService.GetById(id)
	// var respondent model.Respondent
	respondent, err := respondentService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find respondent with [id:%s]: %s", id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if !*rUser.IsAdmin && rUser.ID.String() != respondent.UserID.String() {
		message := fmt.Sprintf("Access denied: %s", "You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	uploadDir := config.UPLOAD_DIR
	if config.RESPONDER_DOCUMENT_UPLOAD_DIR != "" {
		uploadDir = config.RESPONDER_DOCUMENT_UPLOAD_DIR
	}
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	var documents []model.ImageDocument

	// Multipart form
	form, _ := c.MultipartForm()
	if form != nil {
		files := form.File["documents[]"]
		if len(files) != 0 {

			for _, file := range files {
				// log.Println(file.Filename)
				ext := filepath.Ext(file.Filename)
				newName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

				uploadPath := fmt.Sprintf("%s/responder/%s/documents/%s", uploadDir, respondent.ID, newName)
				dst := fmt.Sprintf("%s/%s", config.ASSET_DIR, uploadPath)
				// Upload the file to specific dst.
				err = c.SaveUploadedFile(file, dst)
				if err != nil {
					message := fmt.Sprintf("An error occurred uploading file: %s", err.Error())
					c.JSON(http.StatusOK, gin.H{"message": message, "status": "error"})
					return
				}

				docType := "DRIVER_LICENSE" //file.Header.Get("docType")

				document := model.ImageDocument{
					ResponderID:  &respondent.ID,
					DocType:      &docType,
					FilePath:     uploadPath,
					Size:         file.Size,
					OriginalName: file.Filename,
					Extension:    ext,
				}
				documents = append(documents, document)
			}
		}
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if !*rUser.IsAdmin {
			*respondent.Active = false
		}
		if input.Ssn != nil {
			respondent.Ssn = input.Ssn
			tx.Save(respondent)
		}
		// CLEAR EXISTING DOCUMENTS

		// if err := tx.Where("responder_id = ?", respondent.ID).Delete(&model.ImageDocument{}).Error; nil != err {
		// 	return err
		// }
		if len(documents) > 0 {
			if err = respondent.ClearDocuments(); err != nil {
				return err
			}
			for _, document := range documents {
				if err := tx.Create(&document).Error; nil != err {
					return err
				}
			}
		}
		return nil
	})

	if nil != err {
		message := fmt.Sprintf("An error occured: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondent})

	// c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))

	// var page pagination.Page

	// query := model.DB.Where("responder_id = ?", respondent.ID)
	// if docType != "" {
	// 	query = query.Where("doc_type = ?", docType)
	// }

	// query = query.Scopes(pagination.Paginate(documents, &page, query))

	// if err = query.Find(&documents).Error; err != nil {
	// 	message := fmt.Sprintf("An error occured: %s", err.Error())
	// 	c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
	// }

	// page.Rows = documents

}
