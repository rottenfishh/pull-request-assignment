package model

type UserReviewsCount struct {
	User         User `json:"user_id"`
	ReviewsCount int  `json:"reviews_count"`
}
