package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	models "github.com/vnlab/makeshop-payment/src/domain/models"
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
		
		// Check if token is blacklisted
		if jwtService.IsBlacklisted(tokenString) {
			c.Next()
			return
		}
		
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Set authentication information in context
		c.Set("authenticated", true)
		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("roleId", claims.RoleID)
		c.Set("roleCode", claims.RoleCode)
		c.Set("token", tokenString) // Save token in context for logout

		c.Next()
	}
}

// WithAuth creates a GraphQL resolver context with auth information
func WithAuth(ctx context.Context, c *gin.Context) context.Context {
	for _, key := range []string{"authenticated", "userId", "email", "roleId", "roleCode", "token"} {
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

// GetUserID extracts the user ID from context
func GetUserID(ctx context.Context) (int, error) {
    if err := CheckAuth(ctx); err != nil {
        return 0, err
    }
    
    userID, ok := ctx.Value("userId").(int)
    if !ok {
        return 0, errors.New("user ID not found in context")
    }
    
    return userID, nil
}

// GetUserEmail extracts the user email from context
func GetUserEmail(ctx context.Context) (string, error) {
    if err := CheckAuth(ctx); err != nil {
        return "", err
    }
    
    email, ok := ctx.Value("email").(string)
    if !ok {
        return "", errors.New("user email not found in context")
    }
    
    return email, nil
}

// CheckRoleCode verifies if the user has the required role code
func CheckRoleCode(ctx context.Context, requiredCode string) error {
    if err := CheckAuth(ctx); err != nil {
        return err
    }
    
    roleCode, ok := ctx.Value("roleCode").(string)
    if !ok || roleCode != requiredCode {
        return fmt.Errorf("permission denied: %s role required", requiredCode)
    }
    
    return nil
}

// IsAdminRole checks if the authenticated user has admin role
func IsAdminRole(ctx context.Context) bool {
    roleCode, ok := ctx.Value("roleCode").(string)
    return ok && roleCode == string(models.RoleCodeAdmin)
}
