package aut

import (
	"testing"
	"time"
	"github.com/google/uuid"
	"net/http"
)


func TestCheckHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	token := "Bearer hopuga"
	req.Header.Set("Authorization", token)
	
	req2, _ := http.NewRequest("GET", "", nil)

	tests := []struct {
		name    string
		token 	string
		header  http.Header
		wantErr bool
	}{
		{
			name:     "Correct token",
			token: 		token,
			header:     req.Header,
			wantErr:  	false,
		},
		{
			name:     "No header",
			token: 		"",
			header:     req2.Header,
			wantErr:  	true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func TestCheckPasswordJWT(t *testing.T) {
	user := uuid.New()
	expire, _ := time.ParseDuration("10h")
	password1 := "asdbg123"
	password2 := "lalaboska!"
	jwt1, _ := MakeJWT(user, password1, expire)
	jwt2, _ := MakeJWT(user, password2, expire)

	tests := []struct {
		name     string
		password string
		jwt     string
		user	uuid.UUID
		wantErr  bool
	}{
		{
			name:     "Correct jwt",
			password: 	password1,
			jwt:     	jwt1,
			user:		user,
			wantErr:  	false,
		},
		{
			name:     "Correct jwt 2",
			password: 	password2,
			jwt:     	jwt2,
			user:		user,
			wantErr:  	false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			jwt:     	jwt1,
			user:		user,
			wantErr:  	true,
		},
		{
			name:     "Password doesn't match different jwt",
			password: 	password1,
			jwt:     	jwt2,
			user:		user,
			wantErr:  	true,
		},
		{
			name:     "Invalid jwt",
			password: 	password1,
			jwt:     	"invalid jwt",
			user:		user,
			wantErr:  	true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateJWT(tt.jwt, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func TestCheckPasswordHash(t *testing.T) {
	password1 := "asdbg123"
	password2 := "chakachaka"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
