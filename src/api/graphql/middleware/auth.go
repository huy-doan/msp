package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
)

// GraphQLAuthMiddleware creates a middleware for GraphQL authentication
func GraphQLAuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// For GraphQL, we don't want to abort the request if authentication fails
		// Instead, we just set context values that resolvers can check
		c.Set("authenticated", false)

		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check if header has the correct format
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.Next()
			return
		}

		// Parse and validate the token
		tokenString := headerParts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Set authentication information in context
		c.Set("authenticated", true)
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// WithAuth creates a GraphQL resolver context with auth information
func WithAuth(ctx context.Context, c *gin.Context) context.Context {
	for _, key := range []string{"authenticated", "userId", "username", "role"} {
		if value, exists := c.Get(key); exists {
			ctx = context.WithValue(ctx, key, value)
		}
	}
	return ctx
}

// CheckAuth checks if user is authenticated in GraphQL resolver context
func CheckAuth(ctx context.Context) error {
	authenticated, ok := ctx.Value("authenticated").(bool)
	if !ok || !authenticated {
		return errors.New("not authenticated")
	}
	return nil
}
