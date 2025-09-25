package mqttchat

import (
	"fmt"
	"strings"
	"time"
)

type Second int

func sleep(n Second) {
	time.Sleep(time.Duration(n) * time.Second)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func getRandomId() int {
	return int(time.Now().UnixMilli())
}

func printHeader(screenName string) {
	screenName = strings.ToUpper(screenName)
	clearScreen()
	padding := (TOTAL_WIDTH - len(screenName)) / 2
	rightPadding := TOTAL_WIDTH - len(screenName) - padding
	fmt.Printf("%s%s%s\n", strings.Repeat("=-", padding/2), screenName, strings.Repeat("-=", rightPadding/2))
	fmt.Print("\n<0> Back\n\n")
}

func findChatIndexByUser(user string) int {
	for i, chat := range data.Chats {
		if chat.AnotherUser == user {
			return i
		}
	}
	return -1
}

func findChatGroupByName(name string) int {
	for i, chat := range data.Chats {
		if chat.IsGroup && chat.GroupName == name {
			return i
		}
	}
	return -1
}

func findGroupInAllGroupsByName(name string) int {
	for i, group := range data.AllGroups {
		if group.GroupName == name {
			return i
		}
	}
	return -1
}

func findUserByName(name string) int {
	for i, user := range data.AllUsers {
		if user.Name == name {
			return i
		}
	}
	return -1
}
