package handler

import (
	"fmt"
	"net/http"
	"pr-assignment/internal/adapter/in/http/dto"
	"pr-assignment/internal/model"
	"pr-assignment/internal/service"

	"github.com/gin-gonic/gin"
)

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
// @Success      200  {object}   dto.UserResponse
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

	user, err := h.userService.SetUserActive(ctx, query.UserID, query.IsActive)
	if err != nil {
		errResp := model.ParseErrorResponse(err)
		if errResp.Error.Code == model.NotFound {
			c.IndentedJSON(http.StatusNotFound, model.ParseErrorResponse(err))
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, errResp)
		return
	}

	if !query.IsActive {
		err = h.prService.ReassignReviewsAfterDeath(ctx, query.UserID)
		if err != nil {
			fmt.Println(err)
		}
	}

	response := dto.UserResponse{User: *user}
	c.IndentedJSON(http.StatusOK, response)
}

// GetReview godoc
// @Summary      get prs where user is reviewer
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id query string true "user id"
// @Success      200  {object}   dto.UserPrsResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /users/getReview [get]
func (h *UserHandler) GetReviews(c *gin.Context) {
	ctx := c.Request.Context()

	var userID dto.UserIDQuery
	if err := c.BindQuery(&userID); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prs, err := h.prService.GetPRsByUser(ctx, userID.UserID)
	if err != nil {
		errResp := model.ParseErrorResponse(err)
		if errResp.Error.Code == model.NotFound {
			c.IndentedJSON(http.StatusNotFound, model.ParseErrorResponse(err))
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, errResp)
		return
	}

	responses := []model.PullRequestShort{}
	for _, pr := range prs {
		prResponse := model.PullRequestShort{PullRequestID: pr.PullRequestID, PullRequestName: pr.PullRequestName,
			AuthorID: pr.AuthorID, Status: pr.Status}
		responses = append(responses, prResponse)
	}

	resp := dto.UserPrsResponse{UserID: userID.UserID, PullRequests: responses}
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
		errorResp := model.ParseErrorResponse(err)
		if errorResp.Error.Code == model.TeamExists {
			c.IndentedJSON(400, model.ParseErrorResponse(err))
			return
		}
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
// @Param        team_name query string true "team_name"
// @Success      200  {object}   model.Team
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

	fmt.Println("q", query.TeamName)
	team, err := h.userService.GetTeam(ctx, query.TeamName)
	if err != nil {
		errorResp := model.ParseErrorResponse(err)
		if errorResp.Error.Code == model.NotFound {
			c.IndentedJSON(http.StatusNotFound, model.ParseErrorResponse(err))
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, model.ParseErrorResponse(err))
		return
	}
	c.IndentedJSON(http.StatusOK, team)
}

// KillTeam godoc
// @Summary      deactivate all users in team
// @Description  set users status to not active by a given team name
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        team_name body dto.TeamName true "team_name"
// @Success      200  {object}   model.Team
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /team/kill [post]
func (h *UserHandler) KillTeam(c *gin.Context) {
	ctx := c.Request.Context()

	var query dto.TeamName
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.userService.KillTeam(ctx, query.TeamName)
	if err != nil {
		errResp := model.ParseErrorResponse(err)
		if errResp.Error.Code == model.NotFound {
			c.IndentedJSON(http.StatusNotFound, model.ParseErrorResponse(err))
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, model.ParseErrorResponse(err))
		return
	}

	c.IndentedJSON(http.StatusOK, team)
}
