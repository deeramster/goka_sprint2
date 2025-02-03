package models

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type BlockCommand struct {
	User      string `json:"user"`
	BlockUser string `json:"block"`
}

type BlockedUsers struct {
	Users []string `json:"users"`
}
