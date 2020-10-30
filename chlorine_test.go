package chlorine

import (
	"context"
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
	resp := struct {
		Data struct {
			User struct {
				UserID   string `json:"user_id"`
				UserName string `json:"user_name"`
			} `json:"user"`
		} `json:"data,omitempty"`
		Errors []ClError `json:"errors,omitempty"`
	}{}
	req := NewRequest(q)
	//req.Var("userID", "be8ca64d-105d-4551-9b15-5d8fb2585b50")
	req.Var("userID", "1")
	//req.Var("userIID", "1")
	err := client.Run(context.Background(), req, &resp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestUpdateUser(t *testing.T) {
	q := `mutation update_user($userID: ID!, $userName: String){
	user(user_id:$userID, user_name:$userName){
		user_id
		user_name
	}
}
`
	resp := struct {
		Data struct {
			User struct {
				UserID   string `json:"user_id"`
				UserName string `json:"user_name"`
			} `json:"user"`
		} `json:"data,omitempty"`
		Errors []ClError `json:"errors,omitempty"`
	}{}
	req := NewRequest(q)
	req.Var("userID", "be8ca64d-105d-4551-9b15-5d8fb2585b50")
	req.Var("userName", "PJ")
	//req.Var("userID", "1")
	err := client.Run(context.Background(), req, &resp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}
