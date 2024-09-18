package main

import (
	"encoding/json"
	"mailer/pkg/mail"
	"os"
	"testing"
)

func TestStruct(t *testing.T) {
	f, err := os.Create("msg.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			t.Error(err)
		}
	}()

	if err = json.NewEncoder(f).Encode(mail.Parsable{
		To:          []string{"matvey-sizov@mail.ru"},
		Subject:     "thank you!",
		CopyTo:      nil,
		BlindCopyTo: nil,
		Parts: []mail.Part{
			{
				ContentType: mail.TextPlain,
				Body:        []byte("hello, friend! I am very stupid {{ .String }}"),
			},
		},
		PartValues: map[string]any{
			"String": "(secret \"NOT\")",
		},
		Files: []*mail.File{
			{
				FilePath: "Dockerfile",
			},
		},
		Settings: nil,
	}); err != nil {
		t.Error(err)
	}
}
