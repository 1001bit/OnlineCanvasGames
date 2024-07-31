package token

func RefreshTokens(refreshToken string) (string, string, error) {
	// validate refresh token
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// generate new tokens
	newAccessToken, err := GenerateAccessToken(claims.Username)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(claims.Username)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
