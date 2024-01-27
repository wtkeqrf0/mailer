package main

import (
	"github.com/goccy/go-json"
	"mailer/internal/consumer"
	"os"
)

// Used only for test data
func main() {
	email := consumer.Email{
		To:          []string{"xolvpt@hldrive.com"},
		Subject:     "New Go Email",
		CopyTo:      []string{"yum@yandex.ru"},
		BlindCopyTo: []string{"sus@mail.ru", "admin@gmail.com"},
		Sender:      "Biba from Exnode <best@mail.ru>",
		ReplyTo:     "myemail@mail.ru",
		Parts: []consumer.Part{
			{
				Body:        "Hello friend!",
				ContentType: consumer.TextPlain,
			},
		},
		PartValues: map[string]any{"link": "https://gmail.com"},
		Files: []*consumer.File{
			{
				FilePath: "pussy.png",
				Name:     "image.jpg",
				MimeType: ".png",
				B64Data:  "something",
				Data:     []byte("something"),
				Inline:   true,
			},
		},
		Settings: consumer.ServiceSettings{
			Name:   "hello",
			From:   "singleSender",
			Locale: "ru",
		},
	}

	b, err := json.MarshalIndent(email, "  ", "   ")
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile("me.txt", b, 0644); err != nil {
		panic(err)
	}
}
