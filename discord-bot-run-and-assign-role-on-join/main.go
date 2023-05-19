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
	Token      string            `json:"token"`
	InviteRole map[string]string `json:"inviteRole"` // Map to store invite links and their corresponding roles
}

type MemberList struct {
	Members        map[string]bool
	InviteCounters map[string]int
	InviteRole     map[string]string
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

func getInviteRole(inviteLink string, inviteRoleMap map[string]string) string {
	if role, ok := inviteRoleMap[inviteLink]; ok {
		return role
	}
	return "" // Replace with the default role ID or name
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
		memberList.Members[member.User.Username] = true
	}

	invites, err := session.GuildInvites(guildID)
	if err != nil {
		fmt.Println("Error fetching guild invites:", err)
		return
	}

	for _, invite := range invites {
		memberList.InviteCounters[invite.Code] = invite.Uses
		memberList.InviteRole[invite.Code] = "1108944121630044172" // Replace "roleID" with the corresponding role ID or name
	}

	fmt.Println("Existing members:")
	for member := range memberList.Members {
		fmt.Println(member)
	}
}

func onEvent(session *discordgo.Session, event interface{}, memberList *MemberList) {
	switch ev := event.(type) {
	case *discordgo.MessageCreate:
		if ev.Author.Bot {
			return
		}

		fmt.Println("new message...")
		authorID := ev.Author.ID

		if memberList.Members[authorID] {
			return
		}

		memberList.Members[authorID] = true

		guildID := ev.GuildID
		invites, err := session.GuildInvites(guildID)
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
			fmt.Println("No invite found for the user:", ev.Author.Username)
			return
		}

		inviteCode := usedInvite.Code
		fmt.Print("invite code:")
		fmt.Println(inviteCode)
		role := getInviteRole(inviteCode, memberList.InviteRole)
		if role == "" {
			fmt.Println("No role assigned for invite link:", inviteCode)
			return
		}
		guilds := session.State.Guilds
		if len(guilds) > 0 {
			guildID := guilds[0].ID
			assignRoleToUser(session, guildID, ev.Author.Username, "Customer")
		} else {
			fmt.Println("Error: No guilds available")
		}

		//fmt.Println(guildID)
		//fmt.Println(ev.Author.Username)
		//fmt.Println(role)
		/*
			err = session.GuildMemberRoleAdd(guildID, ev.Author.ID, "1108944121630044172")
			if err != nil {
				fmt.Println("Error assigning role to user:", err)
				return
			}
		*/
		fmt.Printf("Assigned role '%s' to user %s\n", role, ev.Author.Username)
	}
}

func main() {
	token, err := getTokenFromCredentials("credentials.json")
	if err != nil {
		log.Fatal("Error retrieving bot token:", err)
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	memberList := &MemberList{
		Members:        make(map[string]bool),
		InviteCounters: make(map[string]int),
		InviteRole:     make(map[string]string),
	}

	session.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		onReady(s, event, memberList)
	})

	session.AddHandler(func(s *discordgo.Session, event *discordgo.MessageCreate) {
		onEvent(s, event, memberList)
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
