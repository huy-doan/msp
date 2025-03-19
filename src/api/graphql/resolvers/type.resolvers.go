package resolvers

type typeResolver struct {
    *Resolver
}

// If you need to process specific fields, provide them as separate functions
// For example:
// func (r *Resolver) GetUserPosts(ctx context.Context, user *entities.User) ([]*Post, error) {
//     // Custom logic here to fetch user's posts
// }