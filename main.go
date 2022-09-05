package main

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	clientId := os.Getenv("CLIENT_ID")
	u := m.Author
	fmt.Printf("%20s %20s(%20s) > %s\n", m.ChannelID, u.Username, u.ID, m.Content)
	if u.ID != clientId {
		fmt.Println(m.Type)
		sendMessage(s, m.ChannelID, u.Mention()+"なんか喋った!")
		sendReply(s, m.ChannelID, "test", m.Reference())
	}
}

func onInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	fmt.Printf("interaction: %#v\n", interaction.Interaction)
	interactionType := interaction.Type
	interactionData := interaction.Data

	fmt.Printf("type:%s data:%s\n", interactionType, interactionData)
	sendMessage(session, interaction.ChannelID, "receive command")
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "command received!",
		},
	})
}

func sendMessage(s *discordgo.Session, channelID string, msg string) {
	_, err := s.ChannelMessageSend(channelID, msg)
	fmt.Println(">>> " + msg)
	if err != nil {
		fmt.Println("Error sending message: ", err)
	}
}

func sendReply(s *discordgo.Session, channelID string, msg string, reference *discordgo.MessageReference) {
	_, err := s.ChannelMessageSendReply(channelID, msg, reference)
	if err != nil {
		fmt.Println("Error sending message: ", err)
	}
 }

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
        fmt.Println("cannot load")
    }

	fmt.Println("start")
	discord, err := discordgo.New(os.Getenv("TOKEN"))
	if err != nil {
        fmt.Println("cannot create")
    }
	
	fmt.Println("created")
	err = discord.Open()
	if err != nil {
        fmt.Println("cannot open")
		fmt.Println(err)
    }
	fmt.Println("open")

	discord.AddHandler(onMessageCreate)
	discord.AddHandler(onInteractionCreate)

	for{}
}
