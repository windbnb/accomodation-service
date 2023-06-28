package util

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func SaveHeaderFileImages(files []*multipart.FileHeader) ([]string, error) {
	if len(files) == 0 {
		return nil, errors.New("you have to provide at least one image")
	}

	fileNames := []string{}
	for _, fileHeader := range files {

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			return nil, err
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			return nil, errors.New("the provided file format is not allowed, please upload a JPEG or PNG image")
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}

		filenameTokens := strings.Split(filepath.Ext(fileHeader.Filename), ".")
		fileExtension := filenameTokens[len(filenameTokens)-1]
		saveFileName := uuid.New().String() + "." + fileExtension
		fileNames = append(fileNames, saveFileName)
		f, err := os.Create(fmt.Sprintf("/app/images/%s", saveFileName))

		if err != nil {
			return nil, err
		}

		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			return nil, err
		}
	}

	return fileNames, nil
}

