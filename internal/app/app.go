package app

import (
	"context"
	"fmt"
	_ "pr-assignment/docs"
	"pr-assignment/internal/adapter/in/http/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//endpoints:
///team/add
///team/get
//
///users/setIsActive
///users/getReview
//
///pullRequest/create
///pullRequest/merge
///pullRequest/reassign

type Server struct {
	prHandler   *handler.PullRequestHandler
	userHandler *handler.UserHandler
	statHandler *handler.StatHandler
}

func NewServer(prHandler *handler.PullRequestHandler, userHandler *handler.UserHandler, statHandler *handler.StatHandler) *Server {
	return &Server{prHandler: prHandler, userHandler: userHandler, statHandler: statHandler}
}

func (s *Server) RunServer(ctx context.Context) error {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.GET("/team/get", s.userHandler.GetTeam)
	router.POST("/team/add", s.userHandler.AddTeam)

	router.POST("/users/setIsActive", s.userHandler.SetIsUserActive)
	router.GET("/users/getReview", s.userHandler.GetReviews)

	router.POST("/pullRequest/create", s.prHandler.CreatePullRequest)
	router.POST("/pullRequest/merge", s.prHandler.MergePullRequest)
	router.POST("/pullRequest/reassign", s.prHandler.ReassignPullRequest)

	router.GET("/stat/pull_request/reviewers", s.statHandler.GetReviewersCountedByPR)
	router.GET("/stat/users/reviews", s.statHandler.GetReviewsCountedByUser)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
