package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	tele "gopkg.in/telebot.v3"
)

var infotienda = `
Sos el asistente virtual de "Ropa Style", una tienda online de ropa en Buenos Aires.
Información del negocio:
- Productos: remeras ($15000), pantalones ($25000), camperas ($45000)
- Talles disponibles: S, M, L, XL
- Horario de atención humana: lunes a viernes de 9 a 18hs
- Envíos: a todo el país por Andreani, gratis en CABA
- Medios de pago: MercadoPago y transferencia bancaria
- Instagram: @ropastyle
Nunca inventes productos, precios ni información que no esté en este texto. Si no tenés el dato, decí que lo vas a consultar con el equipo.
`

var historiales = map[int64][]openai.ChatCompletionMessage{}
var clienteGroq *openai.Client

func obtenerRespuesta(chatID int64, mensaje string) string {
	if _, ok := historiales[chatID]; !ok {
		historiales[chatID] = []openai.ChatCompletionMessage{
			{Role: "system", Content: infotienda},
		}
	}

	historiales[chatID] = append(historiales[chatID], openai.ChatCompletionMessage{
		Role:    "user",
		Content: mensaje,
	})

	resp, err := clienteGroq.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    "llama-3.1-8b-instant",
			Messages: historiales[chatID],
		},
	)
	if err != nil {
		log.Println("Error al llamar a Groq:", err)
		return "Lo siento, hubo un error. Intentá de nuevo."
	}

	respuesta := resp.Choices[0].Message.Content

	historiales[chatID] = append(historiales[chatID], openai.ChatCompletionMessage{
		Role:    "assistant",
		Content: respuesta,
	})

	return respuesta
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	config := openai.DefaultConfig(os.Getenv("GROQ_API_KEY"))
	config.BaseURL = "https://api.groq.com/openai/v1"
	clienteGroq = openai.NewClientWithConfig(config)

	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("👋 ¡Hola! Soy el asistente virtual de *Ropa Style*.\n\nPuedo ayudarte con:\n• Precios y productos\n• Talles disponibles\n• Información de envíos\n• Medios de pago\n\n¿En qué te puedo ayudar hoy?", &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		mensaje := c.Text()
		chatID := c.Chat().ID
		respuesta := obtenerRespuesta(chatID, mensaje)
		return c.Send(respuesta)
	})

	log.Println("Bot corriendo...")
	bot.Start()
}
