package schema

type User struct {
	Username      string    `firestore:"username, omitempty"`
	Password      string    `firestore:"password, omitempty"`
	IsGroupMember bool      `firestore:"isGroupMember, omitempty"`
	Messages      []Message `firestore:"messages, omitempty"`
}

type Message struct {
	Author  string `firestore:"author, omitempty"`
	Content string `firestore:"content, omitempty"`
}
