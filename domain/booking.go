package domain

type Booking struct {
	Id         int
	Day        int     `json:"day"`
	Hours      float32 `json:"hours"`
	InvoiceId  int     `json:"invoiceId"`
	ProjectId  int     `json:"projectId"`
	ActivityId int     `json:"activityId"`
}
