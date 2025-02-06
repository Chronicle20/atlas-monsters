package _map

import (
	consumer2 "atlas-monsters/kafka/consumer"
	"atlas-monsters/monster"
	"context"
	"github.com/Chronicle20/atlas-kafka/consumer"
	"github.com/Chronicle20/atlas-kafka/handler"
	"github.com/Chronicle20/atlas-kafka/message"
	"github.com/Chronicle20/atlas-kafka/topic"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/sirupsen/logrus"
)

func InitConsumers(l logrus.FieldLogger) func(func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
	return func(rf func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
		return func(consumerGroupId string) {
			rf(consumer2.NewConfig(l)("map_status_event")(EnvEventTopicMapStatus)(consumerGroupId), consumer.SetHeaderParsers(consumer.SpanHeaderParser, consumer.TenantHeaderParser))
		}
	}
}

func InitHandlers(l logrus.FieldLogger) func(rf func(topic string, handler handler.Handler) (string, error)) {
	return func(rf func(topic string, handler handler.Handler) (string, error)) {
		var t string
		t, _ = topic.EnvProvider(l)(EnvEventTopicMapStatus)()
		_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventCharacterEnter)))
		_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventCharacterExit)))
	}
}

func handleStatusEventCharacterEnter(l logrus.FieldLogger, ctx context.Context, event statusEvent[characterEnter]) {
	if event.Type != EventTopicMapStatusTypeCharacterEnter {
		return
	}

	provider := monster.NotControlledInMapProvider(ctx)(event.WorldId, event.ChannelId, event.MapId)
	_ = model.ForEachSlice(provider, monster.FindNextController(l)(ctx), model.ParallelExecute())
}

func handleStatusEventCharacterExit(l logrus.FieldLogger, ctx context.Context, event statusEvent[characterExit]) {
	if event.Type != EventTopicMapStatusTypeCharacterExit {
		return
	}

	provider := monster.ControlledByCharacterInMapProvider(ctx)(event.WorldId, event.ChannelId, event.MapId, event.Body.CharacterId)
	_ = model.ForEachSlice(provider, monster.StopControl(l)(ctx), model.ParallelExecute())
	_ = model.ForEachSlice(provider, monster.FindNextController(l)(ctx))
}
