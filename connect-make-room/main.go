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

		// Create the private welcome room
		createPrivateWelcomeRoom(s, "774195371903025168", "1109193687491162143")
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

func createPrivateWelcomeRoom(session *discordgo.Session, guildID string, memberID string) {
	// Create the category for the group
	category, err := session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     "Welcome to Beksoft",
		Type:     discordgo.ChannelTypeGuildCategory,
		Position: 1,
	})
	if err != nil {
		fmt.Println("Error creating category:", err)
		return
	}

	// Create the chat room under the category
	chatRoom, err := session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     "chat",
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: category.ID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   guildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
			{
				ID:    memberID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel,
			},
		},
	})
	if err != nil {
		fmt.Println("Error creating chat room:", err)
		return
	}

	// Create the voice room under the category
	voiceRoom, err := session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     "voice",
		Type:     discordgo.ChannelTypeGuildVoice,
		ParentID: category.ID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   guildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
			{
				ID:    memberID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel,
			},
		},
	})
	if err != nil {
		fmt.Println("Error creating voice room:", err)
		return
	}

	fmt.Println("Welcome group created successfully!")
	fmt.Printf("Category: %s\n", category.Name)
	fmt.Printf("Chat Room: %s\n", chatRoom.Name)
	fmt.Printf("Voice Room: %s\n", voiceRoom.Name)
}
