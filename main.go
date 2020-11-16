package main

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"

	"github.com/jdxj/bilibili/config"
	"github.com/jdxj/bilibili/models"
	"github.com/jdxj/bilibili/utils"
)

func main() {
	cro := cron.New()
	users := config.GetUsers()

	for i, user := range users {
		b, err := models.NewBiliBili(user.Cookie)
		if err != nil {
			log.Printf("i: %d, err: %s\n", i, err)
			return
		}

		// 避免使用 for 的 user 局部变量
		addr := user.Email
		entryID, err := cro.AddFunc("0 8 * * *", func() {
			if err := b.Run(); err != nil {
				_ = utils.SendMessage(addr, b.Subject(), fmt.Sprintf("恰硬币失败: %s", err))
			} else {
				_ = utils.SendMessage(addr, b.Subject(), b.Content())
			}
		})

		if err != nil {
			log.Printf("%s\n", err)
			return
		}
		log.Printf("add func: %d\n", entryID)
	}

	cro.Run()
}
