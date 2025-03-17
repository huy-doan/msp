package resolvers

import (
	"github.com/vnlab/makeshop-payment/src/application/services"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
)

// Resolver là resolver gốc
type Resolver struct {
    userService    *services.UserService
    jwtService     *auth.JWTService
}

// NewResolver tạo resolver mới
func NewResolver(
    userService *services.UserService,
    jwtService *auth.JWTService,
) *Resolver {
    return &Resolver{
        userService:    userService,
        jwtService:     jwtService,
    }
}
