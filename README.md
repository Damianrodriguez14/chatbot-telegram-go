# Ropa Style - Asistente Virtual IA para Telegram

Bot de Telegram con inteligencia artificial para atención al cliente de comercios. Responde consultas sobre productos, precios, envíos y medios de pago de forma automática.

## Tecnologías
- Go
- Telegram Bot API
- Groq AI (LLaMA 3.1)
- Variables de entorno para configuración segura

## Funcionalidades
- Responde consultas de clientes automáticamente
- Mantiene historial de conversación por usuario
- Comando /start con mensaje de bienvenida
- Fácilmente configurable para cualquier comercio

## Configuración
1. Clonar el repositorio
2. Crear archivo `.env` con las variables:
TELEGRAM_TOKEN=tu_token
GROQ_API_KEY=tu_api_key
3. Ejecutar con `go run main.go`

## Casos de uso
Ideal para tiendas online, restaurantes y comercios que quieran automatizar la atención al cliente por Telegram.