package mqttchat

import (
	"bufio"
	"fmt"
	"strconv"
)

func home(scanner *bufio.Scanner) int {
	clearScreen()
	fmt.Print("\n=-=-=-=-=-=-=-=-=-=-  H O M E  =-=-=-=-=-=-=-=-=-=-\n\n")
	fmt.Println("<1> New private chat")
	fmt.Println("<2> Chats requests")
	fmt.Println("<3> Groups config")
	fmt.Println("<0> To close app")
	fmt.Println("\nYour chats")
	printChats()
	scanner.Scan()
	input, _ := strconv.ParseInt(scanner.Text(), 10, 64)
	return int(input)
}

func printChats() {
	for index, chat := range data.Chats {
		var lastMsg string
		if len(chat.UserMessages) > 0 {
			lastMsg = chat.UserMessages[len(chat.UserMessages)-1].Text
		}
		status := "Offline"
		if chat.Online {
			status = "Online"
		}
		fixedLen := len(fmt.Sprintf("  <%d> [%s] (%s): ", index, chat.AnotherUser, status)) +
			len(fmt.Sprintf(" (%d)", chat.UnreadMessages))
		if chat.IsGroup {
			fixedLen = len(fmt.Sprintf("  <%d> [%s]: ", index, chat.GroupName)) + len(fmt.Sprintf(" (%d)", chat.UnreadMessages))
		}

		previewSize := max(TOTAL_WIDTH-fixedLen, 0)
		if len(lastMsg) > previewSize {
			lastMsg = lastMsg[:previewSize-3] + "..."
		} else {
			lastMsg = fmt.Sprintf("%-*s", previewSize, lastMsg)
		}

		if chat.IsGroup {
			fmt.Printf("  <%d> [%s]: %s (%d)\n",
				index+MAGIC_NUMBER,
				chat.GroupName,
				lastMsg,
				chat.UnreadMessages,
			)
			continue
		}

		fmt.Printf("  <%d> [%s] (%s): %s (%d)\n",
			index+MAGIC_NUMBER,
			chat.AnotherUser,
			status,
			lastMsg,
			chat.UnreadMessages,
		)
	}
}
