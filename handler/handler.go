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

func (h *Handler) UpdateAccommodationAcceptReservationType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	accomodationId, _ := strconv.Atoi(params["id"])

	var acceptReservationType *model.AcceptReservationTypeDTO
	err := json.NewDecoder(r.Body).Decode(&acceptReservationType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	tokenString := r.Header.Get("Authorization")
	userResponse, err := client.AuthorizeHost(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusUnauthorized})
		return
	}
	if userResponse.Role != "HOST" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusUnauthorized})
		return
	}

	accommodation, err := h.Service.UpdateAccommodationAcceptReservationType(uint(accomodationId), acceptReservationType.AcceptReservationType, userResponse.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	json.NewEncoder(w).Encode(accommodation)
}

func (h *Handler) FindAccommodationById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	accomodationId, _ := strconv.Atoi(params["id"])

	accomodation, err := h.Service.FindAccomodationById(uint(accomodationId))
	availalbleTerms, err2 := h.Service.FindAvailableTerms(uint(accomodationId))
	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusNotFound})
		return
	}

	var returnedValue = model.AccommodationBasicDTO{
		Id:                    accomodation.ID,
		MaximumGuests:         accomodation.MaximumGuests,
		MinimimGuests:         accomodation.MinimimGuests,
		AvailableTerms:        availalbleTerms,
		AcceptReservationType: accomodation.AcceptReservationType}

	json.NewEncoder(w).Encode(returnedValue)
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

	http.ServeFile(w, r, "./images/"+filename)
}

func (h *Handler) CreatePrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createPricesDTO []model.CreatePriceDTO
	json.NewDecoder(r.Body).Decode(&createPricesDTO)

	var pricesDTO []model.PriceDTO

	for _, createPriceDTO := range createPricesDTO {
		newPrice := util.FromCreatePriceDTOToPrice(createPriceDTO)
		_, err := h.Service.FindAccomodationById(newPrice.AccomodationID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
			return
		}
		savedPrice := h.Service.SavePrice(newPrice)
		priceDTO := savedPrice.ToDTO()
		pricesDTO = append(pricesDTO, priceDTO)
	}

	json.NewEncoder(w).Encode(pricesDTO)

}

func (h *Handler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	priceId, _ := strconv.ParseUint(params["id"], 10, 32)

	var updatePriceDTO model.UpdatePriceDTO
	json.NewDecoder(r.Body).Decode(&updatePriceDTO)

	newPrice := util.FromUpdatePriceDTOToPrice(updatePriceDTO)

	savedPrice, err := h.Service.UpdatePrice(newPrice, priceId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}
	priceDTO := savedPrice.ToDTO()

	json.NewEncoder(w).Encode(priceDTO)

}

func (h *Handler) DeletePrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	priceId, _ := strconv.ParseUint(params["id"], 10, 32)

	err := h.Service.DeletePrice(priceId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusNotFound})
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *Handler) CreateAvailableTerm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createAvailableTermsDTO []model.CreateAvailableTermDTO
	json.NewDecoder(r.Body).Decode(&createAvailableTermsDTO)

	var availableTermsDTO []model.AvailableTermDTO

	for _, createAvailableTermDTO := range createAvailableTermsDTO {
		newAvailableTerm := util.FromCreateAvailableTermDTOToAvailableTerm(createAvailableTermDTO)
		_, err := h.Service.FindAccomodationById(newAvailableTerm.AccomodationID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
			return
		}
		savedAvailableTerm := h.Service.SaveAvailableTerm(newAvailableTerm)
		availableTermDTO := savedAvailableTerm.ToDTO()
		availableTermsDTO = append(availableTermsDTO, availableTermDTO)
	}

	json.NewEncoder(w).Encode(availableTermsDTO)

}

func (h *Handler) UpdateAvailableTerm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	availableTermId, _ := strconv.ParseUint(params["id"], 10, 32)

	var updateAvailableTermDTO model.UpdateAvailableTermDTO
	json.NewDecoder(r.Body).Decode(&updateAvailableTermDTO)

	newAvailableTerm := util.FromUpdateAvailableTermDTOToAvailableTerm(updateAvailableTermDTO)

	savedAvailableTerm, err := h.Service.UpdateAvailableTerm(newAvailableTerm, availableTermId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}
	availableTermDTO := savedAvailableTerm.ToDTO()

	json.NewEncoder(w).Encode(availableTermDTO)

}

func (h *Handler) DeleteAvailableTerm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	availableTermId, _ := strconv.ParseUint(params["id"], 10, 32)

	err := h.Service.DeleteAvailableTerm(availableTermId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusNotFound})
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *Handler) CreateReservedTerm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createReservedTermDTO model.CreateReservedTermDTO
	json.NewDecoder(r.Body).Decode(&createReservedTermDTO)

	newReservedTerm := util.FromCreateReservedTermDTOToReservedTerm(createReservedTermDTO)
	_, err := h.Service.FindAccomodationById(newReservedTerm.AccomodationID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}
	savedReservedTerm := h.Service.SaveReservedTerm(newReservedTerm)
	reservedTermDTO := savedReservedTerm.ToDTO()

	json.NewEncoder(w).Encode(reservedTermDTO)

}

func (h *Handler) DeleteReservedTerm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	reservedTermId, _ := strconv.ParseUint(params["id"], 10, 32)

	err := h.Service.DeleteReservedTerm(reservedTermId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusNotFound})
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *Handler) DeleteHostAccomodation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	hostId, err := strconv.ParseUint(params["hostId"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "cannot parse host id", StatusCode: http.StatusBadRequest})
		return
	}

	err = h.Service.DeleteHostAccomodation(uint(hostId))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusNotFound})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
