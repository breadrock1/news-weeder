package http

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"news-weeder/internal/weeder"
)

func (s *HttpServer) CreateSearchGroup() error {
	group := s.server.Group("/weeder")
	group.POST("/search", s.Search)
	group.PUT("/store", s.StoreDocument)
	return nil
}

// Search
// @Summary Search similar news
// @Description Search similar news articles by semantic
// @ID search
// @Tags search
// @Accept  json
// @Produce json
// @Param jsonQuery body weeder.SearchParams true "Embeddings to search similar news"
// @Success 200 {object} []weeder.Document "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /weeder/search [post]
func (s *HttpServer) Search(c echo.Context) error {
	jsonForm := &weeder.SearchParams{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return err
	}

	docs, err := s.weeder.Weeder.Search(jsonForm)
	if err != nil {
		return err
	}

	return c.JSON(200, docs)
}

// StoreDocument
// @Summary Store document to storage
// @Description Store document to storage for similar searching
// @ID store
// @Tags store
// @Accept  json
// @Produce json
// @Param jsonQuery body weeder.Document true "Document data to store"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /weeder/store [put]
func (s *HttpServer) StoreDocument(c echo.Context) error {
	jsonForm := &weeder.Document{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return err
	}

	if err := s.weeder.Weeder.Append(jsonForm); err != nil {
		return err
	}

	resp := createStatusResponse(200, "Done")
	return c.JSON(200, resp)
}
