package functions

import (
	"fmt"
	"strings"

	db "github.com/MahanthMohan/GopherChat/pkg/database"
	schema "github.com/MahanthMohan/GopherChat/pkg/schema"
)

var (
	usr schema.User
)

func LaunchApp() {
	var command string
	fmt.Println("<<>>-<<>>-<<>>-<<>>-<<>>-  Welcome to GopherChat, A Terminal Chat App <<>>-<<>>-<<>>-<<>>-<<>>-")
	fmt.Print("New (n/new) or Exisiting (e/existing) User: ")
	fmt.Scan(&command)
	for {
		switch command {
		case "n":
		case "new":
			RegisterNewUser()

		case "e":
		case "existing":
			LoginUser()

		default:
			fmt.Println("** Invalid Choice **")
			LaunchApp()
		}
	}
}

func validateUserCredentials(usr schema.User) {
	if len(usr.Username) == 0 || len(usr.Password) == 0 {
		fmt.Println("** Empty/Blank Password **")
		RegisterNewUser()
	} else if strings.HasSuffix(usr.Username, ".com") || strings.Contains(usr.Username, "@") {
		fmt.Println("** Username must not be an email, but a name **")
		RegisterNewUser()
	} else if len(usr.Password) < 6 {
		fmt.Println("** Password too short (Min 6 characters) **")
		RegisterNewUser()
	} else {
		fmt.Println("** User Validation Successful **")
		LoginUser()
	}
}

func RegisterNewUser() {
	fmt.Println("<<>>- Registration Screen -<<>>")
	fmt.Print("Username/Name: ")
	fmt.Scan(&usr.Username)
	fmt.Print("Password: ")
	fmt.Scan(&usr.Password)
	fmt.Print("Want to be a group member (true, false): ")
	fmt.Scanf("%t", &usr.IsGroupMember)
	usr.Messages = []schema.Message{}
	validateUserCredentials(usr)
}

func LoginUser() {
	fmt.Println("<<>>- Login Screen -<<>>")
	var uname, pw string
	fmt.Print("Username/Name: ")
	fmt.Scan(&uname)
	fmt.Print("Password: ")
	fmt.Scan(&pw)
	if uname == usr.Username && pw == usr.Password {
		fmt.Println("** Login Successful **")
		if !(usr.IsGroupMember) {
			viewAllMessages(usr.Username)
		} else {
			viewAllMessages("Group")
			viewAllMessages(usr.Username)
		}
	} else {
		fmt.Println("** Please Try Again **")
		// Repeat LoginUser() a maximum of 3 times before redirecting user back to registration screen
		for i := 0; i < 3; i++ {
			LoginUser()
		}
		RegisterNewUser()
	}
}

func viewAllMessages(username string) {
	fmt.Printf("<<>>- %s's Messages -<<>>\n", username)
	for _, msg := range db.GetAllMessages(username) {
		fmt.Println(msg.Author, ": ", msg.Content)
	}
}

func sendUserMessages() {
	fmt.Println("<<>>- Send Messages -<<>>")
	fmt.Println("--- List of Users ---")
	var groupMessages, dmMessages []schema.Message
	for _, user := range db.GetAllUsernames() {
		fmt.Println(user)
	}
	for {
		var userChoice, reciever string
		fmt.Print("Your Choice (msg/dm/(q/quit)): ")
		fmt.Scan(&userChoice)
		if userChoice == "msg" {
			var groupMessage schema.Message
			groupMessage.Author = usr.Username
			fmt.Print("Your Group Message: ")
			fmt.Scan(&groupMessage.Content)
			groupMessages = append(groupMessages, groupMessage)
			db.SendUserMessage("Group", groupMessages)
			viewAllMessages("Group")
		} else if userChoice == "dm" {
			var dmMessage schema.Message
			dmMessage.Author = usr.Username
			fmt.Println("--- List of Users ---")
			for _, user := range db.GetAllUsernames() {
				fmt.Println(user)
			}
			fmt.Print("Reciever: ")
			fmt.Scan(&reciever)
			fmt.Print("Your Direct Message: ")
			fmt.Scan(&dmMessage.Content)
			dmMessages = append(dmMessages, dmMessage)
			db.SendUserMessage(reciever, dmMessages)
			viewAllMessages(usr.Username)
		} else if userChoice == "q" || userChoice == "quit" {
			fmt.Println("<<>>- Hope to see you soon! -<<>>")
			break
		} else {
			fmt.Println("** Invalid Choice **")
			sendUserMessages()
			break
		}
	}
}
