package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type NotificationData struct {
	Amount   string `json:"amount,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	Code     string `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	Tx       string `json:"tx,omitempty"`
	Name     string `json:"name,omitempty"`
	Referral string `json:"referral,omitempty"`
	Coupon   string `json:"coupon,omitempty"`
}

// Notification holds the schema definition for the Notification entity.
type Notification struct {
	ent.Schema
}

// Fields of the Notification.
func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("user_id", uuid.UUID{}),
		field.String("title_key"),
		field.JSON("data", NotificationData{}),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
		field.Bool("read").Default(false),
	}
}

// Edges of the Notification.
func (Notification) Edges() []ent.Edge {
	return nil
}
