package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ZermeloData struct {
	School string
	Start string
	End string
	Key string
	Appointments Appointments
}

type AppointmentWrapper struct {
	Response *Appointments
}

func (z *ZermeloData) GetAppointments() (Appointments, error) {
	fmt.Println("lol")
	var appointments Appointments

	var url strings.Builder
	fmt.Fprintf(&url, "https://%s.zportal.nl/api/v3/", z.School)
	url.WriteString("appointments?user=~me")
	url.WriteString("&start="+z.Start)
	url.WriteString("&end="+z.End)
	url.WriteString("&access_token="+z.Key)

	fmt.Println("url: ", url.String())

	resp, err := http.Get(url.String())
	if err != nil {
		fmt.Println("lol1")
		return appointments, errors.Wrap(err, "Error getting json data")
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("lol2")
		return appointments, errors.New("Wrong statuscode returned")
	}

	defer resp.Body.Close()

	var appointmentsWrapper AppointmentWrapper

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&appointmentsWrapper)
	if err != nil {
		fmt.Print(err)
		return appointments, errors.Wrap(err, "Error decoding json")
	}

	appointments = *appointmentsWrapper.Response

	if appointments.Status != 200 {
		fmt.Println("lol3")
		return appointments, errors.New("Returned status code isn't 200")
	}

	fmt.Println("lol4")



	return appointments, nil
}

type Appointments struct {
	Status int
	TotalRows int
	Data []Lesson
}

type Lesson struct {
	Subjects []string
	Teachers []string
	Locations []string
	Cancelled bool
	Valid bool
	ChangeDescription string
}

func main() {

	// testing the api functions
	now := time.Now()
	start := now.Unix()
	end := now.Add(time.Hour * 24 * 2).Unix()
	accessToken := "u5sv8au3gt7j1tgdinv60nbu09"

	base := ZermeloData{
		Start:strconv.Itoa(int(start)),
		End:strconv.Itoa(int(end)),
		School:"ccg",
		Key:accessToken,
	}

	appointments, err := base.GetAppointments()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, lesson := range appointments.Data {
		if lesson.Cancelled == true || lesson.Valid == false {
			fmt.Println("cancelled")
			continue
		}
		fmt.Println(lesson.Subjects)
	}
}
