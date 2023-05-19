package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

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

func main() {
	// Discord bot token
	token := getTokenFromCredentials("credentials.json")

	// Create a new Discord session
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Register the handler for the Ready event
	session.AddHandler(onReady)

	// Open a websocket connection to Discord
	err = session.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	// Wait for a CTRL-C signal to gracefully close the bot
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close the Discord session
	session.Close()
}

// Handler for the Ready event
func onReady(session *discordgo.Session, event *discordgo.Ready) {
	// Get the guild ID from the first available guild
	guilds := session.State.Guilds
	if len(guilds) > 0 {
		guildID := guilds[0].ID
		assignRoleToUser(session, guildID, "LarsTest", "Customer")
	} else {
		fmt.Println("Error: No guilds available")
	}
}

func assignRoleToUser(session *discordgo.Session, guildID, username, roleName string) {
	// Find the member by username
	members, err := session.GuildMembers(guildID, "", 1000)
	if err != nil {
		fmt.Println("Error retrieving guild members:", err)
		return
	}

	var member *discordgo.Member
	for _, m := range members {
		if m.User.Username == username {
			member = m
			break
		}
	}

	// Find the role by name
	roles, err := session.GuildRoles(guildID)
	if err != nil {
		fmt.Println("Error retrieving guild roles:", err)
		return
	}

	var roleID string
	for _, role := range roles {
		if role.Name == roleName {
			roleID = role.ID
			break
		}
	}

	// If the member and role are found, assign the role to the member
	if member != nil && roleID != "" {
		err = session.GuildMemberRoleAdd(guildID, member.User.ID, roleID)
		if err != nil {
			fmt.Println("Error assigning role to user:", err)
			return
		}
		fmt.Println("Assigned role", roleName, "to user", username)
	} else {
		fmt.Println("Error: Member or role not found")
	}
}
