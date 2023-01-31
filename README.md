## gzhu_jdsd

### ![Github stars](https://img.shields.io/badge/Language-go-brightgreen)

### Introduction

The program implemented by Go

Aims to complete the gzhu_jdsd mini apps.:bulb:

It will run automatically at :nine: CST every day :alarm_clock:

Supported for sending the email for your mailbox:mailbox:

Supported for exec more than one user.:star2:

### How to use

First,you should have sth knowledge about Packet capture,using tool like Fiddler，Charles to get the miniprog's key.:closed_book:

Then you should set you key and email into file `main.go`.(userInfos):key:

If you need the program to send the success or failure email for your mailbox,you should put your mailbox in the variable `emailInfo`

in `main.go`.:email:

Then if you set all the config in the program,and you can execute `go build main.go` to start the program.

:link:Email setting reference:https://cloud.tencent.com/developer/article/2177098

### Tips

- If you have sth knowledge of deploy,you can deploy the program for you  server.
- But it has sth problem in deploying,if you deploy in docker,it can't run automatically by the gocron,you can try to use ==crontab== to controller it.
- The variable `emailInfo` is the info of sender，the receiver's mailbox is set in the variables `userInfos` 
