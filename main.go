package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

func printError(err error) {
	log.Output(2, err.Error())
}

func JSON(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		printError(err)
		return
	}
}

// domain model
type User struct {
	ID   int    `json:"id"`
	Name string `json:"userName"`
}

// graphql output difinition
var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// root query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"users": &graphql.Field{
			Type:        graphql.NewList(UserType),
			Description: "get users",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// TODO use infra function to return response
				return Users(), nil
			},
		},
	},
})

// infra
func Users() []User {
	return []User{
		{ID: 1, Name: "gorilla"},
		{ID: 2, Name: "cat"},
		{ID: 3, Name: "dog"},
		{ID: 4, Name: "bird"},
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		printError(err)
		return
	}

	// new schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		printError(err)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: string(body),
	})

	if len(result.Errors) > 0 {
		printError(fmt.Errorf("wrong result, unexpected errors: %v", result.Errors))
		return
	}

	JSON(w, result.Data)
}

func main() {
	http.HandleFunc("/", handle)
	log.Println("server start")
	log.Fatal(http.ListenAndServe(":80", nil))
}
