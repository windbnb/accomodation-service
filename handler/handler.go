package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/windbnb/accomodation-service/client"
	"github.com/windbnb/accomodation-service/model"
	"github.com/windbnb/accomodation-service/service"
	"github.com/windbnb/accomodation-service/tracer"
	"github.com/windbnb/accomodation-service/util"
)

type Handler struct {
	Service *service.AccomodationService
	Tracer  opentracing.Tracer
	Closer  io.Closer
}

func (handler *Handler) Healthcheck(w http.ResponseWriter, _ *http.Request) {
    _, _ = fmt.Fprintln(w, "Healthy!")
}

func (handler *Handler) Ready(w http.ResponseWriter, _ *http.Request) {
    _, _ = fmt.Fprintln(w, "Ready!")
}

func (h *Handler) CreateAccomodation(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("createAccomodationHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling create accomodation at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, _ := strconv.ParseUint(r.MultipartForm.Value["userId"][0], 10, 32)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	userResponse, err := client.GetUserById(uint(userId))
	if err != nil {
		tracer.LogError(span, err)
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
	savedAccomodation := h.Service.SaveAccomodation(newAccomodation, ctx)

	files := r.MultipartForm.File["images"]

	fileNames, err := util.SaveHeaderFileImages(files)
	if err != nil {
		tracer.LogError(span, err)
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
	span := tracer.StartSpanFromRequest("acceptReservationTypeHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling accept reservation type at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	accomodationId, _ := strconv.Atoi(params["id"])

	ctx := tracer.ContextWithSpan(context.Background(), span)

	var acceptReservationType *model.AcceptReservationTypeDTO
	err := json.NewDecoder(r.Body).Decode(&acceptReservationType)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	tokenString := r.Header.Get("Authorization")
	userResponse, err := client.AuthorizeHost(tokenString)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusUnauthorized})
		return
	}

	if userResponse.Role != "HOST" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusUnauthorized})
		return
	}

	accommodation, err := h.Service.UpdateAccommodationAcceptReservationType(uint(accomodationId), acceptReservationType.AcceptReservationType, userResponse.Id, ctx)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	json.NewEncoder(w).Encode(accommodation)
}

func (h *Handler) FindAccommodationById(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("findAccomodationByIdHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling find accomodation by id at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	accomodationId, _ := strconv.Atoi(params["id"])

	ctx := tracer.ContextWithSpan(context.Background(), span)

	accomodation, err := h.Service.FindAccomodationById(uint(accomodationId), ctx)
	availalbleTerms, err2 := h.Service.FindAvailableTerms(uint(accomodationId), ctx)
	if err != nil || err2 != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusNotFound})
		return
	}

	var returnedValue = model.AccommodationBasicDTO{
		Id:                    accomodation.ID,
		MaximumGuests:         accomodation.MaximumGuests,
		MinimimGuests:         accomodation.MinimimGuests,
		AvailableTerms:        availalbleTerms,
		UserID:                accomodation.UserId,
		AcceptReservationType: accomodation.AcceptReservationType,
		Name:                  accomodation.Name}

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

	// http.ServeFile(w, r, "/app/images/"+filename)
	bytes, err := ioutil.ReadFile("/app/images/" + filename)
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string

	mimeType := http.DetectContentType(bytes)
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	base64Encoding += base64.StdEncoding.EncodeToString(bytes)
	json.NewEncoder(w).Encode(base64Encoding)

}

func (h *Handler) CreatePrice(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("createPriceHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling create price at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	var createPricesDTO []model.CreatePriceDTO
	json.NewDecoder(r.Body).Decode(&createPricesDTO)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	tokenString := r.Header.Get("Authorization")
	userResponse, err := client.AuthorizeHost(tokenString)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusUnauthorized})
		return
	}

	if userResponse.Role != "HOST" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusUnauthorized})
		return
	}

	var pricesDTO []model.PriceDTO

	for _, createPriceDTO := range createPricesDTO {
		newPrice := util.FromCreatePriceDTOToPrice(createPriceDTO)
		newPrice.Active = true
		_, err := h.Service.FindAccomodationById(newPrice.AccomodationID, ctx)
		if err != nil {
			tracer.LogError(span, err)
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
	span := tracer.StartSpanFromRequest("updatePriceHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling update price at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	priceId, _ := strconv.ParseUint(params["id"], 10, 32)

	var updatePriceDTO model.UpdatePriceDTO
	json.NewDecoder(r.Body).Decode(&updatePriceDTO)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	newPrice := util.FromUpdatePriceDTOToPrice(updatePriceDTO)

	savedPrice, err := h.Service.UpdatePrice(newPrice, priceId, ctx)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}
	priceDTO := savedPrice.ToDTO()

	json.NewEncoder(w).Encode(priceDTO)

}

func (h *Handler) DeletePrice(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("deletePriceHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling delete price at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	priceId, _ := strconv.ParseUint(params["id"], 10, 32)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	err := h.Service.DeletePrice(priceId, ctx)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusNotFound})
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *Handler) CreateAvailableTerm(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("createAvailableTermHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling create available term at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	var createAvailableTermsDTO []model.CreateAvailableTermDTO
	json.NewDecoder(r.Body).Decode(&createAvailableTermsDTO)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	tokenString := r.Header.Get("Authorization")
	userResponse, err := client.AuthorizeHost(tokenString)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusUnauthorized})
		return
	}

	if userResponse.Role != "HOST" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusUnauthorized})
		return
	}

	var availableTermsDTO []model.AvailableTermDTO

	for _, createAvailableTermDTO := range createAvailableTermsDTO {
		newAvailableTerm := util.FromCreateAvailableTermDTOToAvailableTerm(createAvailableTermDTO)
		_, err := h.Service.FindAccomodationById(newAvailableTerm.AccomodationID, ctx)
		if err != nil {
			tracer.LogError(span, err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
			return
		}
		savedAvailableTerm := h.Service.SaveAvailableTerm(newAvailableTerm, ctx)
		availableTermDTO := savedAvailableTerm.ToDTO()
		availableTermsDTO = append(availableTermsDTO, availableTermDTO)
	}

	json.NewEncoder(w).Encode(availableTermsDTO)

}

