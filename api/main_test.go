package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	regen "github.com/zach-klippenstein/goregen"
	"net/http"
	"net/http/httptest"
	"testing"
	"xm-task/api"
	"xm-task/conf"
	"xm-task/dao"
	"xm-task/services"
	"xm-task/smodels"
)

func TestCreateCompanyIntegration(t *testing.T) {
	cfg := conf.Config{
		API: conf.API{
			ListenOnPort:       7777,
			CORSAllowedOrigins: []string{"*"},
		},
		Postgres: conf.Postgres{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "1234",
			Database: "xm-test",
			SSLMode:  "disable",
		},
	}
	d, err1 := dao.New(cfg, false)
	require.NoError(t, err1)

	service, err2 := services.NewService(cfg, d)
	require.NoError(t, err2)

	a, err3 := api.NewAPI(cfg, service)
	require.NoError(t, err3)

	ts := httptest.NewServer(a.Router())
	defer ts.Close()

	// checks protection of unauthorized access
	t.Run("it should return 401 status code", func(t *testing.T) {
		company := smodels.Company{
			Name:       "Test Company",
			Employees:  100,
			Registered: true,
			Type:       "Corporation",
		}
		requestBody, err4 := json.Marshal(company)
		require.NoError(t, err4)

		resp, err5 := http.Post(ts.URL+"/auth/companies", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err5)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// checks that users with real access token can create companies and auto-creating uuid
	t.Run("it should return 200 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@gmail.com",
			Password: "23rwgfds",
		}
		requestBody, err6 := json.Marshal(user)
		require.NoError(t, err6)

		userResp, err7 := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err7)
		defer userResp.Body.Close()

		var tt smodels.TestTokenDetails
		err8 := json.NewDecoder(userResp.Body).Decode(&tt)
		require.NoError(t, err8)

		name, err88 := regen.Generate(fmt.Sprintf("[a-zA-Z0-9]{%d,%d}", 15, 15))
		require.NoError(t, err88)
		company := smodels.Company{
			ID:          "0753913b-8910-40de-827f-6c0085dec47e",
			Name:        name,
			Description: "some description",
			Employees:   100,
			Registered:  true,
			Type:        "Corporations",
		}
		requestBody2, err9 := json.Marshal(company)
		require.NoError(t, err9)

		r, err := http.NewRequest("POST", ts.URL+"/auth/companies", bytes.NewBuffer(requestBody2))
		require.NoError(t, err)
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.AccessToken))
		r.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err10 := client.Do(r)
		require.NoError(t, err10)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// checks protection of non-existing types
	t.Run("it should return 500 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@gmail.com",
			Password: "23rwgfds",
		}
		requestBody, err11 := json.Marshal(user)
		require.NoError(t, err11)

		userResp, err12 := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err12)
		defer userResp.Body.Close()

		var tt smodels.TestTokenDetails
		err13 := json.NewDecoder(userResp.Body).Decode(&tt)
		require.NoError(t, err13)

		name, err131 := regen.Generate(fmt.Sprintf("[a-zA-Z0-9]{%d,%d}", 15, 15))
		require.NoError(t, err131)
		company := smodels.Company{
			Name:        name,
			Description: "some description",
			Employees:   100,
			Registered:  true,
			Type:        "random value",
		}
		requestBody2, err14 := json.Marshal(company)
		require.NoError(t, err14)

		r, err15 := http.NewRequest("POST", ts.URL+"/auth/companies", bytes.NewBuffer(requestBody2))
		require.NoError(t, err15)
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.AccessToken))
		r.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err16 := client.Do(r)
		require.NoError(t, err16)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestCreateAndGetAndDeleteCompanyIntegration(t *testing.T) {
	cfg := conf.Config{
		API: conf.API{
			ListenOnPort:       7777,
			CORSAllowedOrigins: []string{"*"},
		},
		Postgres: conf.Postgres{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "1234",
			Database: "xm-test",
			SSLMode:  "disable",
		},
	}
	d, err1 := dao.New(cfg, false)
	require.NoError(t, err1)

	service, err2 := services.NewService(cfg, d)
	require.NoError(t, err2)

	a, err3 := api.NewAPI(cfg, service)
	require.NoError(t, err3)

	ts := httptest.NewServer(a.Router())
	defer ts.Close()

	// checks that users with real access token can create companies and auto-creating uuid
	// and then test GET and DELETE endpoint
	t.Run("it should return 200 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@gmail.com",
			Password: "23rwgfds",
		}
		requestBody, err4 := json.Marshal(user)
		require.NoError(t, err4)

		userResp, err5 := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err5)
		defer userResp.Body.Close()
		assert.Equal(t, http.StatusOK, userResp.StatusCode)

		var tt smodels.TestTokenDetails
		err6 := json.NewDecoder(userResp.Body).Decode(&tt)
		require.NoError(t, err6)

		name, err7 := regen.Generate(fmt.Sprintf("[a-zA-Z0-9]{%d,%d}", 15, 15))
		require.NoError(t, err7)
		company := smodels.Company{
			Name:        name,
			Description: "some description",
			Employees:   100,
			Registered:  true,
			Type:        "Corporations",
		}
		requestBody2, err8 := json.Marshal(company)
		require.NoError(t, err8)

		r, err9 := http.NewRequest("POST", ts.URL+"/auth/companies", bytes.NewBuffer(requestBody2))
		require.NoError(t, err9)
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.AccessToken))
		r.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err10 := client.Do(r)
		require.NoError(t, err10)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var createdCompany smodels.Company
		err11 := json.NewDecoder(resp.Body).Decode(&createdCompany)
		require.NoError(t, err11)

		getResp, err12 := http.Get(ts.URL + "/companies/" + createdCompany.ID)
		require.NoError(t, err12)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		dr, err13 := http.NewRequest("DELETE", ts.URL+"/auth/companies/"+createdCompany.ID, nil)
		require.NoError(t, err13)
		dr.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.AccessToken))
		delResp, err14 := client.Do(dr)
		require.NoError(t, err14)
		defer delResp.Body.Close()
		assert.Equal(t, http.StatusOK, delResp.StatusCode)
	})
}

