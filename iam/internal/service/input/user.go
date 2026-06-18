package input

import "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"

type RegisterInput struct {
	Info model.UserRegistrationInfo
}

type LoginInput struct {
	Login    string
	Password string
}
