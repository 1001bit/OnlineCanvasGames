package auth

import "errors"

var ErrNoClaims = errors.New("no claims found in context")
