package chlorine

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
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
	_, err := client.Run(context.Background(), req, &resp)
	if err != nil {
		t.Fatal(err)
	}
	result, _ := json.Marshal(resp)
	fmt.Println(string(result))
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
<<<<<<< HEAD
	q := `query($ids: [ID!]){
=======
	q := ` query($ids: [ID!]){
>>>>>>> 1bcc76ccd1becf2ad5725ebf7ef79449344e5439
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