func TestSignInIntegration(t *testing.T) {
	cfg := conf.Config{
		API: conf.API{
			ListenOnPort:       7777,
			CORSAllowedOrigins: []string{"*"},
		},
		Postgres: conf.Postgres{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "1234",
			Database: "xm-test",
			SSLMode:  "disable",
		},
	}
	d, err1 := dao.New(cfg, false)
	require.NoError(t, err1)

	service, err2 := services.NewService(cfg, d)
	require.NoError(t, err2)

	a, err3 := api.NewAPI(cfg, service)
	require.NoError(t, err3)

	ts := httptest.NewServer(a.Router())
	defer ts.Close()

	// checks proper sign-in/register
	t.Run("it should return 200 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@gmail.com",
			Password: "23rwgfds",
		}
		requestBody, err4 := json.Marshal(user)
		require.NoError(t, err4)

		resp, err5 := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err5)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// checks incorrect password
	t.Run("it should return 401 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@gmail.com",
			Password: "mypass",
		}
		requestBody, err6 := json.Marshal(user)
		require.NoError(t, err6)

		resp, err7 := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err7)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// checks incorrect email
	t.Run("it should return 400 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@.com",
			Password: "23rwgfds",
		}
		requestBody, err8 := json.Marshal(user)
		require.NoError(t, err8)

		resp, err9 := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err9)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// checks too long password
	t.Run("it should return 500 status code (too long pass)", func(t *testing.T) {
		user := smodels.User{
			Email:    "newuser@gmail.com",
			Password: "237fuhaerug7fhu347fhwi87r8e7yw8yf7wyihyfkwiw5iwefhhfihisfiehisefhihsefihfeie4rye7y4rre4djsijdj",
		}
		requestBody, err10 := json.Marshal(user)
		require.NoError(t, err10)

		resp, err11 := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err11)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
