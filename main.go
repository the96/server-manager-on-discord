package main

import (
	"fmt"
	"os"
	"net"
	"github.com/joho/godotenv"
	"github.com/bwmarrin/discordgo"
	"github.com/go-ping/ping"
)

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	u := m.Author
	fmt.Printf("%20s %20s(%20s) > %s\n", m.ChannelID, u.Username, u.ID, m.Content)
}

func pingServer() bool {
	address := os.Getenv("SERVER_IP_ADDRESS")
	pinger, err := ping.NewPinger(address)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	count := 1
	pinger.Count = count
	pinger.Timeout = 1_000_000_000 // nano seconds
	pinger.Run() // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats
	fmt.Println(stats)

	return (stats.PacketsRecv == count)
}

func wakeOnLan() bool {
	// https://monakaice88.hatenablog.com/entry/2019/10/29/070000
	targetMacAddress := os.Getenv("SERVER_MAC_ADDRESS")
	// localIpAddress   := net.ResolveUDPAddr("udp", ":0")
	// conn, err := net.Dial("udp", os.Getenv("SERVER_IP_ADDRESS")+":9")

	targetIpAddress, _ := net.ResolveUDPAddr("udp", os.Getenv("SERVER_IP_ADDRESS")+":9")
    conn, err := net.DialUDP("udp", nil, targetIpAddress)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer func() {
		_ = conn.Close()
	}()

	packetPrefix := make([]byte, 6)
	for i := range packetPrefix {
		packetPrefix[i] = 0xFF
	}

	hw, err := net.ParseMAC(targetMacAddress)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	sendPacket := make([]byte, 0)
	sendPacket = append(sendPacket, packetPrefix...)
	for i := 0; i < 16; i++ {
		sendPacket = append(sendPacket, hw...)
	}

	_, err = conn.Write(sendPacket)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}


func onInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// fmt.Printf("interaction: %#v\n", interaction.Interaction)
	// interactionType := interaction.Interaction.Type
	// interactionData := interaction.Interaction.Data
	// fmt.Printf("type:%s data:%s\n", interactionType, interactionData)
	
	interactionType := interaction.Interaction.Type
	fmt.Printf("type: %s\n", interactionType)

	data := interaction.Interaction.Data.(discordgo.ApplicationCommandInteractionData)
	options := data.Options
	if (len(options) != 1) {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "sorry, the command you sent is unknown command...",
			},
		})
	}
	
	name := options[0].Name
	response := ""
	switch name {
	case "status":
		status := "server is sleeped. :zzz:"
		if (pingServer()) {
			status = "server is running! :sunny:"
		}
		response = fmt.Sprintf("server status: %s\n", status)
	case "wake":
		wakeOnLan()
		response = "awaking server! please wait a moment."
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
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
