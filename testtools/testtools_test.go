package testtools

import (
	"github.com/cyrilix/robocar-base/mqttdevice"
	"testing"
)

func TestFakePublisher_Publish(t *testing.T) {
	p := NewFakePublisher()

	cases := []struct {
		topic string
		topicPublished string
		value mqttdevice.MqttValue
		expected string
	}{
		{"test/topic1", "test/topic1", mqttdevice.NewMqttValue(1) , "1" },
		{"test/topic2", "test/invalid", mqttdevice.NewMqttValue(1) , "" },
	}

	for _, c := range cases{
		p.Publish(c.topic, c.value)
		val := p.PublishedEvent(c.topicPublished)
		if v, _ := val.StringValue(); v != c.expected {
			t.Errorf("FakePublisher.Publish(%v, %v): %v, wants %v", c.topic, string(c.value), v, c.expected)
		}

	}
}
