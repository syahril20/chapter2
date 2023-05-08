package middleware

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("img") // MENERIMA FILE

		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if file == nil {
			fmt.Println("Kosong")
		}

		src, err := file.Open() // MEMBUKA FILE

		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		defer src.Close() // MENUTUP FILE

		tempFile, err := ioutil.TempFile("upload", "image-*.png") // MEMBBUAT FILE BARU OTOMATIS ADA DIRECTORY UPLOAD
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		defer tempFile.Close() // MENUTUP FUNC SETELAH SELESAI DIJALANKAN

		if _, err = io.Copy(tempFile, src); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		data := tempFile.Name() // upload/image-3e10e160.png
		fileName := data[7:]    // image-3e10e160.png

		c.Set("dataFile", fileName)

		return next(c)
	}
}
