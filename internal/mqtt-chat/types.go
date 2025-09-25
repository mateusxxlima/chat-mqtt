package mqttchat

type UserMessage struct {
	Read      bool
	From      string
	To        string
	Text      string
	IsGroup   bool
	GroupName string
}

type SysMessage struct {
	Action       string
	From         string
	PrivateTopic string
	Online       bool
	GroupName    string
}

type ChatRequest struct {
	From   string
	To     string
	Status string
}

type Chat struct {
	Topic          string
	IsGroup        bool
	GroupName      string
	AnotherUser    string
	UserMessages   []UserMessage
	UnreadMessages int
	Online         bool
}

type User struct {
	Name   string
	Status string
	Online bool
}

type Group struct {
	GroupName      string
	Owner          string
	Members        []User
	Topic          string
	Status         string
	RequestsToJoin []User
}

type AppData struct {
	ChatRequestsToMe []ChatRequest
	Chats            []Chat
	MyGroupsChats    []Group
	AllGroups        []Group
	AllUsers         []User
}
