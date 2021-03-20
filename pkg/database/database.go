package database

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/MahanthMohan/GopherChat/pkg/schema"
	"github.com/fatih/color"
	"google.golang.org/api/option"
)

var (
	myCollection = "GopherChat"
	db           *firestore.Client
)

func init() {
	// Use the option WithCredentialsFile, to use the json file for firebase credentials
	opt := option.WithCredentialsFile("./credentials.json")

	// Initialize a new app
	app, err := firebase.NewApp(context.Background(), nil, opt)

	// Whenever there is an error, panic - stopping the goroutine
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}

	// Initialize a firestore database
	db, err = app.Firestore(context.Background())
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}
}

func CreateUserDocument(usr schema.User) {
	_, err := db.Collection(myCollection).Doc(usr.Username).Set(context.Background(), usr)
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}
}

func UpdateMemberStatus(username string, isGroupMember bool) {
	_, err := db.Collection(myCollection).Doc(username).Update(context.Background(), []firestore.Update{
		{
			Path:  "isGroupMember",
			Value: isGroupMember,
		},
	})
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}
}

func GetMemberStatus(username string) bool {
	docSnap, _ := db.Collection(myCollection).Doc(username).Get(context.Background())

	data, err := json.Marshal(docSnap.Data())
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}

	var userDocument schema.User
	json.Unmarshal(data, &userDocument)

	return userDocument.IsGroupMember
}

func SendUserMessage(reciever string, messages []string) {
	_, err := db.Collection(myCollection).Doc(reciever).Update(context.Background(), []firestore.Update{
		{
			Path:  "messages",
			Value: messages,
		},
	})

	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}
}

func GetAllMessages(username string) []interface{} {
	docSnap, err := db.Collection(myCollection).Doc(username).Get(context.Background())
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}

	data := docSnap.Data()
	messages := data["messages"].([]interface{})

	return messages
}

func ValidateUserLoginCredentials(username string, password string) bool {
	docSnap, err := db.Collection(myCollection).Doc(username).Get(context.Background())
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		print("\n** Username does not exist **")
	}

	data, err := json.Marshal(docSnap.Data())
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}

	var userDocument schema.User
	json.Unmarshal(data, &userDocument)

	if (username == userDocument.Username) && (password == userDocument.Password) {
		return true
	}

	return false
}

func GetAllUsernames() <-chan string {
	documents, err := db.Collection(myCollection).Documents(context.Background()).GetAll()
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}

	out := make(chan string, len(documents))

	go func() {
		for _, doc := range documents {
			data, err := json.Marshal(doc.Data())
			if err != nil {
				color.Set(color.FgHiRed, color.Bold)
				panic(err)
			}
			var userDocument schema.User
			json.Unmarshal(data, &userDocument)
			if len(userDocument.Username) != 0 {
				out <- userDocument.Username
			}
		}
		close(out)
	}()

	return out
}

func Close() {
	defer db.Close()
}
