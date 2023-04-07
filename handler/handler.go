package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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

	name := r.MultipartForm.Value["name"][0]
	address := r.MultipartForm.Value["address"][0]
	hasWifi, _ := strconv.ParseBool(r.MultipartForm.Value["hasWifi"][0])
	hasKitchen, _ := strconv.ParseBool(r.MultipartForm.Value["hasKitchen"][0])
	hasAirConditioning, _ := strconv.ParseBool(r.MultipartForm.Value["hasAirConditioning"][0])
	hasFreeParking, _ := strconv.ParseBool(r.MultipartForm.Value["hasFreeParking"][0])
	minimimGuests, _ := strconv.ParseUint(r.MultipartForm.Value["minimumGuests"][0], 10, 32)
	maximumGuests, _ := strconv.ParseUint(r.MultipartForm.Value["maximumGuests"][0], 10, 32)

	savedAccomodation := h.Service.SaveAccomodation(model.Accomodation{Name: name, Address: address, HasWifi: hasWifi, HasKitchen: hasKitchen, HasAirConditioning: hasAirConditioning, HasFreeParking: hasFreeParking, MinimimGuests: uint(minimimGuests), MaximumGuests: uint(maximumGuests)})

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
