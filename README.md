# xm-test
Test project for XM

## Installation
Open your console.
To install this project, all you need is just download this repo:
`git clone https://github.com/Uberlegenheit/xm-test`.
Move to created folder `cd xm-test` and then call the command: `make deploy`.


## Public endpoints
### / (GET)
Returns the name of backend service

### /health (GET)
Returns status of db

### /companies/:id (GET)
Returns company with specified id

### /sign-in (POST)
Register new user or return existing. Returns access and refresh tokens.
Example:
```json
{
"email": "email@gmail.com",
"password": "2378fyu9u8we8u8"
}
```

### /refresh (POST)
Send Authorization header set with `Bearer *your_refresh_token*`.

## Protected routes
All of them require Authorization header set with `Bearer *your_access_token*`.

### /auth/companies (POST)
Creates new company.
```json
{
  "name": "someName",
  "description": "description",
  "employees_count": 23,
  "registered": true,
  "type": "Cooperative"
}
```

### /auth/companies/:id (PATCH)
Similar to creation, but requires id specified in uri.
```json
{
  "name": "someName",
  "description": "updated",
  "employees_count": 23,
  "registered": true,
  "type": "Cooperative"
}
```

### /auth/companies/:id (DELETE)
Delete company with id specified in uri.
```json
{
  "name": "someName",
  "description": "updated",
  "employees_count": 23,
  "registered": true,
  "type": "Cooperative"
}
```

### /auth/logout (POST)
Delete active session.

## Tests
To run tests, you have to be in root folder and run command `go test ./api/main_test.go`