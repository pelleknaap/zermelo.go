package main

import (
	"fmt"
	"github.com/pelleknaap/zermelo"
	"strconv"
	"time"
)

func main() {
	// getting the times in seconds
	now := time.Now()
	start := now.Unix()
	end := now.Add(time.Hour * 24 * 2).Unix()

	// filling in all the data so we can call methods on it
	z := zermelo.ZermeloData{
		Start: strconv.Itoa(int(start)),
		End:   strconv.Itoa(int(end)),

		// fill-in your own school
		School: "ccg",

		// fill-in your own apikey, or use AuthCode to generate an apikey
		//Key:"your api key here",

		AuthCode: "your auth code here",
	}

	// call z.GetApiKey() if you want to get a api key with an auth code or "koppel code"
	// if you already got an api key yourself, you don't have to call the "GetApiKey" function
	err := z.GetApiKey()
	if err != nil {
		fmt.Println("error getting apikey")
		fmt.Println(err)
		return
	}

	fmt.Println(z.Key)

	// getting announcements
	err = z.GetAnnouncements()
	if err != nil {
		fmt.Println("error getting announcements")
		fmt.Println(err)
		return
	}

	for _, announcement := range z.Announcements.Data {
		fmt.Println(announcement.Title)
	}

	// getting the appointments
	err = z.GetAppointments()
	if err != nil {
		fmt.Println(err)
		return
	}

	// ranging over the appointment data
	// printing all the cancelled lessons
	for _, lesson := range z.Appointments.Data {
		if lesson.Cancelled == true || lesson.Valid == false {
			fmt.Println(lesson.Subject)
			fmt.Println("This leson has been cancelled")
			continue
		}
	}
}
