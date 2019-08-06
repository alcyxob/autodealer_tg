package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
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
	url := "https://auto.ria.com/search/?categories.main.id=1&brand.id[0]=6&model.id[0]=49&year[0].gte=2009&year[0].lte=2012&region.id[0]=7&region.id[1]=10&gearbox.id[0]=2&gearbox.id[1]=3&gearbox.id[2]=4&gearbox.id[3]=5&drive.id[0]=1&abroad.not=0&custom.not=1&page=0&size=100"
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
