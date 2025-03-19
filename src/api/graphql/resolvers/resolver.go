package resolvers

import (
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// Root Resolver
type Resolver struct {
    userUsecase    *usecase.UserUsecase
    jwtService     *auth.JWTService
}

// NewResolver
func NewResolver(
    userUsecase *usecase.UserUsecase,
    jwtService *auth.JWTService,
) *Resolver {
    return &Resolver{
        userUsecase: userUsecase,
        jwtService:  jwtService,
    }
}
