package aircraft

type Aircraft struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"alias"`
}
