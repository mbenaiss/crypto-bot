package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/mbenaiss/crypto-bot/internal/provider"
	"github.com/mbenaiss/crypto-bot/internal/service"
	"github.com/mbenaiss/crypto-bot/models"
)

func uploadCsv(svc *service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Join(os.TempDir(), filepath.Base(file.Filename))
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		err = svc.ReadFromFile(filename, provider.ToProviderName(c.Param("provider")))
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("unable to read file: %s", err.Error()))
			return
		}
		c.Status(http.StatusCreated)
	}
}

func addProvider(svc *service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		p := models.Provider{}
		err := c.BindJSON(&p)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("unable to parse provider: %s", err.Error()))
			return
		}
		err = svc.AddProvider(p)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("unable to add provider: %s", err.Error()))
			return
		}
		c.Status(http.StatusCreated)
	}
}
