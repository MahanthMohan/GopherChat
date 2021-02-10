package schema

type User struct {
	Username      string   `firestore:"username,omitempty"`
	Password      string   `firestore:"password,omitempty"`
	IsGroupMember bool     `firestore:"isGroupMember,omitempty"`
	Messages      []string `firestore:"messages"`
}
