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

type InviteLink struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	RoleID string `json:"roleID"`
}

type InviteRoleMap struct {
	Links []InviteLink `json:"links"`
}

type MemberList struct {
	Members        map[string]bool
	InviteCounters map[string]int
	InviteRole     map[string]string
}

var inviteRoleMap InviteRoleMap

var guildID = "774195371903025168" // Set your server id here

func main() {
	//token, err := getTokenFromCredentials("./credentials.json") // test server
	token, err := getTokenFromCredentials("/var/www/discord-autorole-bot/credentials.json") // production
	if err != nil {
		log.Fatal("Error retrieving bot token:", err)
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	// (Only 1 "hardcoded" server id is supported for now) Please note this will only work for one server, the script needs to be improved and extended upon to add more guilds
	//inviteRoleMap, err = getInviteRoleMapFromJSON("./invite-roles.json") // test server
	inviteRoleMap, err = getInviteRoleMapFromJSON("/var/www/discord-autorole-bot/invite-roles.json") //production
	if err != nil {
		log.Fatal("Error reading invited-roles.json:", err)
	}

	// (Only 1 "hardcoded" server id is supported for now) Please note this will only work for one server, the script needs to be improved and extended upon to add more guilds
	memberList := &MemberList{
		Members:        make(map[string]bool),
		InviteCounters: make(map[string]int),
		InviteRole:     make(map[string]string),
	}

	session.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		for _, guild := range event.Guilds {
			// ensure it only runs on the selected "hardcoded" server ID
			if guild.ID == guildID {
				onReady(s, event, memberList)
			}
		}
	})

	session.AddHandler(func(s *discordgo.Session, event *discordgo.MessageCreate) {
		// ensure it only runs on the selected "hardcoded" server ID
		if event.GuildID == guildID {
			onEvent(s, event, memberList)
		}
	})

	err = session.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()
}

func getTokenFromCredentials(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var jsonData Data
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return "", err
	}

	return jsonData.Token, nil
}

func getInviteRoleMapFromJSON(filePath string) (InviteRoleMap, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return InviteRoleMap{}, err
	}

	var inviteRoleMap InviteRoleMap
	err = json.Unmarshal(data, &inviteRoleMap)
	if err != nil {
		return InviteRoleMap{}, err
	}

	return inviteRoleMap, nil
}

func getInviteRole(inviteRoleMap InviteRoleMap, value string) string {
	for _, link := range inviteRoleMap.Links {
		if link.Value == value {
			return link.RoleID
		}
	}
	return ""
}

func onReady(session *discordgo.Session, event *discordgo.Ready, memberList *MemberList) {
	fmt.Println("Bot is ready!")

	guilds := session.State.Guilds
	if len(guilds) == 0 {
		fmt.Println("No guilds found")
		return
	}

	guildID := guilds[0].ID

	members, err := session.GuildMembers(guildID, "", 1000)
	if err != nil {
		fmt.Println("Error fetching guild members:", err)
		return
	}

	for _, member := range members {
		memberList.Members[member.User.ID] = true
	}

	invites, err := session.GuildInvites(guildID)
	if err != nil {
		fmt.Println("Error fetching guild invites:", err)
		return
	}

	for _, invite := range invites {
		//fmt.Printf("Invite code: %v \n", invite.Code)
		memberList.InviteCounters[invite.Code] = invite.Uses
		memberList.InviteRole[invite.Code] = "" // Replace "roleID" with the corresponding role ID or name
	}

	fmt.Println("Existing members:")
	fmt.Println("----------------------")
	for member := range memberList.Members {
		fmt.Println(member)
	}
	fmt.Println("All member id's listed")
	fmt.Println("----------------------")
}

