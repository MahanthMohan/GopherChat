package functions

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	db "github.com/MahanthMohan/GopherChat/pkg/database"
	schema "github.com/MahanthMohan/GopherChat/pkg/schema"
	"golang.org/x/term"
)

var (
	usr     schema.User
	scanner = bufio.NewScanner(os.Stdin)
	uname   string
	pw      string
)

func LaunchApp() {
	var command string
	fmt.Println("<<>>- <<>>- <<>>-  Welcome to GopherChat, A Terminal Chat App -<<>> -<<>> -<<>>")
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
	var groupMemberChoice string
	fmt.Print("\nRemain a group member (y/N): ")
	fmt.Scan(&groupMemberChoice)
	if groupMemberChoice == "y" {
		usr.IsGroupMember = true
	} else {
		usr.IsGroupMember = false
	}
	if db.ValidateUserLoginCredentials(uname, pw) {
		fmt.Println("** Login Successful **")
		if !(usr.IsGroupMember) {
			viewAllMessages(uname)
			sendUserMessages()
		} else {
			viewAllMessages("Group")
			viewAllMessages(uname)
			sendUserMessages()
		}
	} else {
		fmt.Println("\n** Please Try Again **")
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
	if len(messages) == 0 {
		fmt.Println("No Messages Yet!")
	} else {
		for _, msg := range messages {
			if msg != nil {
				fmt.Println(msg.(string))
			} else {
				fmt.Println("No Messages Yet!")
				break
			}
		}
	}
}

func sendUserMessages() {
	fmt.Println("<<>>- Send Messages -<<>>")
	fmt.Println("--- List of Users ---")
	var groupMessages, dmMessages []string
	for user := range db.GetAllUsernames() {
		fmt.Println(user)
	}
	var userChoice, reciever string
	for {
		fmt.Print("Your Choice (msg/dm/(q/quit)): ")
		fmt.Scan(&userChoice)
		if userChoice == "msg" {
			var groupMessage string
			fmt.Print("Your Group Message: ")
			scanner.Scan()
			groupMessage = scanner.Text()
			groupMessage = fmt.Sprintf("%s: %s", uname, groupMessage)
			groupMessages = append(groupMessages, groupMessage)
			db.SendUserMessage("Group", groupMessages)
			fmt.Println("")
			viewAllMessages("Group")
		} else if userChoice == "dm" {
			var dmMessage string
			fmt.Print("Reciever: ")
			fmt.Scan(&reciever)
			fmt.Print("Your Direct Message: ")
			scanner.Scan()
			dmMessage = scanner.Text()
			dmMessage = fmt.Sprintf("%s: %s", uname, dmMessage)
			dmMessages = append(dmMessages, dmMessage)
			db.SendUserMessage(reciever, dmMessages)
			fmt.Println("")
			viewAllMessages(uname)
		} else if userChoice == "q" || userChoice == "quit" {
			fmt.Println("<<>>- Sad to see you go -<<>>")
			syscall.Exit(0)
		} else {
			fmt.Println("** Invalid Choice **")
			sendUserMessages()
		}
	}
}