func (h *Handler) UpdateAvailableTerm(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("updateAvailableTermHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling update available term at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	availableTermId, _ := strconv.ParseUint(params["id"], 10, 32)

	var updateAvailableTermDTO model.UpdateAvailableTermDTO
	json.NewDecoder(r.Body).Decode(&updateAvailableTermDTO)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	tokenString := r.Header.Get("Authorization")
	userResponse, err := client.AuthorizeHost(tokenString)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusUnauthorized})
		return
	}

	if userResponse.Role != "HOST" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusUnauthorized})
		return
	}

	newAvailableTerm := util.FromUpdateAvailableTermDTOToAvailableTerm(updateAvailableTermDTO)

	savedAvailableTerm, err := h.Service.UpdateAvailableTerm(newAvailableTerm, availableTermId, ctx)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}
	availableTermDTO := savedAvailableTerm.ToDTO()

	json.NewEncoder(w).Encode(availableTermDTO)

}

func (h *Handler) DeleteAvailableTerm(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("deleteAvailableTermHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling delete available term at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	availableTermId, _ := strconv.ParseUint(params["id"], 10, 32)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	tokenString := r.Header.Get("Authorization")
	userResponse, err := client.AuthorizeHost(tokenString)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusUnauthorized})
		return
	}

	if userResponse.Role != "HOST" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusUnauthorized})
		return
	}

	er := h.Service.DeleteAvailableTerm(availableTermId, ctx)
	if er != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusNotFound})
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *Handler) CreateReservedTerm(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("createReservedTermHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling create reserved term at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	var createReservedTermDTO model.CreateReservedTermDTO
	json.NewDecoder(r.Body).Decode(&createReservedTermDTO)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	newReservedTerm := util.FromCreateReservedTermDTOToReservedTerm(createReservedTermDTO)
	_, err := h.Service.FindAccomodationById(newReservedTerm.AccomodationID, ctx)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}
	savedReservedTerm := h.Service.SaveReservedTerm(newReservedTerm, ctx)
	reservedTermDTO := savedReservedTerm.ToDTO()

	json.NewEncoder(w).Encode(reservedTermDTO)

}

func (h *Handler) DeleteReservedTerm(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("deleteReservedTermHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling delete reserved term at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	reservedTermId, _ := strconv.ParseUint(params["id"], 10, 32)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	err := h.Service.DeleteReservedTerm(reservedTermId, ctx)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusNotFound})
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *Handler) DeleteHostAccomodation(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("deleteHostAccomodationHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling delete host accomodation at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	hostId, err := strconv.ParseUint(params["hostId"], 10, 32)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "cannot parse host id", StatusCode: http.StatusBadRequest})
		return
	}

	err = h.Service.DeleteHostAccomodation(uint(hostId), ctx)

	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusNotFound})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SearchAccomodation(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("searchAccomodationHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling search accomodation at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	var searchAccomodationDTO model.SearchAccomodationDTO
	json.NewDecoder(r.Body).Decode(&searchAccomodationDTO)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	accomodationsDTO := h.Service.SearchAccomodations(searchAccomodationDTO, ctx)
	json.NewEncoder(w).Encode(accomodationsDTO)

}

func (h *Handler) FindAccommodationsForHost(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("findAccomodationsForHostHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling find accomodations for host at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	hostId, err := strconv.ParseUint(params["hostId"], 10, 32)

	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "cannot parse host id", StatusCode: http.StatusBadRequest})
		return
	}

	ctx := tracer.ContextWithSpan(context.Background(), span)

	accomodationsDTO := h.Service.FindAccommodationsForHost(uint(hostId), ctx)
	json.NewEncoder(w).Encode(accomodationsDTO)

}

func (h *Handler) GetAvailableTermsForAccomodation(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("getAvailableTermsForAccomodationHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get available terms for accomodation at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	accomodationId, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "cannot parse accomodation id", StatusCode: http.StatusBadRequest})
		return
	}

	ctx := tracer.ContextWithSpan(context.Background(), span)

	availableTermsDTO := h.Service.GetAvailableTermsForAccomodation(uint(accomodationId), ctx)
	json.NewEncoder(w).Encode(availableTermsDTO)

}

func (h *Handler) GetPricesForAccomodation(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("getPricesForAccomodationHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get prices for accomodation at %s\n", r.URL.Path)),
	)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	accomodationId, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "cannot parse accomodation id", StatusCode: http.StatusBadRequest})
		return
	}

	ctx := tracer.ContextWithSpan(context.Background(), span)

	pricesDTO := h.Service.GetPricesForAccomodation(uint(accomodationId), ctx)
	json.NewEncoder(w).Encode(pricesDTO)

}
