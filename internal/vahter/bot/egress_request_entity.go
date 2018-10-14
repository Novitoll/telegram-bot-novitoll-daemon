package bot

import (
	"os"
)

/*
	Telegram Outbound (sendMessage command) request body

	Example: URL / Text
	{
	  "chat_id": 345019684,
	  "text": "Bot reply",
	  "parse_mode": "Markdown",
	  "disable_web_page_preview": true,
	  "disable_notification": false,
	  "reply_to_message_id": 28
	}
*/

var (
	GET = "GET"
	POST = "POST"
	TELEGRAM_URL = "https://api.telegram.org/bot%s/sendMessage"
	TELEGRAM_TOKEN = "123"
	ParseModeMarkdown = "Markdown"
	ParseModeHTML = "HTML"
)

func init() {
	if u, ok := os.LookupEnv("TELEGRAM_URL"); ok {
		TELEGRAM_URL = u
	}
	if t, ok := os.LookupEnv("TELEGRAM_TOKEN"); ok {
		TELEGRAM_TOKEN = t
	}
}

// https://core.telegram.org/bots/api#sendmessage
type BotEgressRequest struct {
	ChatId			int `json:"chat_id"`
	Text			string `json:"text"`
	ParseMode		string `json:"parse_mode"`
	DisableWebPagePreview	bool `json:"disable_web_page_preview"`
	DisableNotification	bool `json:"disable_notification"`
	ReplyToMessageId	int `json:"reply_to_message_id"`
	ReplyMarkup 		*BotForceReply
}

// https://core.telegram.org/bots/api#forcereply
type BotForceReply struct {
	ForceReply 	bool `json:"force_reply"`
	Selective 	bool `json:"selective"`
}