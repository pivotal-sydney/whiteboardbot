package model

type Standup struct {
	Id int				`json:"id"`
	TimeZone string		`json:"time_zone_name_iana"`
	Title string		`json:"title"`
}