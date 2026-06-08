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

	if len(historiales[chatID]) > 12 {
		historiales[chatID] = append(
			historiales[chatID][:1],
			historiales[chatID][len(historiales[chatID])-10:]...,
		)
	}

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

	bot.Handle("/ayuda", func(c tele.Context) error {
		return c.Send("🤖 *Comandos disponibles:*\n\n/start → Iniciar conversación\n/ayuda → Ver esta lista\n/reset → Reiniciar conversación\n\nO simplemente escribime tu consulta y te respondo.", &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	})

	bot.Handle("/reset", func(c tele.Context) error {
		delete(historiales, c.Chat().ID)
		return c.Send("🔄 Conversación reiniciada. ¿En qué te puedo ayudar?")
	})

	bot.Handle(tele.OnPhoto, func(c tele.Context) error {
		return c.Send("📝 Por el momento solo proceso texto. Escribime tu consulta y te ayudo.")
	})

	bot.Handle(tele.OnDocument, func(c tele.Context) error {
		return c.Send("📝 Por el momento solo proceso texto. Escribime tu consulta y te ayudo.")
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		bot.Notify(c.Chat(), tele.Typing)
		mensaje := c.Text()
		chatID := c.Chat().ID
		respuesta := obtenerRespuesta(chatID, mensaje)
		return c.Send(respuesta)
	})

	log.Println("Bot corriendo...")
	bot.Start()
}
