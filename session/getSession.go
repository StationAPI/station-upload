package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type SessionObject struct {
	Id string `json:"cookie"`
}

type User struct {
	GithubId int `json:"github_id"`
}

func GetSessionCache(sid string) string {
	ip := os.Getenv("STATION_SESSION_CACHE_SERVICE_HOST")
	port := os.Getenv("STATION_SESSION_CACHE_SERVICE_PORT")

	url := fmt.Sprintf("http://%s:%s/get?sid=%s", ip, port, sid)

	return url 
}

func GetSession(sid string) (int, error) {
  client := http.Client{}

  req, err := http.NewRequest(
    "GET",
    GetSessionCache(sid),
		nil,
  )

  if err != nil {
    return 0, err
  }

  res, err := client.Do(req)

  if err != nil {
    return 0, err
  }

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return 0, errors.New(err.Error())
	}

  if string(body) == "null" {
    return 0, errors.New("there was an error fetching the session") 
  }

	user := User{}

	err = json.Unmarshal(body, &user)

	if err != nil {
		return 0, errors.New("there was an error processing the session")
	}

	return user.GithubId, nil
}

