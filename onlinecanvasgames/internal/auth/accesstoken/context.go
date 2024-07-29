package accesstoken

import "context"

type claimsKeyType string

var ClaimsKey claimsKeyType = "claims"

// get new context from context and claims
func ContextWithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, ClaimsKey, claims)
}