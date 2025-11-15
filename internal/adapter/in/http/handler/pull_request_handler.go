package handler

import (
	"fmt"
	"net/http"
	"pr-assignment/internal/adapter/in/http/dto"
	"pr-assignment/internal/model"
	"pr-assignment/internal/service"

	"github.com/gin-gonic/gin"
)

///pullRequest/create
///pullRequest/merge
///pullRequest/reassign

type PullRequestHandler struct {
	prService *service.PullRequestService
}

func NewPullRequestHandler(prService *service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{prService: prService}
}

// CreatePullRequest godoc
// @Summary      Create new Pull Request
// @Description  create new pr and assign reviewers automatically
// @Tags         pull requests
// @Accept       json
// @Produce      json
// @Param        query body dto.PullRequestQuery true "PR DATA"
// @Success      201  {object}   dto.PrResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /pullRequest/create [post]
func (h *PullRequestHandler) CreatePullRequest(c *gin.Context) {
	ctx := c.Request.Context()

	var query dto.PullRequestQuery
	if err := c.BindJSON(&query); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pr, err := h.prService.CreatePR(ctx, query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.ParseErrorResponse(err))
		fmt.Println(err)
		return
	}

	newPr := dto.PrResponse{
		PullRequestShort:  pr.PullRequestShort,
		AssignedReviewers: pr.AssignedReviewers,
	}

	c.IndentedJSON(http.StatusCreated, newPr)
}

// MergePullRequest godoc
// @Summary      Merge Pull Request
// @Tags         pull requests
// @Accept       json
// @Produce      json
// @Param        query body dto.PullRequestIdQuery true "PR ID"
// @Success      200  {array}   dto.PrMergedResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /pullRequest/merge [post]
func (h *PullRequestHandler) MergePullRequest(c *gin.Context) {
	ctx := c.Request.Context()
	var prIdQuery dto.PullRequestIdQuery
	if err := c.BindJSON(&prIdQuery); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.ParseErrorResponse(err))
		return
	}

	pr, err := h.prService.MergePR(ctx, prIdQuery.PrId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.ParseErrorResponse(err))
		return
	}

	prResponse := dto.PrResponse{PullRequestShort: pr.PullRequestShort, AssignedReviewers: pr.AssignedReviewers}
	updatedPr := dto.PrMergedResponse{PrMerged: dto.PrMerged{prResponse, pr.MergedAt}}

	c.IndentedJSON(http.StatusOK, updatedPr)
}

// ReassignPullRequest godoc
// @Summary      Reassign reviewer Pull Request
// @Tags         pull requests
// @Accept       json
// @Produce      json
// @Param        query body dto.PrReassignQuery true "Pr id, old reviewer id"
// @Success      200  {array}   dto.PrReassignResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /pullRequest/reassign [post]
func (h *PullRequestHandler) ReassignPullRequest(c *gin.Context) {
	ctx := c.Request.Context()
	var query dto.PrReassignQuery
	if err := c.BindJSON(&query); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.ParseErrorResponse(err))
		return
	}

	result, err := h.prService.ChangeReviewer(ctx, query.PullRequestId, query.OldReviewerId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.ParseErrorResponse(err))
		return
	}

	prResponse := dto.PrResponse{PullRequestShort: result.PullRequest.PullRequestShort,
		AssignedReviewers: result.PullRequest.AssignedReviewers}

	prReassignResponse := dto.PrReassignResponse{PrResponse: prResponse, ReplacedBy: result.NewReviewerId}

	c.IndentedJSON(http.StatusOK, prReassignResponse)
}
