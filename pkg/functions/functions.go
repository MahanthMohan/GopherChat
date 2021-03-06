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
	usr                       schema.User
	uname                     string
	pw                        string
	isGroupMember             bool
	groupMessages, dmMessages []string
	failCount                 int            = 0
	scanner                   *bufio.Scanner = bufio.NewScanner(os.Stdin)
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
		fmt.Println("** Registration Successful **")
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
	isGroupMember = db.GetMemberStatus(uname)
	if db.ValidateUserLoginCredentials(uname, pw) {
		color.Set(color.FgHiGreen, color.Bold)
		fmt.Println("\n** Login Successful **")
		if !(isGroupMember) {
			viewAllMessages(uname)
			sendUserMessages()
		} else {
			viewAllMessages("Group")
			viewAllMessages(uname)
			sendUserMessages()
		}
	} else {
		color.Set(color.FgHiRed, color.Bold)
		fmt.Println("\n** Please Try Again **")
		if failCount < 3 {
			failCount += 1
			LoginUser()
		} else {
			RegisterNewUser()
		}
	}
}

func viewAllMessages(username string) {
	if username == "Group" {
		color.Set(color.FgHiGreen, color.Bold)
	} else {
		color.Set(color.FgHiMagenta, color.Bold)
	}
	fmt.Printf("<<>>- %s's Messages -<<>>\n", username)
	messages := db.GetAllMessages(username)
	for msg := range messages {
		fmt.Println(msg)
		fmt.Println("-----------------------------------")
	}
}

func sendUserMessages() {
	color.Set(color.FgHiMagenta, color.Bold)
	fmt.Println("<<>>- Send Messages -<<>>")
	fmt.Println("--- List of Users ---")
	for user := range db.GetAllUsernames() {
		if uname != user {
			fmt.Println(user)
		}
	}

	for groupMsg := range db.GetAllMessages("Group") {
		groupMessages = append(groupMessages, groupMsg)
	}
	for dmMsg := range db.GetAllMessages(uname) {
		dmMessages = append(dmMessages, dmMsg)
	}

	var userChoice, reciever string
	for {
		fmt.Print("Your Choice (msg/dm/status/(q/quit)): ")
		fmt.Scan(&userChoice)
		if userChoice == "msg" {
			if isGroupMember {
				color.Set(color.FgHiGreen, color.Bold)
				var n int
				fmt.Print("# of words: ")
				fmt.Scan(&n)
				msgContent := make([]string, n)
				fmt.Print("Your Group Message: ")
				for i := 0; i < n; i++ {
					fmt.Scan(&msgContent[i])
				}
				groupMessage := fmt.Sprintf("%s: %s", uname, strings.Join(msgContent, " "))
				groupMessages = append(groupMessages, groupMessage)
				db.SendUserMessage("Group", groupMessages)
				viewAllMessages("Group")
			} else {
				color.Set(color.FgHiRed, color.Bold)
				fmt.Println("** You are not a group member **")
				sendUserMessages()
			}
		} else if userChoice == "dm" {
			color.Set(color.FgHiMagenta, color.Bold)
			fmt.Print("Reciever: ")
			fmt.Scan(&reciever)
			if reciever == uname {
				color.Set(color.FgHiRed, color.Bold)
				fmt.Println("** Cannot send a message to yourself! **")
				sendUserMessages()
			}
			var n int
			fmt.Print("# of words: ")
			fmt.Scan(&n)
			msgContent := make([]string, n)
			fmt.Print("Your Direct Message: ")
			for i := 0; i < n; i++ {
				fmt.Scan(&msgContent[i])
			}
			dmMessage := fmt.Sprintf("%s: %s", uname, strings.Join(msgContent, " "))
			dmMessages = append(dmMessages, dmMessage)
			db.SendUserMessage(reciever, dmMessages)
			viewAllMessages(uname)
		} else if userChoice == "q" || userChoice == "quit" {
			color.Set(color.FgHiRed, color.Bold)
			fmt.Println("<<>>- Sad to see you go -<<>>")
			db.Close()
			color.Unset()
			syscall.Exit(0)
		} else if userChoice == "status" {
			var groupMemberChoice string
			fmt.Print("Remain a group member (y/N): ")
			fmt.Scan(&groupMemberChoice)
			if groupMemberChoice == "y" {
				isGroupMember = true
				db.UpdateMemberStatus(uname, isGroupMember)
				sendUserMessages()
			} else {
				isGroupMember = false
				db.UpdateMemberStatus(uname, isGroupMember)
				sendUserMessages()
			}
		} else {
			color.Set(color.FgHiRed, color.Bold)
			fmt.Println("** Invalid Choice **")
			sendUserMessages()
		}
	}
}
