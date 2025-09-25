package mqttchat

import (
	"encoding/json"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var sysMsgListener MQTT.MessageHandler = func(conn MQTT.Client, msg MQTT.Message) {
	var message SysMessage
	err := json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		panic(err)
	}
	if message.From == SELF_ID {
		return
	}
	if message.Action == "REQUEST_PRIVATE_CHAT" {
		chatReq := ChatRequest{
			From:   message.From,
			To:     SELF_ID,
			Status: "PENDING",
		}
		data.ChatRequestsToMe = append(data.ChatRequestsToMe, chatReq)
	}
	if message.Action == "PRIVATE_CHAT_ACCEPT" {
		newChat := Chat{
			Topic:          message.PrivateTopic,
			AnotherUser:    message.From,
			UserMessages:   []UserMessage{},
			UnreadMessages: 0,
			IsGroup:        false,
		}
		subInTopic(conn, message.PrivateTopic, userMsgListener)
		data.Chats = append(data.Chats, newChat)
		index := findUserByName(message.From)
		data.AllUsers[index].Status = "FRIENDS"
	}
	if message.Action == "REQUEST_JOIN_GROUP" {
		request := User{
			Name:   message.From,
			Status: "PENDING",
		}
		for i, group := range data.MyGroupsChats {
			if group.GroupName == message.GroupName {
				data.MyGroupsChats[i].RequestsToJoin = append(data.MyGroupsChats[i].RequestsToJoin, request)
			}
		}
	}
	if message.Action == "GROUP_JOIN_ACCEPTED" {
		newChat := Chat{
			Topic:          message.PrivateTopic,
			IsGroup:        true,
			GroupName:      message.GroupName,
			UnreadMessages: 0,
		}
		data.Chats = append(data.Chats, newChat)
		subInTopic(conn, message.PrivateTopic, userMsgListener)
		index := findGroupInAllGroupsByName(message.GroupName)
		data.AllGroups[index].Status = "ACCEPTED"
	}
}

var userMsgListener MQTT.MessageHandler = func(conn MQTT.Client, msg MQTT.Message) {
	var message UserMessage
	err := json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		panic(err)
	}
	if message.From == SELF_ID {
		return
	}
	if len(message.GroupName) > 0 {
		index := findChatGroupByName(message.GroupName)
		if index == -1 {
			return
		}
		chat := &data.Chats[index]
		chat.UserMessages = append(chat.UserMessages, message)
		if CHAT_NOW != nil && chat.GroupName == CHAT_NOW.GroupName {
			fmt.Println()
			fmt.Printf("%s: %s\n\n", message.From, message.Text)
			fmt.Printf("%*s", TOTAL_WIDTH/2, " ")
			return
		}
		chat.UnreadMessages++
		return
	}
	index := findChatIndexByUser(message.From)
	if index == -1 {
		return
	}
	chat := &data.Chats[index]
	chat.UserMessages = append(chat.UserMessages, message)
	if CHAT_NOW != nil && chat.AnotherUser == CHAT_NOW.AnotherUser {
		fmt.Println()
		fmt.Printf("%s: %s\n\n", message.From, message.Text)
		fmt.Printf("%*s", TOTAL_WIDTH/2, " ")
		return
	}
	chat.UnreadMessages++
}

var groupsBroadcast MQTT.MessageHandler = func(_ MQTT.Client, msg MQTT.Message) {
	var group Group
	err := json.Unmarshal(msg.Payload(), &group)
	if err != nil {
		panic(err)
	}
	if group.Owner == SELF_ID {
		return
	}
	for i, g := range data.AllGroups {
		if g.GroupName == group.GroupName {
			data.AllGroups[i].Members = group.Members
			return
		}
	}
	group.Status = "NOT_REQUESTED"
	data.AllGroups = append(data.AllGroups, group)
}

var usersOnlineBroadcast MQTT.MessageHandler = func(_ MQTT.Client, msg MQTT.Message) {
	var message SysMessage
	err := json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		panic(err)
	}
	if message.From == SELF_ID {
		return
	}
	index := findUserByName(message.From)
	if index != -1 {
		data.AllUsers[index].Online = message.Online
		chatIndex := findChatIndexByUser(message.From)
		if chatIndex == -1 {
			return
		}
		data.Chats[chatIndex].Online = message.Online
		return
	}
	data.AllUsers = append(data.AllUsers, User{Name: message.From, Online: message.Online, Status: "NOT_REQUESTED"})
}
