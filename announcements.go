package zermelo

import "fmt"

type AnnouncementsAPI struct {
	c *core
}

type Announcement struct {
	Id    int
	Start int
	End   int
	Title string
	Text  string
}

type announcementReplyWrapper struct {
	Response struct{
		TotalCount int `json:"totalCount"`
		Offset     int
		Count      int
		Data       []Announcement
	}
}

func newAnnouncements(co *core) *AnnouncementsAPI {
	return &AnnouncementsAPI{c: co}
}

func (a *AnnouncementsAPI) Get(current bool) ([]Announcement, error){
	var announcementsReply announcementReplyWrapper

	err := a.c.Get(fmt.Sprintf("announcements?user=~me&current=%t", current), &announcementsReply)
	if err != nil {
		return nil, err
	}

	return announcementsReply.Response.Data, nil
}
