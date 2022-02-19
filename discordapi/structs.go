package discordapi

/*type Guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}*/

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type uint8  `json:"type"`
	LMI  string `json:"last_message_id"`
}

type Message struct {
	ID        string           `json:"id"`
	Author    *User            `json:"author"`
	Content   string           `json:"content"`
	Files     []Attachment     `json:"attachments"`
	Reference MessageReference `json:"referenced_message"`
}

type MessageSender struct {
	Content string `json:"content"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bot      bool   `json:"bot"`
}

type MessageReference struct {
	ID      string       `json:"id"`
	Author  *User        `json:"author"`
	Content string       `json:"content"`
	Files   []Attachment `json:"attachments"`
}
