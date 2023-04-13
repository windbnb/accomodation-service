package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/windbnb/accomodation-service/client"
	"github.com/windbnb/accomodation-service/model"
	"github.com/windbnb/accomodation-service/service"
	"github.com/windbnb/accomodation-service/util"
)

type Handler struct {
	Service *service.AccomodationService
}

func (h *Handler) CreateAccomodation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, _ := strconv.ParseUint(r.MultipartForm.Value["userId"][0], 10, 32)

	userResponse, err := client.GetUserById(uint(userId))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadGateway})
		return
	}

	if userResponse.Role != "HOST" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusBadRequest})
		return
	}

	newAccomodation := util.ParseMultipartAccomodation(r)
	newAccomodation.UserId = uint(userId)
	savedAccomodation := h.Service.SaveAccomodation(newAccomodation)

	files := r.MultipartForm.File["images"]

	fileNames, err := util.SaveHeaderFileImages(files)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	accomodationDTO := savedAccomodation.ToDTO()
	for _, imageName := range fileNames {
		h.Service.SaveAccomodationImage(model.AccomodationImage{ImageName: imageName, AccomodationID: savedAccomodation.ID})
		accomodationDTO.Images = append(accomodationDTO.Images, imageName)
	}
	
	json.NewEncoder(w).Encode(accomodationDTO)
}

func (h *Handler) ImageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	filename := params["filename"]

	filenameTokens := strings.Split(filename, ".")
	fileExtension := filenameTokens[len(filenameTokens)-1]

	if fileExtension != "jpg" && fileExtension != "jpeg" && fileExtension != "png" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "unsupported file type", StatusCode: http.StatusBadRequest})
		return
	}

	http.ServeFile(w, r, "./images/" + filename)
}
