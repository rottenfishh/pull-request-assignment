package model

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}
