package functions

import (
	"fmt"
	"strings"
	"syscall"

	db "github.com/MahanthMohan/GopherChat/pkg/database"
	schema "github.com/MahanthMohan/GopherChat/pkg/schema"
	"golang.org/x/term"
)

var (
	usr   schema.User
	uname string
	pw    string
)

func LaunchApp() {
	var command string
	fmt.Println("<<>>-<<>>-<<>>-<<>>-<<>>-  Welcome to GopherChat, A Terminal Chat App <<>>-<<>>-<<>>-<<>>-<<>>-")
	fmt.Print("New (n/new) or Existing (e/existing) User: ")
	fmt.Scan(&command)
	if command == "n" || command == "new" {
		RegisterNewUser()
	} else if command == "e" || command == "existing" {
		LoginUser()
	} else {
		fmt.Println("** Invalid Choice **")
		LaunchApp()
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
		db.CreateUserDocument(usr)
		LoginUser()
	}
}

func RegisterNewUser() {
	fmt.Println("<<>>- Registration Screen -<<>>")
	fmt.Print("Username/Name (No Spaces, Single Word): ")
	fmt.Scan(&usr.Username)
	fmt.Print("Password: ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	usr.Password = string(bytepw)
	var groupMemberChoice string
	fmt.Print("\nWant to be a group member (y/N): ")
	fmt.Scan(&groupMemberChoice)
	if groupMemberChoice == "y" {
		usr.IsGroupMember = true
	} else if groupMemberChoice == "N" {
		usr.IsGroupMember = false
	}
	usr.Messages = []string{}
	validateUserCredentials(usr)
}

func LoginUser() {
	fmt.Println("<<>>- Login Screen -<<>>")
	fmt.Print("Username/Name: ")
	fmt.Scan(&uname)
	fmt.Print("Password: ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	pw = string(bytepw)
	if db.ValidateUserLoginCredentials(uname, pw) {
		fmt.Println("\n** Login Successful **")
		if !(usr.IsGroupMember) {
			viewAllMessages(uname)
		} else {
			viewAllMessages("Group")
			viewAllMessages(uname)
		}
	} else {
		fmt.Println("** Please Try Again **")
		// Repeat LoginUser() a maximum of 3 times before redirecting user back to registration screen
		for i := 0; i < 3; i++ {
			LoginUser()
		}
		fmt.Println("** Login limit: 3 tries **")
		RegisterNewUser()
	}
}

func viewAllMessages(username string) {
	fmt.Printf("<<>>- %s's Messages -<<>>\n", username)
	messages := db.GetAllMessages(username)
	for _, msg := range messages {
		fmt.Println(msg.(string))
	}
}

func sendUserMessages() {
	fmt.Println("<<>>- Send Messages -<<>>")
	fmt.Println("--- List of Users ---")
	var groupMessages, dmMessages []string
	for _, user := range db.GetAllUsernames() {
		fmt.Println(user)
	}
	var userChoice, reciever string
	for !(userChoice == "q" || userChoice == "quit") {
		fmt.Print("Your Choice (msg/dm/(q/quit)): ")
		fmt.Scan(&userChoice)
		if userChoice == "msg" {
			var groupMessage string
			fmt.Print("Your Group Message: ")
			fmt.Scan(groupMessage)
			groupMessage = fmt.Sprintf("%s: %s", uname, groupMessage)
			groupMessages = append(groupMessages, groupMessage)
			db.SendUserMessage("Group", groupMessages)
			viewAllMessages("Group")
			sendUserMessages()
		} else if userChoice == "dm" {
			var dmMessage string
			fmt.Println("--- List of Users ---")
			for _, user := range db.GetAllUsernames() {
				fmt.Println(user)
			}
			fmt.Print("Reciever: ")
			fmt.Scan(&reciever)
			fmt.Print("Your Direct Message: ")
			fmt.Scan(&dmMessage)
			dmMessage = fmt.Sprintf("%s: %s", uname, dmMessage)
			dmMessages = append(dmMessages, dmMessage)
			db.SendUserMessage(reciever, dmMessages)
			viewAllMessages(uname)
			sendUserMessages()
		} else {
			fmt.Println("** Invalid Choice **")
			sendUserMessages()
			break
		}
	}
}
