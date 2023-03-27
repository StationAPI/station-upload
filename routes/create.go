package routes

import (
	"encoding/json"
	"net/http"
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

  website.Owner = githubId
	
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
