package url

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// LinksHandler handles URL shortening and retrieval operations.
type LinksHandler struct {
	log     *slog.Logger
	service *services.LinkService
}

// NewLinkHandler creates a new LinksHandler instance.
func NewLinkHandler(log *slog.Logger, service *services.LinkService) *LinksHandler {
	return &LinksHandler{log: log, service: service}
}

// GetRequest represents a request to retrieve the original URL.
type GetRequest struct {
}

// GetResponse represents the response for retrieving the original URL.
type GetResponse struct {
	resp.Response
	OriginalURL string `json:"link,omitempty"`
}

// GetLink retrieves the original URL for a given short URL.
//	@Summary		Retrieve the original URL
//	@Description	Retrieves the original URL associated with the provided short URL.
//	@Tags			url
//	@Accept			json
//	@Produce		json
//	@Param			link	path		string	true	"Short URL"
//	@Success		200		{object}	GetResponse
//	@Failure		400		{object}	resp.Response
//	@Failure		404		{object}	resp.Response
//	@Failure		500		{object}	resp.Response
//	@Router			/{link} [get]
func (h *LinksHandler) GetLink(c *gin.Context) {
	const op = "handlers.url.GetLink"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", c.GetString("request_id")),
	)

	shortUrl := c.Param("link")

	originalURL, err := h.service.GetOriginalURL(shortUrl)
	if errors.Is(err, services.ErrNotFound) {
		log.Info("url was not found", slog.String("shortURL", shortUrl))
		c.JSON(http.StatusNotFound, resp.NotFound("url was not found"))
		return
	}
	if errors.Is(err, services.ErrInvalidLink) {
		log.Info("passed incorrect link", slog.String("shortURL", shortUrl))
		c.JSON(http.StatusBadRequest, resp.Response{
			Status: resp.StatusBadRequest,
			Error:  "Passed invalid short link",
		},
		)
		return
	}
	if err != nil {
		log.Error("failed to find url", sl.Err(err))
		c.JSON(http.StatusInternalServerError, resp.InternalError("failed to find url"))
		return
	}

	log.Info("original url retrieved", slog.String("originalURL", originalURL))

	c.JSON(http.StatusOK, GetResponse{
		Response:    resp.OK(),
		OriginalURL: originalURL,
	})
}

// SaveRequest represents a request to save a new short URL.
type SaveRequest struct {
	OriginalURL string `json:"url" validate:"required,url"`
}

// SaveResponse represents the response for saving a new short URL.
type SaveResponse struct {
	resp.Response
	ShortURL string `json:"link,omitempty"`
}

// SaveLink saves a new short URL for the provided original URL.
//	@Summary		Save a new short URL
//	@Description	Saves a new short URL for the provided original URL.
//	@Tags			url
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SaveRequest	true	"Original URL to shorten"
//	@Success		200		{object}	SaveResponse
//	@Failure		400		{object}	resp.Response
//	@Failure		500		{object}	resp.Response
//	@Router			/ [post]
func (h *LinksHandler) SaveLink(c *gin.Context) {
	const op = "handlers.url.SaveLink"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", c.GetString("request_id")),
	)

	var req SaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			c.JSON(http.StatusBadRequest, resp.Error("empty request"))
			return
		}

		log.Error("failed to decode request body", sl.Err(err))
		c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))
		return
	}

	log.Info("request body decoded", slog.Any("request", req))
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	if err := validate.Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("invalid request", sl.Err(err))
		c.JSON(http.StatusBadRequest, resp.ValidationError(validateErr))
		return
	}

	shortURL, err := h.service.Save(req.OriginalURL, 5)
	if errors.Is(err, services.ErrInvalidURL) {
		log.Info("passed incorrect link", slog.String("originalURL", req.OriginalURL))
		c.JSON(http.StatusBadRequest, resp.Response{
			Status: resp.StatusBadRequest,
			Error:  "Passed invalid url to shorten",
		},
		)
		return
	}
	if errors.Is(err, services.ErrMaxRetriesExceeded) {
		log.Info("max retries exceeded, could not save", slog.String("originalURL", req.OriginalURL))
		c.JSON(http.StatusInternalServerError, resp.InternalError("max retries exceeded"))
		return
	}
	if err != nil {
		log.Error("failed to add url", sl.Err(err))
		c.JSON(http.StatusInternalServerError, resp.InternalError("failed to add url"))
		return
	}

	log.Info("url added", slog.String("shortURL", shortURL))

	c.JSON(http.StatusOK, SaveResponse{
		Response: resp.OK(),
		ShortURL: shortURL,
	})
}
