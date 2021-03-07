package models

import uuid "github.com/satori/go.uuid"

type User struct {
	ID            uuid.NullUUID `json:"id" sql:",pk,type:uuid default uuid_generate_v4()"`
	Name          string        `json:"name"`
	Password      string        `json:"password"`
	Avatar        []byte        `json:"avatar"`
	RegisterIP    string        `json:"register_ip"`
	RegisterDate  int64         `json:"register_date"`
	LastLoginIP   string        `json:"last_login_ip"`
	LastLoginDate int64         `json:"last_login_date"`
	LanguageCode  string        `json:"language_code"` // eg. "en"
	RankLevelID   string        `json:"rank_level_id"` // UUID of RankLevel struct containing
} // booleans of allowed functions

type Notification struct {
	ID      uuid.NullUUID `json:"id" sql:",pk,type:uuid default uuid_generate_v4()"`
	User    string        `json:"user" pg:"rel:has-one"`        // To which user this notification belongs
	WasRead bool          `json:"was_read" sql:"default:false"` // If user already read this notification
	Message string        `json:"message"`                      // Message contents
}

type Favourites struct {
	User   *User  `json:"user_id" pg:"rel:has-one"`
	SongID string `json:"song_id"`
}

// API tokens
type Tokens struct {
	Token string `json:"token"` // auto-generated random unique string
	Owner *User  `json:"owner" pg:"rel:belongs-to"`
}