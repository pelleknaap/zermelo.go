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
	Start  string
	End    string

	// the "koppel code", not the api key
	AuthCode string

	// the api key that can be retrieved by using the GetApiKey() function
	Key           string
	Appointments  Appointments
	Announcements Announcements
}

// json wrapper structs
type JSONWrapperAppointments struct {
	Response *Appointments
}

type JSONWrapperAnnouncements struct {
	Response *Announcements
}

type ApiKeyWrapper struct {
	AccessToken string `json:"access_token"`
}

type Appointments struct {
	Status    int
	TotalRows int
	Data      []Lesson
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

type Announcements struct {
	Status    int
	Message   string
	StartRow  int
	EndRow    int
	TotalRows int
	Data      []Announcement
}

type Announcement struct {
	Id    int
	Start int
	End   int
	Title string
	Text  string
}

// GetAppointments
// Gets all the appointments from the Zermelo api
// used z.Start & z.End to determine the period of the appointments
// Needs z.Key to access the API, will return an error if there isn't one
// Makes a request to the Zermelo api and fills the z.Appointments slice
func (z *ZermeloData) GetAppointments() error {
	if z.Start == "" || z.End == "" || z.Key == "" {
		return errors.New("Not all needed variables are present, check the z.Start, z.End & the z.Key variables")
	}

	// make the url
	var reqUrl strings.Builder
	fmt.Fprintf(&reqUrl, "https://%s.zportal.nl/api/v3/", z.School)
	reqUrl.WriteString("appointments?user=~me")
	reqUrl.WriteString("&start=" + z.Start)
	reqUrl.WriteString("&end=" + z.End)
	reqUrl.WriteString("&access_token=" + z.Key)

	// Get the data
	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return errors.Wrap(err, "Error getting json data")
	}

	// Check if nothing went wrong with the request
	if resp.StatusCode != http.StatusOK {
		return errors.New("Wrong statuscode returned")
	}

	defer resp.Body.Close()

	var appointmentsWrapper JSONWrapperAppointments

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&appointmentsWrapper)
	if err != nil {
		return errors.Wrap(err, "Error decoding json")
	}

	z.Appointments = *appointmentsWrapper.Response

	// Check if Zermelo's api also returned 200 as a Status
	if z.Appointments.Status != 200 {
		return errors.New("Returned status code isn't 200, it is: " + strconv.Itoa(z.Appointments.Status))
	}

	return nil
}

// GetAnnouncements
// Gets all announcements from the Zermelo API
// The z.Key variable needs to be present, an error will be returned if it isn't
// fills the z.Announcements slice, returns error if something went wrong
func (z *ZermeloData) GetAnnouncements() error {
	if z.Key == "" {
		return errors.New("Not all needed variables are present, check the z.Start, z.End & the z.Key variables")
	}

	// Create the url for the request
	var reqUrl strings.Builder
	fmt.Fprintf(&reqUrl, "https://%s.zportal.nl/api/v2/", z.School)
	reqUrl.WriteString("announcements?user=~me")
	reqUrl.WriteString("&current=true")
	reqUrl.WriteString("&access_token=" + z.Key)

	// create and execute the request
	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return errors.Wrap(err, "Error getting json data")
	}

	// check if nothing went wrong
	if resp.StatusCode != http.StatusOK {
		return errors.New("Wrong statuscode returned")
	}

	defer resp.Body.Close()

	var announcementsWrapper JSONWrapperAnnouncements

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&announcementsWrapper)
	if err != nil {
		return errors.Wrap(err, "Error decoding json")
	}

	// set the announcements slice to the response data
	z.Announcements = *announcementsWrapper.Response

	if z.Announcements.Status != 200 {
		return errors.New("Returned status code isn't 200, it is: " + strconv.Itoa(z.Announcements.Status))
	}

	return nil
}

// GetApiKey gets the api key needed to interact with the api from the zermelo api
// The z.Authcode variable needs to be present, otherwise an error will be returned
func (z *ZermeloData) GetApiKey() error {
	if z.AuthCode == "" {
		return errors.New("Please fill-in the auth code before trying to get an apikey")
	}

	resp, err := http.PostForm("https://"+z.School+".zportal.nl/api/v2/oauth/token",
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
