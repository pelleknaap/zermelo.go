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
	School: "schoolhere",
	Key:    "keyhere",
}

func TestZermeloData_GetAppointments(t *testing.T) {
	badData := ZermeloData{}

	err := badData.GetAnnouncements()

	if err == nil {
		t.Error("GetAnnouncements() err = nil; want non-nil")
	}

	if err.Error() != "Not all needed variables are present, check the z.Start, z.End & the z.Key variables" {
		t.Errorf("GetAnnouncements() err = %s; want Not all needed variables are present, check the z.Start, z.End & the z.Key variables", err.Error())
	}

	err = z.GetAnnouncements()
	if err == nil {
		t.Error("GetAnnouncements() err = nil; want non-nil")
	}

	if err.Error() != "Wrong statuscode returned" {
		t.Errorf("GetAnnouncements() err = %s; want Wrong statuscode returned", err.Error())
	}
}

func TestZermeloData_GetAnnouncements(t *testing.T) {
	badData := ZermeloData{}

	err := badData.GetAppointments()

	if err == nil {
		t.Error("GetAppointments() err = nil; want non-nil")
	}

	if err.Error() != "Not all needed variables are present, check the z.Start, z.End & the z.Key variables" {
		t.Errorf("GetAppointments() err = %s; want Not all needed variables are present, check the z.Start, z.End & the z.Key variables", err.Error())
	}
}