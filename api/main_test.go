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
			Host:     "postgres-db",
			Port:     "5432",
			User:     "postgres",
			Password: "somesecretpassword1234",
			Database: "xm-test-db",
			SSLMode:  "disable",
		},
	}
	d, err := dao.New(cfg, false)
	require.NoError(t, err)

	service, err := services.NewService(cfg, d)
	require.NoError(t, err)

	a, err := api.NewAPI(cfg, service)
	require.NoError(t, err)

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
		requestBody, err := json.Marshal(company)
		require.NoError(t, err)

		resp, err := http.Post(ts.URL+"/auth/companies", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// checks that users with real access token can create companies and auto-creating uuid
	t.Run("it should return 200 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@gmail.com",
			Password: "23rwgfds",
		}
		requestBody, err := json.Marshal(user)
		require.NoError(t, err)

		userResp, err := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err)
		defer userResp.Body.Close()

		var tt smodels.TestTokenDetails
		err = json.NewDecoder(userResp.Body).Decode(&tt)
		require.NoError(t, err)

		name, errReg := regen.Generate(fmt.Sprintf("[a-zA-Z0-9]{%d,%d}", 15, 15))
		require.NoError(t, errReg)
		company := smodels.Company{
			ID:          "0753913b-8910-40de-827f-6c0085dec47e",
			Name:        name,
			Description: "some description",
			Employees:   100,
			Registered:  true,
			Type:        "Corporations",
		}
		requestBody, err = json.Marshal(company)
		require.NoError(t, err)

		r, err := http.NewRequest("POST", ts.URL+"/auth/companies", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.AccessToken))
		r.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(r)
		require.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// checks protection of non-existing types
	t.Run("it should return 500 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@gmail.com",
			Password: "23rwgfds",
		}
		requestBody, err := json.Marshal(user)
		require.NoError(t, err)

		userResp, err := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err)
		defer userResp.Body.Close()

		var tt smodels.TestTokenDetails
		err = json.NewDecoder(userResp.Body).Decode(&tt)
		require.NoError(t, err)

		name, errReg := regen.Generate(fmt.Sprintf("[a-zA-Z0-9]{%d,%d}", 15, 15))
		require.NoError(t, errReg)
		company := smodels.Company{
			Name:        name,
			Description: "some description",
			Employees:   100,
			Registered:  true,
			Type:        "random value",
		}
		requestBody, err = json.Marshal(company)
		require.NoError(t, err)

		r, err := http.NewRequest("POST", ts.URL+"/auth/companies", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.AccessToken))
		r.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(r)
		require.NoError(t, err)

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
	d, err := dao.New(cfg, false)
	require.NoError(t, err)

	service, err := services.NewService(cfg, d)
	require.NoError(t, err)

	a, err := api.NewAPI(cfg, service)
	require.NoError(t, err)

	ts := httptest.NewServer(a.Router())
	defer ts.Close()

	// checks that users with real access token can create companies and auto-creating uuid
	// and then test GET and DELETE endpoint
	t.Run("it should return 200 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@gmail.com",
			Password: "23rwgfds",
		}
		requestBody, err := json.Marshal(user)
		require.NoError(t, err)

		userResp, err := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err)
		defer userResp.Body.Close()
		assert.Equal(t, http.StatusOK, userResp.StatusCode)

		var tt smodels.TestTokenDetails
		err = json.NewDecoder(userResp.Body).Decode(&tt)
		require.NoError(t, err)

		name, errReg := regen.Generate(fmt.Sprintf("[a-zA-Z0-9]{%d,%d}", 15, 15))
		require.NoError(t, errReg)
		company := smodels.Company{
			Name:        name,
			Description: "some description",
			Employees:   100,
			Registered:  true,
			Type:        "Corporations",
		}
		requestBody, err = json.Marshal(company)
		require.NoError(t, err)

		r, err := http.NewRequest("POST", ts.URL+"/auth/companies", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.AccessToken))
		r.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(r)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var createdCompany smodels.Company
		err = json.NewDecoder(resp.Body).Decode(&createdCompany)
		require.NoError(t, err)

		getResp, err := http.Get(ts.URL + "/companies/" + createdCompany.ID)
		require.NoError(t, err)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		dr, err := http.NewRequest("DELETE", ts.URL+"/auth/companies/"+createdCompany.ID, nil)
		require.NoError(t, err)
		dr.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.AccessToken))
		delResp, err := client.Do(dr)
		require.NoError(t, err)
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
	d, err := dao.New(cfg, false)
	require.NoError(t, err)

	service, err := services.NewService(cfg, d)
	require.NoError(t, err)

	a, err := api.NewAPI(cfg, service)
	require.NoError(t, err)

	ts := httptest.NewServer(a.Router())
	defer ts.Close()

	// checks proper sign-in/register
	t.Run("it should return 200 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@gmail.com",
			Password: "23rwgfds",
		}
		requestBody, err := json.Marshal(user)
		require.NoError(t, err)

		resp, err := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// checks incorrect password
	t.Run("it should return 401 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@gmail.com",
			Password: "mypass",
		}
		requestBody, err := json.Marshal(user)
		require.NoError(t, err)

		resp, err := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// checks incorrect email
	t.Run("it should return 400 status code", func(t *testing.T) {
		user := smodels.User{
			Email:    "email12@.com",
			Password: "23rwgfds",
		}
		requestBody, err := json.Marshal(user)
		require.NoError(t, err)

		resp, err := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// checks too long password
	t.Run("it should return 500 status code (too long pass)", func(t *testing.T) {
		user := smodels.User{
			Email:    "newuser@gmail.com",
			Password: "237fuhaerug7fhu347fhwi87r8e7yw8yf7wyihyfkwiw5iwefhhfihisfiehisefhihsefihfeie4rye7y4rre4djsijdj",
		}
		requestBody, err := json.Marshal(user)
		require.NoError(t, err)

		resp, err := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(requestBody))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