func onEvent(session *discordgo.Session, event interface{}, memberList *MemberList) {
	switch ev := event.(type) {
	case *discordgo.MessageCreate:
		if ev.Author.Bot {
			return
		}

		//fmt.Println("new message...")
		authorID := ev.Author.ID

		if memberList.Members[authorID] {
			return
		}

		memberList.Members[authorID] = true

		eventGuildID := ev.GuildID
		invites, err := session.GuildInvites(eventGuildID)
		if err != nil {
			fmt.Println("Error fetching guild invites:", err)
			return
		}

		var usedInvite *discordgo.Invite
		for _, invite := range invites {
			if invite.Uses > memberList.InviteCounters[invite.Code] {
				usedInvite = invite
				break
			}
		}

		if usedInvite == nil {
			fmt.Println("No invite found for the user:", ev.Author.ID)
			return
		}

		fmt.Println("New Member:")
		inviteCode := usedInvite.Code
		fmt.Print("invite code: ")
		fmt.Println(inviteCode)
		roleID := getInviteRole(inviteRoleMap, inviteCode)
		if roleID == "" {
			fmt.Println("No role assigned for invite code:", inviteCode)
			return
		}

		if roleID == "" {
			fmt.Println("No role assigned for invite link:", inviteCode)
			return
		}

		assignRoleToUser(session, eventGuildID, ev.Author.ID, roleID)

		err = session.GuildMemberRoleAdd(eventGuildID, ev.Author.ID, roleID)
		if err != nil {
			fmt.Println("Error assigning role to user:", err)
			return
		}

		fmt.Printf("Assigned role '%s' to user %s\n", roleID, ev.Author.ID)

		createPrivateWelcomeRoom(session, eventGuildID, ev.Author.ID) // Call the function here

	}
}

func assignRoleToUser(session *discordgo.Session, guildID string, userID string, roleID string) {
	// Find the member by ID
	members, err := session.GuildMembers(guildID, "", 1000)
	if err != nil {
		fmt.Println("Error retrieving guild members:", err)
		return
	}

	var member *discordgo.Member
	for _, m := range members {
		if m.User.ID == userID {
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

	for _, role := range roles {
		if role.ID == roleID {
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
		fmt.Println("Assigned role", roleID, "to user", userID)
	} else {
		fmt.Println("Error: Member or role not found")
	}
}

func createPrivateWelcomeRoom(session *discordgo.Session, guildID string, memberID string) {
	// Create the category for the group
	category, err := session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     "Welcome! Tell us about you",
		Type:     discordgo.ChannelTypeGuildCategory,
		Position: 1,
	})
	if err != nil {
		fmt.Println("Error creating category:", err)
		return
	}

	// Create the chat room under the category
	chatRoom, err := session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     "Chat with Beksoft ApS",
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
			{
				ID:    session.State.User.ID, // Add the bot's user ID here
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel,
			},
		},
	})
	if err != nil {
		fmt.Println("Error creating chat room:", err)
		return
	}

	// Send the welcome message
	welcomeMessage := "Welcome, let us know more about you! / Velkommen, lad os lære hinanden at kende"
	_, err = session.ChannelMessageSend(chatRoom.ID, welcomeMessage)
	if err != nil {
		fmt.Println("Error sending welcome message:", err)
	}

	// Send the welcome message
	welcomeMessage2 := "We answer as soon as possible / Vi svarer hurtigst muligt"
	_, err = session.ChannelMessageSend(chatRoom.ID, welcomeMessage2)
	if err != nil {
		fmt.Println("Error sending welcome message:", err)
	}

	// Create the voice room under the category
	//voiceRoom, err := session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
	//	Name:     "voice",
	//	Type:     discordgo.ChannelTypeGuildVoice,
	//	ParentID: category.ID,
	//	PermissionOverwrites: []*discordgo.PermissionOverwrite{
	//		{
	//			ID:   guildID,
	//			Type: discordgo.PermissionOverwriteTypeRole,
	//			Deny: discordgo.PermissionViewChannel,
	//		},
	//		{
	//			ID:    memberID,
	//			Type:  discordgo.PermissionOverwriteTypeMember,
	//			Allow: discordgo.PermissionViewChannel,
	//		},
	//	},
	//})
	//if err != nil {
	//	fmt.Println("Error creating voice room:", err)
	//	return
	//}

	fmt.Println("Welcome group created successfully!")
	fmt.Printf("Category: %s\n", category.Name)
	fmt.Printf("Chat Room: %s\n", chatRoom.Name)
	//fmt.Printf("Voice Room: %s\n", voiceRoom.Name)
	fmt.Println("----------------------")
}
