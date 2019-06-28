package zermelo

import "fmt"

type Zermelo struct {
	appointments *AppointmentsAPI
	announcements *AnnouncementsAPI
}

func ObtainAccessCode(authcode string, school string) (string, error) {
	c := core{
		school: school,
	}

	err := c.GetAccessToken(authcode)
	if err != nil {
		fmt.Println("err access token")
		fmt.Println(err)
		return "", err
	}

	return c.accessCode, nil
}

func Get(accessCode string, school string) (Zermelo, error) {
	c := core{
		accessCode:accessCode,
		school:school,
	}

	return Zermelo{
		appointments: newAppointments(&c),
		announcements: newAnnouncements(&c),
	}, nil
}

func (z Zermelo) Appointments() *AppointmentsAPI {
	return z.appointments
}

func (z Zermelo) Announcements() *AnnouncementsAPI {
	return z.announcements
}