package communication

type Message struct {
	Type         string `json:"type"`
	UserID       string `json:"id"`
	Body         []byte `json:"body"`
	Name         string `json:"name"`
	Extension    string `json:"extension"`
	NewName      string `json:"new_name"`
	NewExtension string `json:"new_extension"`
}
