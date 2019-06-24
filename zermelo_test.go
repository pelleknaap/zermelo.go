package zermelo

import (
	"strconv"
	"testing"
	"time"
)

// setting up the needed data
var z = ZermeloData{
	Start:  strconv.Itoa(int(time.Now().Unix())),
	End:    strconv.Itoa(int(time.Now().Add(time.Hour * 24 * 3).Unix())),
	School: "school here",
	Key:    "key here",
}

func TestAppointments(t *testing.T) {
	// getting the appointments
	err := z.GetAppointments()
	if err != nil {
		t.Errorf("Expected an appointment slice, however an error occurred: %b", err)
	}

	if len(z.Appointments.Data) < 1 {
		t.Error("Appointment slice is 0, this test might fail if it's a holiday")
	}
}
