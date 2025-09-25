package mqttchat

import (
	"bufio"
	"fmt"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func chats(chatIndex int, scanner *bufio.Scanner, conn MQTT.Client) {
	CHAT_NOW = &data.Chats[chatIndex]
	CHAT_NOW.UnreadMessages = 0
	if CHAT_NOW.IsGroup {
		printHeader(" " + CHAT_NOW.GroupName + " ")
	} else {
		printHeader(" " + CHAT_NOW.AnotherUser + " ")
	}

	for _, msg := range CHAT_NOW.UserMessages {
		if msg.From == SELF_ID {
			fmt.Printf("%s%s\n\n", strings.Repeat(" ", TOTAL_WIDTH/2), msg.Text)
		} else {
			fmt.Printf("%s: %s\n\n", msg.From, msg.Text)
		}
	}
	for {
		fmt.Printf("\n%*s", TOTAL_WIDTH/2, " ")
		scanner.Scan()
		text := scanner.Text()
		if text == "0" {
			CHAT_NOW = nil
			break
		}
		message := UserMessage{
			Read: false,
			From: SELF_ID,
			Text: text,
		}

		if CHAT_NOW.IsGroup {
			message.IsGroup = true
			message.GroupName = CHAT_NOW.GroupName
		} else {
			message.To = CHAT_NOW.AnotherUser
		}

		CHAT_NOW.UserMessages = append(CHAT_NOW.UserMessages, message)
		pubInTopic(conn, CHAT_NOW.Topic, message)
	}
}
