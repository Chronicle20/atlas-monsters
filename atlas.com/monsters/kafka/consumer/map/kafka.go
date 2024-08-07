package _map

import "atlas-monsters/tenant"

const (
	EnvEventTopicMapStatus                = "EVENT_TOPIC_MAP_STATUS"
	EventTopicMapStatusTypeCharacterEnter = "CHARACTER_ENTER"
	EventTopicMapStatusTypeCharacterExit  = "CHARACTER_EXIT"
)

type statusEvent[E any] struct {
	Tenant    tenant.Model `json:"tenant"`
	WorldId   byte         `json:"worldId"`
	ChannelId byte         `json:"channelId"`
	MapId     uint32       `json:"mapId"`
	Type      string       `json:"type"`
	Body      E            `json:"body"`
}

type characterEnter struct {
	CharacterId uint32 `json:"characterId"`
}

type characterExit struct {
	CharacterId uint32 `json:"characterId"`
}