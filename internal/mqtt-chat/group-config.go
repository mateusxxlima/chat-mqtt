package mqttchat

import (
	"bufio"
	"fmt"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func groupConfig(conn MQTT.Client, scanner *bufio.Scanner) {
	for {
		clearScreen()
		fmt.Print("\n=-=-=-=-=-=-=-=-=-= GROUPS CONFIG =-=-=-=-=-=-=-=-=-\n\n")
		fmt.Println("<1> My groups")
		fmt.Printf("<2> Find groups (%d)\n", len(data.AllGroups))
		fmt.Println("<0> Back")

		scanner.Scan()
		input, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		option := int(input)

		switch option {
		case 1:
			myGroups(conn, scanner)
		case 2:
			findGroups(conn, scanner)
		case 0:
			return
		}
	}
}

func myGroups(conn MQTT.Client, scanner *bufio.Scanner) {
	for {
		screenName := " MY GROUPS "
		printHeader(screenName)

		fmt.Println("<1> Create new Group")
		fmt.Print("\nYour groups:\n\n")
		for i, group := range data.MyGroupsChats {
			fmt.Printf("    [%s] <%d> See group\n\n", group.GroupName, i+2)
		}
		scanner.Scan()
		input, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println(err)
			return
		}
		if input == 0 {
			return
		}
		if input == 1 {
			createGroup(conn, scanner)
			continue
		}
		indexGroup := input - 2
		if indexGroup > len(data.MyGroupsChats) || indexGroup < 0 {
			continue
		}
		seeGroup(indexGroup, conn, scanner)
	}
}

func seeGroup(indexGroup int, conn MQTT.Client, scanner *bufio.Scanner) {
	for {
		group := &data.MyGroupsChats[indexGroup]
		screenName := " " + group.GroupName + " "
		printHeader(screenName)
		for j, member := range group.RequestsToJoin {
			if member.Status != "PENDING" {
				fmt.Printf("[%s] status <%s> \n", member.Name, member.Status)
			}
			if member.Status == "PENDING" {
				fmt.Printf("[%s] wanna join in the group: <%d> Accept; <%d> Refuse\n", member.Name, (j*2)+1, (j*2)+2)
			}
			for _, member := range group.Members {
				fmt.Printf("  - %s\n", member.Name)
			}
			fmt.Println()
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
		if input < 0 || input > len(data.MyGroupsChats) {
			continue
		}
		var index int
		var status string
		if input%2 == 0 {
			index = input/2 - 1
			status = "REFUSED"
		} else {
			index = (input - 1) / 2
			status = "ACCEPTED"
		}

		group.RequestsToJoin[index].Status = status
		if status == "REFUSED" {
			return
		}
		sysMsg := SysMessage{
			Action:       "GROUP_JOIN_ACCEPTED",
			PrivateTopic: group.Topic,
			GroupName:    group.GroupName,
		}

		memberName := group.RequestsToJoin[index].Name
		pubInTopic(conn, "control/"+memberName, sysMsg)
		member := User{Name: memberName, Status: "MEMBER"}
		group.Members = append(group.Members, member)
	}
}

func createGroup(conn MQTT.Client, scanner *bufio.Scanner) {
	fmt.Print("Group Name: ")

	scanner.Scan()
	groupName := scanner.Text()
	if groupName == "0" {
		return
	}
	topic := fmt.Sprintf("%s_%d", groupName, getRandomId())

	newGroup := Group{
		GroupName: groupName,
		Owner:     SELF_ID,
		Topic:     topic,
		Members:   []User{{Name: SELF_ID}},
	}
	newChat := Chat{
		Topic:          topic,
		IsGroup:        true,
		UnreadMessages: 0,
		GroupName:      groupName,
	}

	data.MyGroupsChats = append(data.MyGroupsChats, newGroup)
	data.Chats = append(data.Chats, newChat)
	subInTopic(conn, topic, userMsgListener)

	fmt.Printf("\nGroup %v created!\n", groupName)
	sleep(3)
}

func findGroups(conn MQTT.Client, scanner *bufio.Scanner) {
	for {
		screenName := " ALL GROUPS "
		printHeader(screenName)

		if len(data.AllGroups) == 0 {
			fmt.Print("\n\nNo groups yet :)")
			sleep(3)
			return
		}

		for i, group := range data.AllGroups {
			if group.Status == "NOT_REQUESTED" {
				fmt.Printf("[%s] (Owner: %s) <%d> Ask to join\n", group.GroupName, group.Owner, i+1)
			} else {
				fmt.Printf("[%s] (Owner: %s) Status: %s\n", group.GroupName, group.Owner, group.Status)
			}
			for _, member := range group.Members {
				fmt.Printf("  - %s\n", member.Name)
			}
			fmt.Println()
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
		if input < 0 || input > len(data.AllGroups) {
			continue
		}

		group := &data.AllGroups[input-1]
		message := SysMessage{
			Action:    "REQUEST_JOIN_GROUP",
			From:      SELF_ID,
			GroupName: group.GroupName,
		}
		pubInTopic(conn, "control/"+group.Owner, message)
		group.Status = "REQUESTED"
	}
}
