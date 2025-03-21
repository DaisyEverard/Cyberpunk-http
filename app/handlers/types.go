package handlers

import "go.mongodb.org/mongo-driver/bson/primitive"

type SkillOption struct {
	Name string `bson:"name"`;
	Level int32 `bson:"level"`;
}
type Skill struct {
	Level int32 `bson:"level"`;
	Has_options bool `bool:"has_options"`;
	Stat_type string `bson:"stat_type"`;
	Options map[string]SkillOption `bson:"options"`;
}

type Effect struct {
	Name string `bson:"name"`;
	Description string `bson:"description"`;
	Category string `bson:"category"`;
	Alt string `bson:"alt"`;
}

type CharacterWithID struct {
	Id primitive.ObjectID `bson:"id"`;
	Name string `bson:"name"`;
	Role string `bson:"role"`;
	Stats map[string]int32 `bson:"stats"`;
	HP int32 `bson:"hp"`;
	Humanity int32 `bson:"humanity"`;
	CurrentSkills map[string]Skill `bson:"currentSkills"`;
	CurrentEffects map[string]Effect `bson:"currentEffects"`;
}

type CharacterWithoutID struct {
	Name string `bson:"name"`;
	Role string `bson:"role"`;
	Stats map[string]int32 `bson:"stats"`;
	HP int32 `bson:"hp"`;
	Humanity int32 `bson:"humanity"`;
	CurrentSkills map[string]Skill `bson:"currentSkills"`;
	CurrentEffects map[string]Effect `bson:"currentEffects"`;
}