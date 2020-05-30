package aircraft

type Checklist struct {
	Title string
	Items []string `json:"items"`
}

type Aircraft struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	ShortName  string   `json:"alias"`
	Checklists []string `json:"checklists"`
}
