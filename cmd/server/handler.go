package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/mbenaiss/crypto-bot/internal/service"
)

func uploadCsv(svc *service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println(c.Param("provider"))
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

		err = svc.ReadFromFile(filename)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("unable to read file: %s", err.Error()))
			return
		}
		c.Status(http.StatusCreated)
	}
}
