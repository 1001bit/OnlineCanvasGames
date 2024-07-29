package token

func RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(claims.UserID, claims.Username)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
