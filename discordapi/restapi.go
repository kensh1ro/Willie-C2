package discordapi

import (
	"bytes"
	"encoding/json"
)

/*
func get_guilds() []Guild {
	response := Route("GET", "/users/@me/guilds")
	guilds := []Guild{}
	json.Unmarshal(response, &guilds)
	return guilds
}

func get_guild_channels(guild_id string) []Channel {
	response := Route("GET", "/guilds/"+guild_id+"/channels")
	channels := []Channel{}
	json.Unmarshal(response, &channels)
	return channels
}
*/
func GetChannel(channel_id string) Channel {
	response := Route("GET", "/channels/"+channel_id, nil, nil)
	channel := Channel{}
	json.Unmarshal(response, &channel)
	return channel
}

func GetMessage(channel_id string) Message {
	channel := GetChannel(channel_id)
	response := Route("GET", "/channels/"+channel_id+"/messages/"+channel.LMI, nil, nil)
	message := Message{}
	json.Unmarshal(response, &message)
	return message
}

func GetGateway() {
	response := Route("GET", "/gateway/bot", nil, nil)
	message := Message{}
	json.Unmarshal(response, &message)
}

func SendMessage(channel_id string, content string) {
	m := MessageSender{}
	m.Content = content
	form, _ := json.Marshal(m)
	Route("POST", "/channels/"+channel_id+"/messages", bytes.NewReader(form), nil)
}

func SendFile(channel_id string, filename string, data []byte) {
	a := Attachment{}
	a.Filename = filename
	if filename == "image" {
		a.ContentType = "image/png"
	} else {
		a.ContentType = "application/octet-stream"
	}
	Route("POST", "/channels/"+channel_id+"/messages", bytes.NewReader(data), &a)
}
