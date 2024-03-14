package handler

import (
	"gallery-api/db"
	"gallery-api/model"
	"gallery-api/preview"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func UploadPhoto(c *gin.Context) {
	p := model.PhotoEntity{}
	err := c.ShouldBind(&p)
	if err != nil {
		c.AbortWithStatusJSON(400, "Bad Request")
	}

	uuidString := uuid.New().String()
	uniqueName := uuidString + p.File.Filename
	filePath := filepath.Join("resources/storage/images", uniqueName)
	err = c.SaveUploadedFile(p.File, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Uploading error"})
		return
	}
	thumbName, err := preview.GeneratePreview(filePath, uuidString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Preview generation error"})
	}

	photo := &db.Photo{
		Model:   gorm.Model{},
		Name:    uniqueName,
		Preview: thumbName,
	}
	conn := db.GetDBConnection()
	conn.Create(&photo)

	photo.ImageURL += "http://localhost:8080/take/image/" + strconv.Itoa(int(photo.ID))
	photo.PreviewURL += "http://localhost:8080/take/preview/" + strconv.Itoa(int(photo.ID))
	conn.Save(&photo)

	c.JSON(http.StatusOK, gin.H{"message": "Photo uploaded successfully"})
}

func ListPhotos(c *gin.Context) {
	conn := db.GetDBConnection()
	rows, err := conn.Raw("SELECT id," +
		"name, " +
		"created_at, " +
		"updated_at, " +
		"preview, " +
		"image_url " +
		"FROM photos WHERE deleted_at IS NULL").Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var photos []db.Photo
	var photo db.Photo
	for rows.Next() {
		if err := rows.Scan(
			&photo.ID,
			&photo.Name,
			&photo.CreatedAt,
			&photo.UpdatedAt,
			&photo.Preview,
			&photo.ImageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		photos = append(photos, photo)
	}

	c.JSON(http.StatusOK, gin.H{"photos": photos})
}

func DeletePhoto(c *gin.Context) {
	id := c.Param("id")
	conn := db.GetDBConnection()

	var photo db.Photo
	if err := conn.First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	imagePath := filepath.Join("resources/storage/images", photo.Name)
	err := os.Remove(imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo from storage"})
	}

	thumbPath := filepath.Join("resources/storage/preview", photo.Preview)
	err = os.Remove(thumbPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete preview from storage"})
	}

	err = conn.Delete(&photo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted from db successfully"})
}

func GetImage(c *gin.Context) {
	id := c.Param("id")
	conn := db.GetDBConnection()

	var photo db.Photo
	if err := conn.First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	imagePath := filepath.Join("resources/storage/images", photo.Name)
	c.File(imagePath)
}

func GetPreview(c *gin.Context) {
	id := c.Param("id")
	conn := db.GetDBConnection()

	var photo db.Photo
	if err := conn.First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	previewPath := filepath.Join("resources/storage/preview", photo.Preview)
	c.File(previewPath)
}
