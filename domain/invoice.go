package domain

type Position struct {
	Hours float64
}

type Invoice struct {
	Id        int    `json:"id"`
	Status    string `json:"status"`
	Customer  string `json:"customer"`
	Positions map[int]map[string]Position `json:"positions"`
}

func NewInvoice(customer string) *Invoice {
	return &Invoice{Customer: customer, Positions: make(map[int]map[string]Position)}
}

func (i *Invoice) AddPosition(projectId int, title string, hours float64) {
	if i.Positions[projectId] == nil {
		i.Positions[projectId] = make(map[string]Position)
	}
	i.Positions[projectId][title] = Position{Hours: hours}
}
