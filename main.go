package main

import (
	"io/ioutil"
	"log"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	qrcode "github.com/skip2/go-qrcode"

	utopiago "github.com/Sagleft/utopialib-go"
)

const (
	qrFilePath = "qr.png"
)

type Config struct {
	TelegramBotToken string
	UtpToken         string
	UtpPort          int
	PublicKey        string
}

func main() {

	//read congig
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(configuration.TelegramBotToken)
	fmt.Println(configuration.UtpToken)
	fmt.Println(configuration.UtpPort)

	// bot-token

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// ini channel
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updatesChann := bot.GetUpdatesChan(ucfg)

	//utp
	client := utopiago.UtopiaClient{
		Protocol: "http",
		Token:    configuration.UtpToken,
		Host:     "127.0.0.1",
		Port:     configuration.UtpPort,
	}

	//send bool

	// update
	for {
		select {

		case update := <-updatesChann:
			// User bot
			UserName := update.Message.From.UserName

			// ID chat.

			ChatID := update.Message.Chat.ID

			// Text massage user
			Text := update.Message.Text

			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			//commands

			switch Text {

			case "/ucode":

				//send to channel

				//ucodeEncode

				Ucodius64, err := client.UCodeEncode(configuration.PublicKey, "BASE64", "PNG", 256)
				if err != nil {
					log.Println(err)
				}
				//decode
				log.Println(Ucodius64)
				log.Println("Ucode COMPLIIIITE")

				err = ioutil.WriteFile("sample.b64", []byte(Ucodius64), 0644)

				if err != nil {
					log.Println(err)
				}

				b64Data, err := ioutil.ReadFile("sample.b64")
				if err != nil {
					log.Println(err)
				}

				// base64 to png
				outPngData, err := base64.StdEncoding.DecodeString(string(b64Data))

				if err != nil {
					log.Println(err)
				}

				//send

				photoFileBytes := tgbotapi.FileBytes{
					Name:  "ucode",
					Bytes: outPngData,
				}

				msg := tgbotapi.NewPhoto(ChatID, photoFileBytes)

				bot.Send(msg)

			case "/qrcode":

				//send to channel

				//create Qr

				err := qrcode.WriteFile(configuration.PublicKey, qrcode.Medium, 256, qrFilePath)
				if err != nil {
					log.Println("write error")
				}
				log.Println("CREATE QR")
				//Send Qr

				photoBytes, err := ioutil.ReadFile(qrFilePath)
				if err != nil {
					panic(err)
				}
				photoFileBytes := tgbotapi.FileBytes{
					Name:  "qr",
					Bytes: photoBytes,
				}

				msg := tgbotapi.NewPhoto(ChatID, photoFileBytes)

				bot.Send(msg)

			case "/how":

				log.Println("/how")

				msg := tgbotapi.NewMessage(ChatID, "create config.json and use")

				bot.Send(msg)

			default:

				fmt.Println("commands")

				reply := "Commands:\n /ucode \n /qrcode"
				msg := tgbotapi.NewMessage(ChatID, reply)

				bot.Send(msg)

			}

		}

	}
}
