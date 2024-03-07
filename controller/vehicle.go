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
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// GET /vehicles
// Get all vehicles
func FindVehicles(c *gin.Context) {
	var vehicles []model.Vehicle
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if nil != err {
		message := fmt.Sprintf("Authentication Error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	query := model.DB.Model(&model.Vehicle{}).Preload("Documents").Preload("Model").Preload("Owner")

	if *rUser.IsAdmin {
		if ownerID := c.Query("owner_id"); ownerID != "" {
			query = query.Where("owner_id = ?", ownerID)
		}
	} else {
		query = query.Where("owner_id = ?", rUser.ID)
	}

	query = query.Scopes(pagination.Paginate(vehicles, &page, query))

	if err := query.Find(&vehicles).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	for i, vehicle := range vehicles {
		model.ResolveDocumentSlicePublicPaths(vehicle.Documents)
		vehicles[i] = vehicle
	}

	page.Rows = vehicles
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateVehicle(c *gin.Context) {

	userService := service.GetUserService()
	companyService := service.GetCompanyService()
	vehicleService := service.GetVehicleService()
	modelService := service.GetVehicleModelService()

	var input dto.VehicleCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var makeName, modelName *string
	var vehicleModel *model.VehicleModel
	var owner *model.User
	var company *model.Company

	if input.ModelID != nil {
		vehicleModel, err = modelService.GetById(*input.ModelID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified Model with [id:%s]: %s", input.ModelID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
	} else {
		if input.MakeName == nil || input.ModelName == nil {
			message := fmt.Sprintf("You must provide a valid %s or both %s and %s", "ModelID", "MakeName", "ModelName")
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
		makeName, modelName = input.MakeName, input.ModelName
	}

	if input.OwnerID != nil {
		owner, err = userService.GetById(*input.OwnerID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified User with [id:%s]: %s", input.OwnerID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
		if input.OwnerID.String() != rUser.ID.String() && !*rUser.IsAdmin {
			message := "You are not allowed to provide {ownerId} for this request"
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
			return
		}
	} else {
		if input.CompanyID != nil {
			company, err = companyService.GetById(*input.CompanyID)
			if nil != err {
				message := fmt.Sprintf("Could not resolve the specified Company with [id:%s]: %s", input.OwnerID, err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
				return
			}
		} else {
			owner = rUser
		}

	}

	vehicle := model.Vehicle{
		LicensePlateNumber: *input.LicensePlateNumber,
		Color:              *input.Color,
		Description:        input.Description,
		Published:          &input.Published,
		VinNumber:          input.VinNumber,
	}

	if owner != nil {
		vehicle.OwnerID = &owner.ID
	}

	if company != nil {
		vehicle.CompanyID = &company.ID
	}
	if vehicleModel != nil {
		vehicle.ModelID = &vehicleModel.ID
	} else {
		vehicle.MakeName, vehicle.ModelName = makeName, modelName
	}

	if err := vehicleService.Save(&vehicle); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err = UpdateVehicleDocuments(c, &vehicle, rUser); err != nil {
		message := fmt.Sprintf("An error occured: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	model.ResolveDocumentSlicePublicPaths(vehicle.Documents)
	c.JSON(http.StatusOK, gin.H{"data": vehicle, "status": "success"})
}

func FindVehicle(c *gin.Context) {
	vehicleService := service.GetVehicleService()

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	// Preload("Respondent.User").Preload("Assignments").Preload("Assignments.Assignment.Product")
	vehicle, err := vehicleService.GetById(id, "Owner", "Documents")
	if nil != err {
		message := fmt.Sprintf("Could not find vehicle with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if *rUser.IsAdmin {
		/// USER IS ADMIN - ADMIN CAN VIEW ANY SESSION
	} else {
		// USER IS NOT ADMIN - USER CAN ONLY VIEW HIS/HER VEHICLE
		if rUser.ID.String() != (*vehicle.OwnerID).String() {
			message := fmt.Sprintf("Could not find your vehicle with [id:%s]", id)
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
			return
		}
	}
	model.ResolveDocumentSlicePublicPaths(vehicle.Documents)

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": vehicle})
}

func UpdateVehicle(c *gin.Context) {
	userService := service.GetUserService()
	vehicleService := service.GetVehicleService()
	modelService := service.GetVehicleModelService()

	rUser, err := auth.GetCurrentUser(c)

	if err != nil {
		message := fmt.Sprintf("Authentication Error: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}

	var input dto.VehicleUpdateInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	vehicle, err := vehicleService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find vehicle with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if !*rUser.IsAdmin && (vehicle.OwnerID == nil || rUser.ID.String() != vehicle.OwnerID.String()) {
		message := fmt.Sprintf("Unathorized: You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	if nil != input.VinNumber {
		vehicle.VinNumber = input.VinNumber
	}

	if nil != input.Color {
		vehicle.Color = *input.Color
	}

	if nil != input.ModelName {
		vehicle.ModelName = input.ModelName
	}

	if nil != input.MakeName {
		vehicle.MakeName = input.MakeName
	}

	if nil != input.LicensePlateNumber {
		vehicle.LicensePlateNumber = *input.LicensePlateNumber
	}
	if nil != input.Description {
		vehicle.Description = *input.Description
	}

	if nil != input.ModelID {
		vehicleModel, err := modelService.GetById(*input.ModelID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified Model with [id:%s]: %s", input.ModelID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
		vehicle.ModelID = &vehicleModel.ID
	}

	if *rUser.IsAdmin {
		if nil != input.OwnerID {
			owner, err := userService.GetById(*input.OwnerID)
			if nil != err {
				message := fmt.Sprintf("Could not resolve the specified User with [id:%s]: %s", input.OwnerID, err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
				return
			}
			vehicle.OwnerID = &owner.ID
		}

		if nil != input.Published {
			vehicle.Published = input.Published
		}

	} else {
		*vehicle.Published = false
		/**
		* This is to automatically disable each vehicle after user modification so that admins
		* can review and manually re-enable them.
		 */
	}

	if err := vehicleService.Save(vehicle); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err = UpdateVehicleDocuments(c, vehicle, rUser); err != nil {
		message := fmt.Sprintf("An error occured while updating Images: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	model.ResolveDocumentSlicePublicPaths(vehicle.Documents)

	message := fmt.Sprintf("Updated model \"%s\"", vehicle.ID)
	c.JSON(http.StatusOK, gin.H{"data": vehicle, "status": "success", "message": message})
}

func DeleteVehicle(c *gin.Context) {
	vehicleService := service.GetVehicleService()

	id := c.Param("id")

	var vehicle model.Vehicle

	if err := model.DB.Where("id = ?", id).First(&vehicle).Error; err != nil {
		message := fmt.Sprintf("Could not find vehicle with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if err := vehicleService.Delete(&vehicle); nil != err {
		message := fmt.Sprintf("An error occurred while deleting entry:\"%s\"", err.Error())
		c.JSON(http.StatusOK, gin.H{"data": vehicle, "status": "success", "message": message})
		return
	}

	vehicle.ClearDocuments()

	message := fmt.Sprintf("Deleted vehicle \"%s\"", vehicle.ID)
	c.JSON(http.StatusOK, gin.H{"data": vehicle, "status": "success", "message": message})
}

func FindVehicleDocuments(c *gin.Context) {
	vehicleService := service.GetVehicleService()

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

	vehicle, err := vehicleService.GetById(id, "Owner")
	if nil != err {
		message := fmt.Sprintf("Could not find respondent with [id:%s]: %s", id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if !*rUser.IsAdmin && rUser.ID.String() != vehicle.OwnerID.String() {

		message := fmt.Sprintf("Access denied: %s", "You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	var page pagination.Page
	var documents []*model.ImageDocument

	query := model.DB.Where("vehicle_id = ?", vehicle.ID)
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

func UpdateVehicleDocuments(c *gin.Context, vehicle *model.Vehicle, rUser *model.User) (err error) {
	// respondentService := service.GetRespondentService()

	config, err := initializers.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	uploadDir := config.UPLOAD_DIR
	if config.VEHICLE_DOCUMENT_UPLOAD_DIR != "" {
		uploadDir = config.VEHICLE_DOCUMENT_UPLOAD_DIR
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

				uploadPath := fmt.Sprintf("%s/vehicles/%s/documents/%s", uploadDir, vehicle.ID, newName)
				dst := fmt.Sprintf("%s/%s", config.ASSET_DIR, uploadPath)
				// Upload the file to specific dst.
				err = c.SaveUploadedFile(file, dst)
				if err != nil {
					return err
				}

				docType := "IMAGE" //file.Header.Get("docType")

				document := model.ImageDocument{
					VehicleID:    &vehicle.ID,
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
			V := false
			vehicle.Published = &V
		}
		// CLEAR EXISTING DOCUMENTS

		// if err := tx.Where("responder_id = ?", respondent.ID).Delete(&model.ImageDocument{}).Error; nil != err {
		// 	return err
		// }
		if err = vehicle.ClearDocuments(); err != nil {
			return err
		}
		for _, document := range documents {
			if err := tx.Create(&document).Error; nil != err {
				return err
			}
		}
		return nil
	})

	return err
}

func ValidateVehicleCategory(category model.VehicleCategory) (err error) {
	switch category {
	case model.PRODUCT_FLAT_TIRE_SERVICE:
		return nil
	case model.PRODUCT_FUEL_DELIVERY_SERVICE:
		return nil
	case model.PRODUCT_TIRE_AIR_SERVICE:
		return nil
	case model.PRODUCT_TOWING_SERVICE:
		return nil
	case model.PRODUCT_JUMP_START_SERVICE:
		return nil
	case model.PRODUCT_KEY_UNLOCK_SERVICE:
		return nil
	}
	return (fmt.Errorf("Unsupported Vehicle Category \"%s\"", category))
}
