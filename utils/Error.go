package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
	"net/http"
	"strings"
)

var generator *shortid.Shortid

type clientError struct {
	ID            string `json:"id"`
	MessageToUser string `json:"messageToUser"`
	DeveloperInfo string `json:"developerInfo"`
	Err           string `json:"error"`
	StatusCode    int    `json:"statusCode"`
	IsClientError bool   `json:"isClientError"`
}

func newClientError(err error, statusCode int, messageToUser string, additionalInfoForDevs ...string) *clientError {
	additionalInfoJoined := strings.Join(additionalInfoForDevs, "\n")
	if additionalInfoJoined == "" {
		additionalInfoJoined = messageToUser
	}

	errorID, _ := generator.Generate()
	var errString string
	if err != nil {
		errString = err.Error()
	}
	return &clientError{
		ID:            errorID,
		MessageToUser: messageToUser,
		DeveloperInfo: additionalInfoJoined,
		Err:           errString,
		StatusCode:    statusCode,
		IsClientError: true,
	}
}

// RespondError sends an error message to the API caller and logs the error
func RespondError(w http.ResponseWriter, statusCode int, err error, messageToUser string, additionalInfoForDevs ...string) {
	logrus.Errorf("status: %d, message: %s, err: %+v ", statusCode, messageToUser, err)
	clientError := newClientError(err, statusCode, messageToUser, additionalInfoForDevs...)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientError); err != nil {
		logrus.Errorf("Failed to send error to caller with error: %+v", err)
	}
}
