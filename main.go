package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {
	tgApi := os.Getenv("TG_API")
	bot, err := tgbotapi.NewBotAPI(tgApi)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upd, _ := bot.GetUpdatesChan(ucfg)
	// читаем обновления из канала
	for {
		select {
		case update := <-upd:
			UserName := update.Message.From.UserName

			ChatID := update.Message.Chat.ID

			Text := update.Message.Text + " tvar"

			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			// Ответим пользователю его же сообщением
			//reply := Text
			reply := getAudi()
			// Созадаем сообщение
			msg := tgbotapi.NewMessage(ChatID, reply)
			// и отправляем его
			bot.Send(msg)
		}

	}
}

func getAudi() string {
	url := "https://auto.ria.com/search/?body.id[0]=3&year[0].gte=2010&year[0].lte=2013&categories.main.id=1&brand.id[0]=48&model.id[0]=428&price.currency=1&drive.id[0]=1&abroad.not=0&custom.not=1&page=0&size=100"
	resp, err := http.Get(url)
	response := "not found"
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(body)))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "window.ria.server.ids") {
			response = line
		}
	}

	riaIdsDirty := strings.Split(response, "=")[1]
	riaIdsDirty = strings.Replace(riaIdsDirty, "[", "", -1)
	riaIdsDirty = strings.Replace(riaIdsDirty, "]", "", -1)
	riaIdsDirty = strings.Replace(riaIdsDirty, ";", "", -1)
	riaIdsClean := strings.Replace(riaIdsDirty, "\"", "", -1)
	//riaIds := strings.Split(riaIdsClean, ",")

	return riaIdsClean
}
