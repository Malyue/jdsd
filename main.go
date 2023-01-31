package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"jdsd/jdsd"
)

type UserInfo struct {
	key   string
	email string
}

var userInfos = []UserInfo{
	{key: "", email: ""},
}

var emailInfo = &jdsd.EmailInfo{
	ServerHost:   "smtp.qq.com",
	ServerPort:   25,
	FromEmail:    "",
	FromPassword: "",
	Recipient:    make([]string, 0),
}

func main() {
	s := gocron.NewScheduler()
	s.Every(1).Day().At("9:00").Do(func() {
		for _, user := range userInfos {
			userInfo, sysErr, userErr := jdsd.Exec(user.key)
			re := userInfo["re"].(map[string]interface{})
			emailInfo.Recipient = append(emailInfo.Recipient, user.email)

			if userErr != nil {
				fmt.Println(userErr)
				fmt.Println(sysErr)

				// send the error email to user
				if emailInfo.FromEmail != "" && emailInfo.FromPassword != "" {
					body := fmt.Sprintf("今日经典诵读未完成，出现以下问题，请查阅后根据提示进行：\n")
					body = body + "\n" + userErr.Error()
					err := jdsd.SendEmail("经典诵读工具如", body, emailInfo)
					if err != nil {
						fmt.Println(err)
						fmt.Println(user.email + "邮箱发送错误")
					}
				}
				return
			}
			if emailInfo.FromEmail != "" && emailInfo.FromPassword != "" {
				body := fmt.Sprintf("今日经典诵读已完成，共获得积分：%v/30 ，当前总积分为：%v", re["per_day_credits"], re["credits"])
				err := jdsd.SendEmail("经典诵读工具人", body, emailInfo)
				if err != nil {
					fmt.Println(err)
					fmt.Println(user.email + "邮箱发送错误")
					return
				}
			}
		}
	})
	<-s.Start()
}
