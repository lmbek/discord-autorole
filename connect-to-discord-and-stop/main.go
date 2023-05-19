package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

// Flags
var (
	//GuildID        = flag.String("guild", "", "Test guild ID")
	//StageChannelID = flag.String("stage", "", "Test stage channel ID")
	BotToken = flag.String("token", "", "Bot token")
)

//func init() { flag.Parse() }

type Data struct {
	Token string `json:"token"`
}

func getTokenFromCredentials(filePath string) string {
	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON data
	var jsonData Data
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	// Access the token value
	return jsonData.Token
}

func login(token string) {
	session, _ := discordgo.New("Bot " + token)
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is ready")
	})

	err := session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer session.Close()
}

// To be correctly used, the bot needs to be in a guild.
// All actions must be done on a stage channel event
func main() {
	token := getTokenFromCredentials("credentials.json")
	login(token)
	/*
		// Create a new Stage instance on the previous channel
		si, err := s.StageInstanceCreate(&discordgo.StageInstanceParams{
			ChannelID:             *StageChannelID,
			Topic:                 "Amazing topic",
			PrivacyLevel:          discordgo.StageInstancePrivacyLevelGuildOnly,
			SendStartNotification: true,
		})
		if err != nil {
			log.Fatalf("Cannot create stage instance: %v", err)
		}
		log.Printf("Stage Instance %s has been successfully created", si.Topic)

		// Edit the stage instance with a new Topic
		si, err = s.StageInstanceEdit(*StageChannelID, &discordgo.StageInstanceParams{
			Topic: "New amazing topic",
		})
		if err != nil {
			log.Fatalf("Cannot edit stage instance: %v", err)
		}
		log.Printf("Stage Instance %s has been successfully edited", si.Topic)

		time.Sleep(5 * time.Second)
		if err = s.StageInstanceDelete(*StageChannelID); err != nil {
			log.Fatalf("Cannot delete stage instance: %v", err)
		}
		log.Printf("Stage Instance %s has been successfully deleted", si.Topic)
	*/
}
