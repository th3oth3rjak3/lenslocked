package views

import (
	"log"

	"lenslocked/models/usersModel"
)

const (
	// AlertLevelError is used to render level-danger alerts for bootstrap
	AlertLevelError = "danger"

	// AlertLevelWarning is used to render level-warning alerts for bootstrap
	AlertLevelWarning = "warning"

	// AlertLevelInfo is used to render level-info alerts for bootstrap
	AlertLevelInfo = "info"

	// AlertLevelSuccess is used to render level-success alerts for bootstrap
	AlertLevelSuccess = "success"

	// AlertMessageGeneric is used to provide a generic message when internal
	// server errors or random, unhandled errors are encoutered in our backend code.
	AlertMessageGeneric = "Something went wrong, please try again." +
		" If the problem persists, please contact support."
)

// Alert contains the alert level and message content to be displayed to the user
type Alert struct {
	Level   string
	Message string
}

// Data contains data to be rendered on the template. If an alert exists
// then it will be set in the Alert property. Any other payload data to
// be rendered on the page will be housed in the Payload property.
type Data struct {
	Alert   *Alert
	User    *usersModel.User
	Payload interface{}
}

// setAlert takes any error, both public and private, and sets an alert
// in the Data object based on the error. If logging the alert is desired,
// set logErr to true.
func (d *Data) SetAlert(err error, logErr bool) {
	if logErr {
		log.Println(err)
	}
	d.Alert = &Alert{Level: AlertLevelError}
	if pErr, ok := err.(PublicError); ok {
		d.Alert.Message = pErr.Public()
	} else {
		d.Alert.Message = AlertMessageGeneric
	}
}

type PublicError interface {
	error
	Public() string
}
