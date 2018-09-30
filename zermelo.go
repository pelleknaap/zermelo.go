package zermelo

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ZermeloData is the struct that holds all the methods and data that is needed by the methods
type ZermeloData struct {
	School string
	Start string
	End string

	// the "koppel code", not the api key
	AuthCode string

	// the api key that can be retrieved by using the GetApiKey() function
	Key string
	Appointments Appointments
}

// json wrapper structs
type AppointmentWrapper struct {
	Response *Appointments
}

type ApiKeyWrapper struct {
	AccessToken string `json:"access_token"`
}

type Appointments struct {
	Status int
	TotalRows int
	Data []Lesson
}

type Lesson struct {
	ID                  int      `json:"id"`
	AppointmentInstance int      `json:"appointmentInstance"`
	Start               int      `json:"start"`
	End                 int      `json:"end"`
	StartTimeSlot       int      `json:"startTimeSlot"`
	EndTimeSlot         int      `json:"endTimeSlot"`
	Subjects            []string `json:"subjects"`
	Teachers            []string `json:"teachers"`
	Groups              []string `json:"groups"`
	GroupsInDepartments []int    `json:"groupsInDepartment"`
	Locations           []string `json:"locations"`
	LocationsOfBranch   []int    `json:"locationsOfBranch"`
	Type                string   `json:"type"`
	Remark              string   `json:"remark"`
	Valid               bool     `json:"valid"`
	Cancelled           bool     `json:"cancelled"`
	Modified            bool     `json:"modified"`
	Moved               bool     `json:"moved"`
	New                 bool     `json:"new"`
	ChangeDescription   string   `json:"changeDescription"`
}

func (z *ZermeloData) GetAppointments() (error) {
	var appointments Appointments

	var url strings.Builder
	fmt.Fprintf(&url, "https://%s.zportal.nl/api/v3/", z.School)
	url.WriteString("appointments?user=~me")
	url.WriteString("&start="+z.Start)
	url.WriteString("&end="+z.End)
	url.WriteString("&access_token="+z.Key)

	resp, err := http.Get(url.String())
	if err != nil {
		return errors.Wrap(err, "Error getting json data")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Wrong statuscode returned")
	}

	defer resp.Body.Close()

	var appointmentsWrapper AppointmentWrapper

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&appointmentsWrapper)
	if err != nil {
		return errors.Wrap(err, "Error decoding json")
	}

	z.Appointments = *appointmentsWrapper.Response

	if z.Appointments.Status != 200 {
		return errors.New("Returned status code isn't 200, it is: " + strconv.Itoa(appointments.Status))
	}

	return nil
}

func (z *ZermeloData) GetApiKey() (error) {
	if z.AuthCode == "" {
		return errors.New("Please fill-in the auth code before trying to get an apikey")
	}

	resp, err := http.PostForm("https://" + z.School + ".zportal.nl/api/v2/oauth/token",
		url.Values{"grant_type": {"authorization_code"}, "code": {z.AuthCode}})

	if err != nil {
		return errors.Wrap(err, "Error getting api key")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Wrong statuscode returned, code: " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	var apiKeyWrapper ApiKeyWrapper

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&apiKeyWrapper)
	if err != nil {
		return errors.Wrap(err, "Error decoding json")
	}

	z.Key = apiKeyWrapper.AccessToken
	return nil
}