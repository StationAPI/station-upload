package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"context"

	events "github.com/ferretcode-freelancing/fc-bus"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go"
	neon "github.com/stationapi/station-upload/db"
	"gorm.io/gorm"
)

type message struct {
	name string	
	id string
}

func Create (w http.ResponseWriter, r *http.Request, db gorm.DB) error {
	website := neon.Website{}	

	err := ProcessBody(r.Body, &website)

	if err != nil {
		http.Error(w, "there was an error processing the request body", http.StatusInternalServerError)

		return err
	}

	if website.Id != "" {
		http.Error(w, "the id was already defined in the request body", http.StatusBadRequest)

		return errors.New("the id was already defined in the request body")
	}

	website.Id = uuid.NewString() 

	neon.CreateWebsite(website, db)

	ctx := context.Background()

	bus := events.Bus{
		Channel: "new_website",
		ClientId: uuid.NewString(),
		Context: ctx,
		TransportType: kubemq.TransportTypeGRPC,
	}

	client, err := bus.Connect()

	if err != nil {
		http.Error(w, "there was an error creating the website", http.StatusInternalServerError)

		return err
	}

	message := message{
		name: website.Name,
		id: website.Id,
	}

	stringified, err := json.Marshal(message)

	if err != nil {
		http.Error(w, "there was an error creating the website", http.StatusInternalServerError)

		return err
	}

	_, sendErr := client.Send(ctx, kubemq.NewQueueMessage().
		SetId(uuid.NewString()).
		SetChannel("new_website").
		SetBody(stringified))

	if sendErr != nil {
		http.Error(w, "there was an error creating the website", http.StatusInternalServerError)

		return err
	}

	w.WriteHeader(200)
	w.Write([]byte(website.Id))

	return nil
}
