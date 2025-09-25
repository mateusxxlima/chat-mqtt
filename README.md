# MQTT Chat

**Author:** Mateus de Lima  

This project implements a **chat system in Go** using the **MQTT protocol** as the communication layer.  
It was developed **for academic purposes**, as part of an elective course, with the goal of learning and practicing the **MQTT protocol** and the **Mosquitto broker**.  
The course focused on understanding how MQTT is applied in **IoT scenarios**, and this project also served as an opportunity to explore and practice the **Go programming language**.  

## How it works

- Each message is transmitted through an **MQTT topic**.  
- Users can communicate in **private chats** or **group conversations**.  
- The application keeps all chats, groups, and requests in memory.  
- When the program exits, data can be persisted and restored on the next startup.  
- A text-based interface (terminal) is provided to navigate between chats, messages, and groups.  
- User presence (online/offline) is managed through periodic status updates published to the broker.  

## Prerequisites

1. Install [Go](https://go.dev/doc/install) (version 1.18 or higher recommended).  
2. Have access to a valid **MQTT broker** [Mosquitto](https://mosquitto.org/).  

## Running the application

Clone the repository and enter the project directory:


```bash
* Add a Mosquitto server address to the .env file

git clone https://github.com/mateusxxlima/chat-mqtt.git

cd chat-mqtt

go run . username
