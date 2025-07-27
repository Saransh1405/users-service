package models

type Notification struct {
	Type      string `json:"type" bson:"type"`
	UserID    string `json:"userId" bson:"userId"`
	Message   string `json:"message" bson:"message"`
	Data      any    `json:"data,omitempty" bson:"data,omitempty"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
}
