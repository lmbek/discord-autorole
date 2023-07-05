This one is used for Beksoft discord server
It tracks who is already member of the discord group
When a new member is added, depending on the invite link, they are assigned a role.
They are also assigned into a new room called "Welcome! Tell us about you"

This project will be crosscompiled from windows using:
set GOOS=linux
set GOARCH=amd64
go build -o autorole-bot

These settings (GOOS and GOARCH) will only be temporary

Then the autorole-bot linux executable is sent to a linux server, that should run it 24/7
I have put it on my LMBEK server (for test purposes so far - 05-07-2023)

The credentials.json and invite-roles.json is also sent to the linux server, they should sit in the same folder as the executable