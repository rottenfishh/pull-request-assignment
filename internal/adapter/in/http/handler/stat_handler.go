package handler

import (
	"net/http"
	"pr-assignment/internal/model"
	"pr-assignment/internal/service"

	"github.com/gin-gonic/gin"
)

type StatHandler struct {
	statService *service.StatService
}

func NewStatHandler(statService *service.StatService) *StatHandler {
	return &StatHandler{statService: statService}
}

// GetReviewersCountedByPR godoc
// @Summary      get prs with assigned reviewers
// @Description  get pull requests with reviewers and their number
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.PrReviewersCount
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /stat/pull_request/reviewers [get]
func (h *StatHandler) GetReviewersCountedByPR(c *gin.Context) {
	ctx := c.Request.Context()

	prReviewers, err := h.statService.GetReviewsCountedByPR(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.ParseErrorResponse(err))
	}

	c.IndentedJSON(http.StatusOK, prReviewers)
}

// GetReviewsCountedByUser godoc
// @Summary      get users and number of reviews
// @Description  get users and in how many pull requests they are reviewers
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.UserReviewsCount
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /stat/users/reviews [get]
func (h *StatHandler) GetReviewsCountedByUser(c *gin.Context) {
	ctx := c.Request.Context()

	userPrs, err := h.statService.GetReviewsCountedByUser(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.ParseErrorResponse(err))
	}

	c.IndentedJSON(http.StatusOK, userPrs)
}
