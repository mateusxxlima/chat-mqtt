package mqttchat

import (
	"bufio"
	"fmt"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func chatRequests(scanner *bufio.Scanner, conn MQTT.Client) {
	for {
		screenName := " CHAT REQUESTS TO ME "
		printHeader(screenName)

		for i, item := range data.ChatRequestsToMe {
			if item.Status == "PENDING" {
				fmt.Printf("[%s] <%d> Accept; <%d> Refuse\n\n", item.From, (i*2)+1, (i*2)+2)
			} else {
				index := findUserByName(item.From)
				if index != -1 {
					user := data.AllUsers[index]
					online := "Offline"
					if user.Online {
						online = "Online"
					}
					fmt.Printf("[%s] (%s): %s\n\n", user.Name, online, user.Status)
				}
			}
		}
		scanner.Scan()
		text := scanner.Text()
		input, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println(err)
			return
		}
		if input == 0 {
			return
		}
		if input%2 == 0 {
			index := input/2 - 1
			data.ChatRequestsToMe[index].Status = "REFUSED"
			index2 := findUserByName(data.ChatRequestsToMe[index].From)
			if index2 != -1 {
				data.AllUsers[index2].Status = "REFUSED"
			}
		} else {
			index := (input - 1) / 2
			data.ChatRequestsToMe[index].Status = "ACCEPTED"
			createNewPrivateChat(data.ChatRequestsToMe[index].From, conn)
			index2 := findUserByName(data.ChatRequestsToMe[index].From)
			if index2 != -1 {
				data.AllUsers[index2].Status = "FRIENDS"
			}
		}
	}
}

func createNewPrivateChat(anotherUser string, conn MQTT.Client) {
	privateTopic := SELF_ID + "_" + anotherUser + "_" + strconv.Itoa(getRandomId())
	sysMsg := SysMessage{
		Action:       "PRIVATE_CHAT_ACCEPT",
		From:         SELF_ID,
		PrivateTopic: privateTopic,
	}
	pubInTopic(conn, "control/"+anotherUser, sysMsg)
	subInTopic(conn, privateTopic, userMsgListener)
	chat := Chat{
		Topic:          privateTopic,
		AnotherUser:    anotherUser,
		UserMessages:   []UserMessage{},
		UnreadMessages: 0,
		IsGroup:        false,
	}
	data.Chats = append(data.Chats, chat)
}
