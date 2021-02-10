package database

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/MahanthMohan/GopherChat/pkg/schema"
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
		panic(err)
	}

	// Initialize a firestore database
	db, err = app.Firestore(context.Background())
	if err != nil {
		panic(err)
	}
}

func CreateUserDocument(usr schema.User) {
	_, err := db.Collection(myCollection).Doc(usr.Username).Set(context.Background(), usr)
	if err != nil {
		panic(err)
	}
}

func UpdateMemberStatus(usr schema.User, isGroupMember bool) {
	_, err := db.Collection(myCollection).Doc(usr.Username).Update(context.Background(), []firestore.Update{
		{
			Path:  "isGroupMember",
			Value: isGroupMember,
		},
	})
	if err != nil {
		panic(err)
	}
}

func SendUserMessage(reciever string, messages []string) {
	_, err := db.Collection(myCollection).Doc(reciever).Update(context.Background(), []firestore.Update{
		{
			Path:  "messages",
			Value: messages,
		},
	})

	if err != nil {
		panic(err)
	}
}

func GetAllMessages(username string) []interface{} {
	docSnap, err := db.Collection(myCollection).Doc(username).Get(context.Background())
	if err != nil {
		panic(err)
	}
	data := docSnap.Data()
	messages := data["messages"].([]interface{})
	return messages
}

func ValidateUserLoginCredentials(username string, password string) bool {
	var ret bool
	docSnap, err := db.Collection(myCollection).Doc(username).Get(context.Background())
	if err != nil {
		panic(err)
	}
	data := docSnap.Data()
	actualUsername, actualPassword := data["username"].(string), data["password"].(string)

	if (username == actualUsername) && (password == actualPassword) {
		ret = true
	}

	return ret
}

func GetAllUsernames() []string {
	var names []string
	documents, err := db.Collection(myCollection).Documents(context.Background()).GetAll()
	if err != nil {
		panic(err)
	}
	for _, doc := range documents {
		data := doc.Data()
		username := data["username"].(string)
		if err != nil {
			panic(err)
		}
		names = append(names, username)
	}

	return names
}
