package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID            uuid.UUID       `json:"id" pg:",type:uuid"`
	Name          string          `json:"name" pg:",notnull"`
	Password      string          `json:"password"`
	Avatar        []byte          `json:"avatar"`
	RegisterIP    string          `json:"register_ip"`
	RegisterDate  time.Time       `json:"register_date" pg:",default:now()"`
	LastLoginIP   string          `json:"last_login_ip"`
	LastLoginDate time.Time       `json:"last_login_date" pg:",default:now()"`
	LanguageCode  string          `json:"language_code" pg:",default:'en'"` // eg. "en"
	RankLevelID   string          `json:"rank_level_id"`                    // UUID of RankLevel struct containing
	Tokens        []*Token        `json:"tokens" pg:"rel:has-many,join_fk:token"`
	Notifications []*Notification `json:"notifications" pg:"rel:has-many,join_fk:user"`
	Favourites    []*Favourite    `json:"favourites" pg:"rel:has-many,join_fk:user"`
} // booleans of allowed functions

type Notification struct {
	ID      uuid.UUID `json:"id" sql:",pk,type:uuid"`
	User    *User     `json:"user" pg:"rel:has-one,notnull"`        // To which user this notification belongs
	WasRead bool      `json:"was_read" pg:",notnull,default:false"` // If user already read this notification
	Message string    `json:"message" pg:",notnull"`                // Message contents
}

type Favourite struct {
	ID     int       `json:"id"`
	UserID uuid.UUID `json:"user_id" pg:",type:uuid"`
	SongID string    `json:"song_id"`
	User   *User     `json:"user" pg:",rel:belongs-to"`
}

// API tokens
type Token struct {
	Token string `json:"token" pg:",pk"` // auto-generated random unique string
	Owner *User  `json:"owner" pg:"rel:belongs-to,join_fk:tokens"`
}
