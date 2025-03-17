package graphql

import (
	"github.com/gin-gonic/gin"
	"github.com/vnlab/makeshop-payment/src/application/services"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
)

// GraphQLQuery godoc
// @Summary GraphQL API endpoint
// @Description Endpoint chính cho GraphQL API queries và mutations
// @Tags graphql
// @Accept json
// @Produce json
// @Param query body string true "GraphQL query"
// @Success 200 {object} map[string]interface{} "GraphQL response"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /graphql/query [post]
func GraphQLInfo(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "GraphQL endpoint đang được phát triển",
    })
}

// SetupGraphQL configures GraphQL handlers for the given Gin router
func SetupGraphQL(
	router *gin.Engine,
	userService *services.UserService,
	jwtService *auth.JWTService,
) {
	// Phần GraphQL đang được phát triển
	// Tạm thời tắt để ứng dụng có thể biên dịch
	// Thêm route stub để đánh dấu tính năng đang trong quá trình phát triển
	router.POST("/graphql", GraphQLInfo)
}
