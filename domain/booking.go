package domain

type Booking struct {
	Id         int     `json:"id"`
	Day        int     `json:"day"`
	Hours      float32 `json:"hours"`
	InvoiceId  int     `json:"invoiceId"`  // belongs to invoice
	ProjectId  int     `json:"projectId"`  // belongs to project
	ActivityId int     `json:"activityId"` // belongs to activity
}
