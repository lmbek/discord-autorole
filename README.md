# discord-autorole

## What is discord-autorole go bot?
This discord bot can automatically assign a role to a new user based on invitation link value.
It is made with Go and can be used for your own discord bot. You can freely take it and modify it to your needs.
Since the API token is your own, this script should work if you add a credentials.json and invite-roles.json.

## How to use
Download the files from this repository (and then maybe make your own or keep it local, you decide).
Then set up Go (if you dont know Go, follow some youtube tutorial you will get started within 1 hour i am sure)
After you have installed Go, you can get started with the domain:

You can follow the example jsons. First make a discord application (bot), then find token and put it into credentials.json. 
Then find out what permissions you need, then decide how you want the authentication to work, 
here you can get a link based on permissions and assign the bot to a server directly. 
Then you add in the roleID's and invitation links (only the value of the link, so everything after discord.com etc...).
And lastly after you added the roleID's and invitation links, you test it. (first test it in a separated room, just in case).
Now you have a working bot, you can use on your server, that automatically assigns roles for your new users, customers, friends or whatever.

I use it for my own discord server

This bot has 3 parts:
1) Connect-test
   * Connects to discord and disconnect right after
2) onready-give-role
   * Is a hardcoded script that connects to discord and gives a certain user a certain role, then it runs until being interrupted by user.
3) autorole-onjoin
   * Connects to discord and makes a list of users in the server/guild (only works with 1 server/guild so far), then it use that list of users as data
   * Adds an event-listener for OnEvent -> CreateMessage as this was the only way i could register the new user. It uses the data from the OnReady event with all the users, then it sees the difference between the users. It does this to find the invite link, as the invite link should be the difference between before and now (by checking every member). This operation is custom made and heavy, but it was necessary because i could not get invite link otherwise.
   * Then after it got the invite link for the new user, it checks what invite link was used and assign the role id to the user. This only works when joining the first time or not being member of server/guild on startup of script, it will work again after restarting script.

## Please Note
invite-roles.json: links.value is the value from the discord link that is shared with no expire date. It needs to be manually added in just like role name and role id

## Author
Lars M Bek (Main author) - June 2023

ChatGPT (junior-developer) - June 2023

## Further Development
* Adding functionality like:
  * New room for certain type of role, where only that person and the server owner is assigned
  * Added to certain room based on invite link (so we dont have to assign role manually like everything else in this bot)