package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// IndicatorSignal holds the schema definition for the IndicatorSignal entity.
type UserSetting struct {
	ent.Schema
}

func (UserSetting) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the IndicatorSignal.
func (UserSetting) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").NotEmpty(),
		field.String("type").NotEmpty(),
		field.String("nid").NotEmpty(),
		field.Bool("enabled").Default(true),
	}
}

func (UserSetting) Indexes() []ent.Index {
	return []ent.Index{
		// unique index.
		index.Fields("user_id", "nid", "type").
			Unique(),
	}
}

// Edges of the IndicatorSignal.
func (UserSetting) Edges() []ent.Edge {
	return nil
}
