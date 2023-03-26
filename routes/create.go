package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"os"

	"github.com/google/uuid"
	"github.com/memphisdev/memphis.go"
	neon "github.com/stationapi/station-upload/db"
	"github.com/stationapi/station-upload/session"
	"gorm.io/gorm"
)

type message struct {
	name string
	id   string
}

func Create(w http.ResponseWriter, r *http.Request, db gorm.DB) error {
	cookie, cookieErr := r.Cookie("station")

	if cookieErr != nil {
		http.Error(w, "you are not authenticated", http.StatusForbidden)

		return cookieErr
	} 

	githubId, authErr := session.GetSession(cookie.Value)

	if authErr != nil {
		http.Error(w, "you are not authenticated", http.StatusForbidden)

		return authErr
	}

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

	if website.Bumps > 0 {
		http.Error(w, "the bumps field was already defined in the request body", http.StatusBadRequest)

		return errors.New("the bumps was already defined in the request body")
	}

	if website.Created > 0 {
		http.Error(w, "the bumps field was already defined in the request body", http.StatusBadRequest)

		return errors.New("the created field was already defined in the request body")
	}

	if website.Owner > 0 {
		http.Error(w, "the owner field was already defined in the request body", http.StatusBadRequest)

		return errors.New("the owner field was already defined in the request body")
	} 

  if website.LastBumped > 0 {
    http.Error(w, "the last_bumped field was already defined in the request body", http.StatusBadGateway)

    return errors.New("the last_bumped field was already defined in the request body")
  }

	website.Id = uuid.NewString()
	website.Bumps = 0
	website.Created = int(time.Now().Unix())
	website.Owner = githubId
  website.LastBumped = 0

	neon.CreateWebsite(website, db)

	conn, err := memphis.Connect(
		"memphis-rest-gateway.memphis.svc.cluster.local",
		os.Getenv("USER"),
		os.Getenv("TOKEN"),
	) 

	if err != nil {
		http.Error(w, "there was an error creating the website", http.StatusInternalServerError)

		return err
	}

	defer conn.Close()

	message := message{
		name: website.Name,
		id:   website.Id,
	}

	stringified, err := json.Marshal(message)

	if err != nil {
		http.Error(w, "there was an error creating the website", http.StatusInternalServerError)

		return err
	}

	producer, err := conn.CreateProducer("new_website", uuid.NewString())

	err = producer.Produce(stringified, memphis.MsgHeaders(memphis.Headers{}))

	if err != nil {
		http.Error(w, "there was an error creating the website", http.StatusInternalServerError)

		return err
	}

	w.WriteHeader(200)
	w.Write([]byte(website.Id))

	return nil
}
