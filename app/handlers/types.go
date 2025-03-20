package handlers

type SkillOption struct {
	Name string `string:"name"`;
	Level int32 `int32:"level"`;
}
type Skill struct {
	Level int32 `int32:"level"`;
	Has_options bool `bool:"has_options"`;
	Stat_type string `string:"stat_type"`;
	Options map[string]SkillOption `map[string]SkillOption:"options"`;
}

type Effect struct {
	Name string `string:"name"`;
	Description string `string:"description"`;
	Category string `string:"category"`;
	Alt string `string:"alt"`;
}

type Character struct {
	Id *string `string:"id"`;
	Name string `string:"name"`;
	Role string `string:"role"`;
	Stats map[string]int32 `map[string]int32:"stats"`;
	HP int32 `int32:"HP"`;
	Humanity int32 `int32:"humanity"`;
	CurrentSkills map[string]Skill `map[string]Skill:"currentSkills"`;
	CurrentEffects map[string]Effect `map[string]Effect:"currentEffects"`;
}