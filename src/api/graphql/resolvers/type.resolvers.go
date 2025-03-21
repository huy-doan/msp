package resolvers

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/api/graphql/generated"
	"github.com/vnlab/makeshop-payment/src/domain/models"
)

type typeResolver struct {
    *Resolver
}

// If you need to process specific fields, provide them as separate functions
// For example:
// func (r *Resolver) GetUserPosts(ctx context.Context, user *models.User) ([]*Post, error) {
//     // Custom logic here to fetch user's posts
// }

// User returns UserResolver implementation.
func (r *Resolver) User() generated.UserResolver {
    return &userResolver{r}
}

// Thêm struct userResolver
type userResolver struct {
    *Resolver
}

// MFA implementation
func (r *userResolver) MfaType(ctx context.Context, obj *models.User) (*generated.MFAType, error) {
    // Nếu user không có MFA type được bật
    if obj.MFATypeID == nil || obj.MFAType == nil {
        return nil, nil
    }
    
    // convert from models.MFAType to generated.MFAType
    return &generated.MFAType{
        ID:        obj.MFAType.ID,
        No:        obj.MFAType.No,
        Title:     obj.MFAType.Title,
        IsActive:  obj.MFAType.IsActive,
        CreatedAt: obj.MFAType.CreatedAt,
        UpdatedAt: obj.MFAType.UpdatedAt,
    }, nil
}
