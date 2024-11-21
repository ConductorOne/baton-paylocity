package models

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Local Variables:
// go-tag-args: ("-transform" "snakecase")
// End:
