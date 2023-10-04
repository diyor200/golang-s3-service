package v1

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang-diplom-work/internal/entity"
	"golang-diplom-work/internal/service"
	"io"
	"net/http"
	"os"
)

const maxFileSize = 8 << 20

var imageTypes = map[string]interface{}{
	"image/jpeg": nil,
	"image/png":  nil,
}

type fileRoutes struct {
	fileService service.Files
}

func newFileRoutes(g *echo.Group, files service.Files) {
	var r = &fileRoutes{
		fileService: files,
	}
	g.POST("/upload/file", r.upload)
	//g.GET("/get/files", r.getFiles)
}

func (r *fileRoutes) upload(c echo.Context) error {
	file, fileHeader, err := c.Request().FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	defer file.Close()
	if fileHeader.Size > maxFileSize {
		newErrorResponse(c, http.StatusBadRequest, "file size should be no more than 8MB")
		return ErrBigFileSize
	}
	buffer := make([]byte, fileHeader.Size)
	if _, err := file.Read(buffer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	contentType := http.DetectContentType(buffer)

	//	Validate file type
	if _, ex := imageTypes[contentType]; !ex {
		newErrorResponse(c, http.StatusBadRequest, "file type not supported")
		return err
	}

	f, err := os.OpenFile(fileHeader.Filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to create temp file")
	}
	defer f.Close()

	if _, err := io.Copy(f, bytes.NewReader(buffer)); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to write chunk to temp file")

	}

	url, err := r.fileService.Upload(c.Request().Context(), entity.File{
		Name: f.Name(),
		Size: fileHeader.Size,
		URL:  fileHeader.Filename})
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "cannot send file")
		return err
	}

	if err := os.Remove(fileHeader.Filename); err != nil {
		log.Errorf("failed to delete corrupted temp file: %s", err.Error())
	}

	return c.JSON(http.StatusOK, response{url})

}

type response struct {
	Message string
}
