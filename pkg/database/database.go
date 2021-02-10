package database

import (
	"cloud.google.com/go/firestore"
	"context"
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

func SendUserMessage(reciever string, messages []schema.Message) {
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

func GetAllMessages(username string) []schema.Message {
	docSnap, err := db.Collection(myCollection).Doc(username).Get(context.Background())
	if err != nil {
		panic(err)
	}
	messages, err := docSnap.DataAt("messages")
	if err != nil {
		panic(err)
	}
	return messages.([]schema.Message)
}

func GetAllUsernames() []string {
	var names []string
	documents, err := db.Collection(myCollection).Documents(context.Background()).GetAll()
	if err != nil {
		panic(err)
	}
	for _, doc := range documents {
		username, err := doc.DataAt("username")
		if err != nil {
			panic(err)
		}
		names = append(names, username.(string))
	}

	return names
}
