package database

import (
	"context"

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
	docSnap, err := db.Collection(myCollection).Doc(username).Get(context.Background())
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}

	data := docSnap.Data()
	memberStatus := data["isGroupMember"].(bool)

	return memberStatus
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
	ret := false
	docSnap, err := db.Collection(myCollection).Doc(username).Get(context.Background())
	if err != nil {

		color.Set(color.FgHiRed, color.Bold)
		println("** Username does not exist **")
	} else {
		data := docSnap.Data()
		actualUsername, actualPassword := data["username"].(string), data["password"].(string)

		if (username == actualUsername) && (password == actualPassword) {
			ret = true
		}
	}

	return ret
}

func createChannelOfUsers() <-chan map[string]interface{} {
	// Get all documents in the GopherChat Collection
	documents, err := db.Collection(myCollection).Documents(context.Background()).GetAll()
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		panic(err)
	}

	in := make(chan map[string]interface{})
	go func() {
		for _, doc := range documents {
			in <- doc.Data()
		}
		close(in)
	}()

	return in
}

func GetAllUsernames() <-chan string {
	in := createChannelOfUsers()
	out := make(chan string)

	go func() {
		for doc := range in {
			username := doc["username"]
			if username != nil {
				out <- username.(string)
			}
		}
		close(out)
	}()

	return out
}

func Close() {
	defer db.Close()
}
