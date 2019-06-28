package zermelo

import "fmt"

type AppointmentsAPI struct {
	c *core
}

// Lesson holds all the data of a lesson, it's used in the Appointments struct
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

type appointmentReplyWrapper struct {
	Response struct{
		TotalCount int `json:"totalCount"`
		Offset     int
		Count      int
		Data       []Lesson
	}
}

func newAppointments(co *core) *AppointmentsAPI {
	return &AppointmentsAPI{c: co}
}

func (a *AppointmentsAPI) Get(start string, end string) ([]Lesson, error){
	var appointmentsReply appointmentReplyWrapper

	err := a.c.Get(fmt.Sprintf("appointments?user=~me&start=%s&end=%s", start, end), &appointmentsReply)
	if err != nil {
		return nil, err
	}

	return appointmentsReply.Response.Data, nil
}
