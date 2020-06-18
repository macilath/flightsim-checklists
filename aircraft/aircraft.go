package aircraft

type LiteChecklist struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Checklist struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Items      []string `json:"items"`
	AircraftID int      `json:"aircraft_id"`
}

type Aircraft struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	ShortName  string          `json:"alias"`
	Checklists []LiteChecklist `json:"checklists"`
}
