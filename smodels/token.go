package smodels

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string
	Email      string
}

type TestTokenDetails struct {
	AccessToken string `json:"access_token"`
}
