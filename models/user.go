package models

import uuid "github.com/google/uuid"

type User struct {
	ID            uuid.UUID       `json:"id" sql:",pk,notnull,type:uuid default uuid_generate_v4()"`
	Name          string          `json:"name" pg:",notnull"`
	Password      string          `json:"password"`
	Avatar        []byte          `json:"avatar"`
	RegisterIP    string          `json:"register_ip"`
	RegisterDate  int64           `json:"register_date" sql:"default:now()"`
	LastLoginIP   string          `json:"last_login_ip"`
	LastLoginDate int64           `json:"last_login_date" sql:"default:now()"`
	LanguageCode  string          `json:"language_code" sql:"default:en"` // eg. "en"
	RankLevelID   string          `json:"rank_level_id"`                  // UUID of RankLevel struct containing
	Tokens        []*Token        `json:"tokens" pg:"rel:has-many,join_fk:token"`
	Notifications []*Notification `json:"notifications" pg:"rel:has-many,join_fk:user"`
	Favourites    []*Favourite    `json:"favourites" pg:"rel:has-many,join_fk:user"`
} // booleans of allowed functions

type Notification struct {
	ID      uuid.UUID `json:"id" sql:",pk,notnull,type:uuid default uuid_generate_v4()"`
	User    *User     `json:"user" pg:"rel:has-one,notnull"` // To which user this notification belongs
	WasRead bool      `json:"was_read" sql:"default:false"`  // If user already read this notification
	Message string    `json:"message" pg:",notnull"`         // Message contents
}

type Favourite struct {
	ID     int       `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	SongID string    `json:"song_id"`
	User   *User     `json:"user" pg:"rel:belongs-to"`
}

// API tokens
type Token struct {
	Token string `json:"token" pg:",pk"` // auto-generated random unique string
	Owner *User  `json:"owner" pg:"rel:belongs-to,join_fk:tokens"`
}
