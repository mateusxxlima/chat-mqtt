package mqttchat

import (
	"bufio"
	"fmt"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func newChat(scanner *bufio.Scanner, conn MQTT.Client) {
	for {
		screenName := " NEW CHAT "
		printHeader(screenName)

		fmt.Print("System users:\n\n")
		for i, user := range data.AllUsers {
			online := "Offline"
			if user.Online {
				online = "Online"
			}
			if user.Status == "NOT_REQUESTED" {
				fmt.Printf("[%s] (%s) <%d> Send chat request\n\n", user.Name, online, i+1)
			} else {
				fmt.Printf("[%s] (%s): %s\n\n", user.Name, online, user.Status)
			}
		}

		scanner.Scan()
		input, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		index := int(input)
		if index == 0 {
			return
		}
		if input < 0 || index > len(data.AllUsers) {
			continue
		}

		index--
		user := data.AllUsers[index]

		topic := "control/" + user.Name
		message := SysMessage{
			Action: "REQUEST_PRIVATE_CHAT",
			From:   SELF_ID,
		}
		pubInTopic(conn, topic, message)

		data.AllUsers[index].Status = "PENDING"
		fmt.Println("\nRequest sent to", user.Name)
		sleep(2)
	}
}
