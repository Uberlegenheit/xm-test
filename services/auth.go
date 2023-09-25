package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
	"strings"
	"time"
	"xm-task/smodels"
)

func (s *ServiceFacade) CreateToken(email string) (smodels.TokenDetails, error) {
	var td smodels.TokenDetails

	val, _, err := s.dao.GetAuthToken(fmt.Sprintf("%s_td", email))
	if err != nil {
		return smodels.TokenDetails{}, err
	}

	if val != nil {
		td = val.(smodels.TokenDetails)
		if td.RtExpires != 0 {
			return td, nil
		}
	}

	td.AtExpires = time.Now().Add(time.Minute * 30).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["email"] = email
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return smodels.TokenDetails{}, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["email"] = email
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return smodels.TokenDetails{}, err
	}

	return td, nil
}

func (s *ServiceFacade) CreateAuth(email string, td smodels.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	err := s.dao.AddAuthToken(fmt.Sprintf("%s_td", email), td, at.Sub(now))
	if err != nil {
		return err
	}
	err = s.dao.AddAuthToken(td.AccessUuid, email, at.Sub(now))
	if err != nil {
		return err
	}
	err = s.dao.AddAuthToken(td.RefreshUuid, email, rt.Sub(now))
	if err != nil {
		return err
	}
	err = s.dao.AddAuthToken(fmt.Sprintf("%s_access", td.RefreshUuid), td.AccessUuid, rt.Sub(now))
	if err != nil {
		return err
	}
	err = s.dao.AddAuthToken(fmt.Sprintf("%s_refresh", td.AccessUuid), td.RefreshUuid, rt.Sub(now))
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceFacade) ExtractTokenMetadata(c *gin.Context) (smodels.AccessDetails, error) {
	token, err := s.VerifyToken(c.Request)
	if err != nil {
		return smodels.AccessDetails{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return smodels.AccessDetails{}, err
		}

		email := fmt.Sprintf("%s", claims["email"])
		return smodels.AccessDetails{
			AccessUuid: accessUuid,
			Email:      email,
		}, nil
	}
	return smodels.AccessDetails{}, err
}

func (s *ServiceFacade) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	if tokenString == "" {
		return nil, fmt.Errorf("no token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *ServiceFacade) Refresh(r *http.Request) (smodels.TokenDetails, error) {
	refreshToken := r.Header.Get("Authorization")
	parts := strings.Split(refreshToken, " ")
	if len(parts) != 2 {
		return smodels.TokenDetails{}, fmt.Errorf("error: %s", "cannot get the refresh token")
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
	})
	if err != nil {
		return smodels.TokenDetails{}, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return smodels.TokenDetails{}, fmt.Errorf("error: %s", "invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return smodels.TokenDetails{}, fmt.Errorf("error: %s", "invalid token")
		}
		email := fmt.Sprintf("%s", claims["email"])
		accessUuid, ok, err := s.dao.GetAuthToken(fmt.Sprintf("%s_access", refreshUuid))
		if err != nil || !ok {
			return smodels.TokenDetails{}, fmt.Errorf("error: %s", "cannot get access token: invalid refresh_access token")
		}
		err = s.DeleteAuth(refreshUuid,
			accessUuid.(string),
			fmt.Sprintf("%s_access", refreshUuid),
			fmt.Sprintf("%s_refresh", accessUuid.(string)),
			fmt.Sprintf("%s_td", email),
		)
		if err != nil {
			return smodels.TokenDetails{}, fmt.Errorf("error: %s", "invalid token provided")
		}

		ts, err := s.CreateToken(email)
		if err != nil {
			return smodels.TokenDetails{}, fmt.Errorf("error: %s", "cannot create token")
		}

		saveErr := s.CreateAuth(email, ts)
		if saveErr != nil {
			return smodels.TokenDetails{}, fmt.Errorf("error: %s", "cannot create auth")
		}

		return ts, nil
	}

	return smodels.TokenDetails{}, fmt.Errorf("error: %s", "cannot refresh tokens")
}

func extractToken(r *http.Request) string {
	accessToken := r.Header.Get("Authorization")
	token := strings.Split(accessToken, " ")
	if len(token) != 2 {
		return ""
	} else if token[0] != "Bearer" {
		return ""
	}
	return token[1]
}
