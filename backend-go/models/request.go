package models

type Shift struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Role     string `json:"role"`
	Location string `json:"location"`
}

type Request struct {
	ID      int    `json:"id"`
	ShiftID int    `json:"shift_id"`
	UserID  int    `json:"user_id"`
	Status  string `json:"status"`
}
