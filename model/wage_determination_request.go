package model

type WageDeterminationRequest struct {
	WageDetermination WageDetermination
	Location          string
	County            string
	State             string
	Url               string
	Error             error
}
