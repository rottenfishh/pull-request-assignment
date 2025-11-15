package handler

import (
	"fmt"
	"net/http"
	"pr-assignment/internal/adapter/in/http/dto"
	"pr-assignment/internal/model"
	"pr-assignment/internal/service"

	"github.com/gin-gonic/gin"
)

///users/setIsActive
///users/getReview
//endpoints:
///team/add
///team/get
// do the queries as in spec

type UserHandler struct {
	userService *service.UserService
	prService   *service.PullRequestService
}

func NewUserHandler(userService *service.UserService, prService *service.PullRequestService) *UserHandler {
	return &UserHandler{userService: userService, prService: prService}
}

// SetIsUserActive godoc
// @Summary      set user is active status
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        query body dto.StatusQuery true "user id, status"
// @Success      200  {array}   dto.UserResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /users/setIsActive [post]
func (h *UserHandler) SetIsUserActive(c *gin.Context) {
	ctx := c.Request.Context()

	var query dto.StatusQuery
	if err := c.BindJSON(&query); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.SetUserActive(ctx, query.UserId, query.IsActive)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.UserResponse{*user}
	c.IndentedJSON(http.StatusOK, response)
}

// GetReview godoc
// @Summary      get prs where user is reviewer
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        query query dto.UserIdQuery true "user id"
// @Success      200  {array}   dto.UserPrsResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /users/getReview [get]
func (h *UserHandler) GetReviews(c *gin.Context) {
	ctx := c.Request.Context()

	var userId dto.UserIdQuery
	if err := c.BindQuery(&userId); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prs, err := h.prService.GetPRsByUser(ctx, userId.UserId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.ParseErrorResponse(err))
		return
	}

	responses := []model.PullRequestShort{}
	for _, pr := range prs {
		prResponse := model.PullRequestShort{pr.PullRequestId, pr.PullRequestName, pr.AuthorId, pr.Status}
		responses = append(responses, prResponse)
	}

	resp := dto.UserPrsResponse{userId.UserId, responses}
	c.IndentedJSON(http.StatusOK, resp)
}

// AddTeam godoc
// @Summary      add new team
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        team body model.Team true "team object: team_name {members}"
// @Success      201  {object} model.Team
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /team/add [post]
func (h *UserHandler) AddTeam(c *gin.Context) {
	ctx := c.Request.Context()

	var team model.Team
	if err := c.BindJSON(&team); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.AddTeam(ctx, team)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, model.ParseErrorResponse(err))
		return
	}
	c.IndentedJSON(http.StatusCreated, team)
}

// GetTeam godoc
// @Summary      get existing team
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        query query dto.TeamName true "team_name"
// @Success      200  {array}   model.Team
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /team/get [get]
func (h *UserHandler) GetTeam(c *gin.Context) {
	ctx := c.Request.Context()

	var query dto.TeamName
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.userService.GetTeam(ctx, query.TeamName)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, model.ParseErrorResponse(err))
		return
	}
	c.IndentedJSON(http.StatusOK, team)
}
