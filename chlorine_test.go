package chlorine

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	UserID   string
	UserName string
	Orgs     []Organization
}

type Organization struct {
	OrgID   string
	OrgName string
}

var client = NewClient("https://api.beta.kidsloop.net/user/")

func TestQueryUser(t *testing.T) {
	q := `query user($userID: ID!){
	user(user_id:$userID){
		user_id
		user_name
	}
}
`
	resp := Response{
		Data: &struct {
			User *struct {
				UserID   string `json:"user_id"`
				UserName string `json:"user_name"`
			} `json:"user,omitempty"`
		}{},
	}
	req := NewRequest(q)
	req.Var("userID", "be8ca64d-105d-4551-9b15-5d8fb2585b50")
	//req.Var("userID", "be8ca64d-105d-4551-9b15-5d8fb2585b51")
	//req.Var("userID", "1")
	//req.Var("userIID", "1")
	err := client.Run(context.Background(), req, &resp)
	if err != nil {
		t.Fatal(err)
	}
	result, _ := json.Marshal(resp)
	fmt.Println(string(result))
}

func TestUpdateUser(t *testing.T) {
	q := `mutation update_user($userID: ID!, $userName: String){
	user(user_id:$userID, user_name:$userName){
		user_id
		user_name
	}
}
`
	resp := Response{
		Data: &struct {
			User *struct {
				UserID   string `json:"user_id"`
				UserName string `json:"user_name"`
			} `json:"user,omitempty"`
		}{},
	}
	req := NewRequest(q)
	req.Var("userID", "be8ca64d-105d-4551-9b15-5d8fb2585b50")
	req.Var("userName", "PJ")
	//req.Var("userID", "1")
	err := client.Run(context.Background(), req, &resp)
	if err != nil {
		t.Fatal(err)
	}
	result, _ := json.Marshal(resp)
	fmt.Println(string(result))
}

const publicKeyOwen = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxdHMYTqFobj3oGD/JDYb
DN07icTH/Dj7jBtJSG2clM6hQ1HRLApQUNoqcrcJzA0A7aNqELIJuxMovYAoRtAT
E1pYMWpVyG41inQiJjKFyAkuHsVzL+t2C778BFxlXTC/VWoR6CowWSWJaYlT5fA/
krUew7/+sGW6rjV2lQqxBN3sQsfaDOdN5IGkizsfMpdrETbc5tKksNs6nL6SFRDe
LoS4AH5KI4T0/HC53iLDjgBoka7tJuu3YsOBzxDX22FbYfTFV7MmPyq++8ANbzTL
sgaD2lwWhfWO51cWJnFIPc7gHBq9kMqMK3T2dw0jCHpA4vYEMjsErNSWKjaxF8O/
FwIDAQAB
-----END PUBLIC KEY-----`

var token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImI0MjI3NDNlLTllYmMtNWVkNy1iNzI1LTA2Mjk5NGVjNzdmMiIsImVtYWlsIjoiYnJpbGxpYW50LnlhbmdAYmFkYW5hbXUuY29tLmNuIiwiZXhwIjoxNjA1MTc4Nzk2LCJpc3MiOiJraWRzbG9vcCJ9.sDkGFTIWm-NgEDfNJoMS_3KoKcZs0smnR7whqWY0AMnYLFYX3j_Saj6gHjXHpmZMVewbnaNfv9lYfhSokFBZaCcYyeVBXQo6DHL6nppsMUFwmcTjl-NjqSGwYUvjpV7cmkmL33H8KojEuBUDP8kOK-cF5Km28PC6sV2nFRVBNFBXlcNsdB-CIQEeycCzRhw078GAP64Bpugay8W-77keldN-C1Qnrc6spbSCOKnxMpT94pBRzgB8D-vHdcnvB3zlfPj8RYWFlGE_uufHfPTSgS-nTzrz8vRhiJdOAYdPys90w87jGfmopm1AT-qDSqa4Qf8hMW4bj_UDAa4-1bI-yQ"

func TestParseJWT(t *testing.T) {
	claims := &struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		*jwt.StandardClaims
	}{}
	_, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyOwen))
	})
	if err != nil {
		t.Fatal(err)
	}
	marshal, err := json.Marshal(claims)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(marshal))
}
