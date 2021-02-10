package functions

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	db "github.com/MahanthMohan/GopherChat/pkg/database"
	schema "github.com/MahanthMohan/GopherChat/pkg/schema"
	"github.com/fatih/color"
	"golang.org/x/term"
)

var (
	usr   schema.User
	uname string
	pw    string
)

func LaunchApp() {
	var command string
	color.Set(color.FgHiYellow, color.Bold)
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
	color.Set(color.FgHiRed, color.Bold)
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
		color.Set(color.FgHiGreen, color.Bold)
		fmt.Println("** User Validation Successful **")
		db.CreateUserDocument(usr)
		LoginUser()
	}
}

func RegisterNewUser() {
	color.Set(color.FgHiCyan, color.Bold)
	fmt.Println("<<>>- Registration Screen -<<>>")
	fmt.Print("Username/Name (No Spaces, Single Word): ")
	fmt.Scan(&usr.Username)
	fmt.Print("Password: ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		color.Set(color.BgHiRed, color.Bold)
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
	color.Set(color.FgHiBlue, color.Bold)
	fmt.Println("<<>>- Login Screen -<<>>")
	fmt.Print("Username/Name: ")
	fmt.Scan(&uname)
	fmt.Print("Password: ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		color.Set(color.BgHiRed, color.Bold)
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
		color.Set(color.FgHiGreen, color.Bold)
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
		color.Set(color.FgHiRed, color.Bold)
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
	if username == "Group" {
		color.Set(color.FgHiGreen, color.Bold)
	} else {
		color.Set(color.FgHiMagenta, color.Bold)
	}
	fmt.Printf("\n<<>>- %s's Messages -<<>>\n", username)
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
	color.Set(color.FgHiMagenta, color.Bold)
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
			color.Set(color.FgHiGreen, color.Bold)
			var groupMessage string
			fmt.Print("Your Group Message: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			groupMessage = scanner.Text()
			groupMessage = fmt.Sprintf("%s: %s", uname, groupMessage)
			groupMessages = append(groupMessages, groupMessage)
			db.SendUserMessage("Group", groupMessages)
			viewAllMessages("Group")
		} else if userChoice == "dm" {
			color.Set(color.FgHiMagenta, color.Bold)
			var dmMessage string
			fmt.Print("Reciever: ")
			fmt.Scan(&reciever)
			fmt.Print("Your Direct Message: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			dmMessage = scanner.Text()
			dmMessage = fmt.Sprintf("%s: %s", uname, dmMessage)
			dmMessages = append(dmMessages, dmMessage)
			db.SendUserMessage(reciever, dmMessages)
			viewAllMessages(uname)
		} else if userChoice == "q" || userChoice == "quit" {
			color.Set(color.FgHiRed, color.Bold)
			fmt.Println("<<>>- Sad to see you go -<<>>")
			db.Close()
			color.Unset()
			syscall.Exit(0)
		} else {
			color.Set(color.FgHiRed, color.Bold)
			fmt.Println("** Invalid Choice **")
			sendUserMessages()
		}
	}
}
