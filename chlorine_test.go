package chlorine

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
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

//var token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjIwMTNlNTNlLTUyZGQtNWUxYy1hZjBiLWI1MDNlMzFjOGE1OSIsImVtYWlsIjoiYnJpbGxpYW50LnlhbmdAYmFkYW5hbXUuY29tLmNuIiwiZXhwIjoxNjA2MjAyOTg4LCJpc3MiOiJraWRzbG9vcCJ9.Njz8v5A6pMjNPi3aDoDUSGmYnUYKibv44sTxQtd3PSrzXbWya28qVx_OpUON1UvdIixpy_HD61nmial0C3bnWTb6F9256cmPA9w_GypvC14YCm5jE4-HsOhtGZIVkrQKp4DXP8G1nujC9o_YCR6z5OemfnITmPoPdoO44OCbbHiF2IHNJfpTqFVGvUd32ZXhGf1Njpsk9Q2bFWeyYyNQan8raTVvDsWNlHp_UzgsBVopyRM1BlqO8te6z6mWEt_g851HKPSauPZGZFD1RHq351Lhg5YogsXF0eNf-n6TCujAFs54kYJDG2Q9pJHjsGfYkmg8K6yHzugO9KpJ8Xyd6g"
var token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjE1MGRlMDRiLTY3NzctNGNhYi1hMzgxLTVmMzZhOWI3NTBlZiIsImVtYWlsIjoiYnJpbGxpYW50LnlhbmdAYmFkYW5hbXUuY29tLmNuIiwiZXhwIjoxNjI5OTQ1OTIxLCJpc3MiOiJraWRzbG9vcCJ9.v6HxDf6q002e8vlBHtxGDjFjQYwuKDqxo0lnc9nuy2xuJuu9b9K2RAk76J76Jqkt6gNE9w7laJp4EFq_X9cAliGF3VCeEkXevQFfW67XFozldFbTnRo4v54AoG6zkkOeK4s4PVfp1xQU3hqbLM5n1DA5EeFNsJDhwH6Vb5eulUZg6CD3FO_4FqwI3FhvkDGMFx9Psw7xxDfwT_M1vubOXKvu1lRUbnf5gW6pT_0ppli1rS9da-4E8YmZGpePx52dgUNwWXgu21o1NU158AQ5rEkmIyTzZq_WBeHr-fb3vR_PZQUCgdbvlRuSvr3T3bOoGrT4HbeC_UKb6cuciI_khQ"
var client = NewClient("https://api.alpha.kidsloop.net/user/", WithTimeout(time.Minute))

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
	req := NewRequest(q, ReqToken(token))
	req.Var("userID", "be8ca64d-105d-4551-9b15-5d8fb2585b50")
	//req.Var("userID", "be8ca64d-105d-4551-9b15-5d8fb2585b51")
	//req.Var("userID", "1")
	//req.Var("userIID", "1")
	statusCode, err := client.Run(context.Background(), req, &resp)
	if err != nil {
		t.Fatal(statusCode, err)
	}
	result, _ := json.Marshal(resp)
	fmt.Println(statusCode, string(result))
}

func TestQueryOrgBatchDemo1(t *testing.T) {
	q := `query{
   org0: organization(organization_id: "44fa30db-365a-4589-a304-3c8a801debb5") {organization_name}
   org1: organization(organization_id: "e236d102-5324-4740-8f36-629451557a2a") {organization_name}
   org2: organization(organization_id: "3f135b91-a616-4c80-914a-e4463104dbac") {organization_name}
   org3: organization(organization_id: "c70a525e-e62d-41a0-85a2-91ac9b707a53") {organization_name}
   org4: organization(organization_id: "25b9c6dd-21e8-4718-8f90-d71e124684a8") {organization_name}
   org5: organization(organization_id: "66d85eab-9e15-4d89-9e9d-f4ed37d254dd") {organization_name}
   org6: organization(organization_id: "65f4a766-dc2f-4dad-a08c-d1d4ac02fcf1") {organization_name}
   org7: organization(organization_id: "9657e271-1324-4fc5-a371-46748481d664") {organization_name}
}`
	req := NewRequest(q)
	resp := Response{
		Data: map[string]*struct {
			OrgID   string `json:"organization_id"`
			OrgName string `json:"organization_name"`
		}{},
	}

	_, err := client.Run(context.Background(), req, &resp)
	if err != nil {
		t.Fatal(err)
	}
	result, _ := json.Marshal(resp)
	fmt.Println(string(result))
}

func TestQueryOrgBatchDemo2(t *testing.T) {
	q := `query($ids: [ID!]){
	organizations(organization_ids: $ids){
		organization_id
		organization_name
	}
}`
	req := NewRequest(q)
	req.Var("ids", []string{"44fa30db-365a-4589-a304-3c8a801debb5", "e236d102-5324-4740-8f36-629451557a2a"})
	resp := Response{
		Data: map[string]*struct {
			OrgID   string `json:"organization_id"`
			OrgName string `json:"organization_name"`
		}{},
	}

	_, err := client.Run(context.Background(), req, &resp)
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
	_, err := client.Run(context.Background(), req, &resp)
	if err != nil {
		t.Fatal(err)
	}
	result, _ := json.Marshal(resp)
	fmt.Println(string(result))
}

func TestRequest_SetHeader(t *testing.T) {
	req := NewRequest("")
	req.SetHeader("access", "token")
	fmt.Println(req)
}

func TestRequest_SetHeaders(t *testing.T) {
	req := NewRequest("")
	req.SetHeaders("cookie", []string{"access=token"})
	fmt.Println(req)
}
