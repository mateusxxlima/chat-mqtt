package mqttchat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/mateusxxlima/chat-mqtt/internal/config"
)

const (
	TOTAL_WIDTH        = 60
	MAGIC_NUMBER       = 4
	TOPIC_USERS_ONLINE = "control/users"
	TOPIC_GROUPS       = "control/groups"
	APP_DATA_FOLDER    = "app-data"
)

var SELF_ID string
var CHAT_NOW *Chat
var data = AppData{}

func Start() {
	var err error
	err = config.StartConfig()
	FatalError(err)
	SELF_ID = os.Args[1]

	fileStorageData := fmt.Sprintf(APP_DATA_FOLDER+"/%s.json", SELF_ID)
	conn, scanner := bootConfig()
	err = loadAppData(fileStorageData, conn)
	FatalError(err)

	go startBroadcasting(conn)
	defer appClose(fileStorageData, conn)

	var homeAction int
	for {
		homeAction = home(scanner)
		if homeAction == 1 {
			newChat(scanner, conn)
		} else if homeAction == 2 {
			chatRequests(scanner, conn)
		} else if homeAction == 3 {
			groupConfig(conn, scanner)
		} else if homeAction >= MAGIC_NUMBER && homeAction <= len(data.Chats)+MAGIC_NUMBER {
			index := homeAction - MAGIC_NUMBER
			chats(index, scanner, conn)
		} else if homeAction >= MAGIC_NUMBER && homeAction <= len(data.Chats)+MAGIC_NUMBER {
			index := homeAction - MAGIC_NUMBER
			chats(index, scanner, conn)
		} else if homeAction == 0 {
			break
		}
	}
}

func FatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func bootConfig() (conn MQTT.Client, scanner *bufio.Scanner) {
	scanner = bufio.NewScanner(os.Stdin)
	conn = NewMQTTClient()
	subInTopic(conn, TOPIC_USERS_ONLINE, usersOnlineBroadcast)
	myTopicControl := fmt.Sprintf("control/%s", SELF_ID)
	subInTopic(conn, myTopicControl, nil)
	subInTopic(conn, TOPIC_GROUPS, groupsBroadcast)
	return
}

func appClose(fileStorageData string, conn MQTT.Client) {
	fmt.Println("Closing app ...")
	pubInTopic(conn, TOPIC_USERS_ONLINE, SysMessage{Online: false, From: SELF_ID})
	conn.Disconnect(1000)
	if err := saveAppData(fileStorageData, &data); err != nil {
		fmt.Println("Error saving app data:", err)
	}
	os.Exit(0)
}

func startBroadcasting(conn MQTT.Client) {
	for {
		pubInTopic(conn, TOPIC_USERS_ONLINE, SysMessage{Online: true, From: SELF_ID})
		for _, g := range data.MyGroupsChats {
			group := Group{
				GroupName: g.GroupName,
				Owner:     g.Owner,
				Members:   g.Members,
			}
			pubInTopic(conn, TOPIC_GROUPS, group)
		}
		sleep(10)
	}
}

func saveAppData(fileName string, data *AppData) error {
	if err := os.MkdirAll(APP_DATA_FOLDER, os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func loadAppData(fileName string, conn MQTT.Client) error {
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	reSubChatsTopics(conn)
	return nil
}

func reSubChatsTopics(conn MQTT.Client) {
	for _, chats := range data.Chats {
		subInTopic(conn, chats.Topic, userMsgListener)
	}
}
