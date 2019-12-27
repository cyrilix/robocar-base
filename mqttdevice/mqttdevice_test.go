package mqttdevice

import (
	"github.com/cyrilix/robocar-base/types"
	"github.com/cyrilix/robocar-base/testtools/docker"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"testing"
)

func TestIntegration(t *testing.T) {

	ctx, mqttC, mqttUri := docker.MqttContainer(t)
	defer mqttC.Terminate(ctx)

	t.Run("ConnectAndClose", func(t *testing.T) {
		t.Logf("Mqtt connection %s ready", mqttUri)

		p := pahoMqttPubSub{Uri: mqttUri, ClientId: "TestMqtt", Username: "guest", Password: "guest"}
		p.Connect()
		p.Close()
	})
	t.Run("Publish", func(t *testing.T) {
		options := mqtt.NewClientOptions().AddBroker(mqttUri)
		options.SetUsername("guest")
		options.SetPassword("guest")

		client := mqtt.NewClient(options)
		token := client.Connect()
		defer client.Disconnect(100)
		token.Wait()
		if token.Error() != nil {
			t.Fatalf("unable to connect to mqtt broker: %v\n", token.Error())
		}

		c := make(chan string)
		defer close(c)
		client.Subscribe("test/publish", 0, func(client mqtt.Client, message mqtt.Message) {
			c <- string(message.Payload())
		}).Wait()

		p := pahoMqttPubSub{Uri: mqttUri, ClientId: "TestMqtt", Username: "guest", Password: "guest"}
		p.Connect()
		defer p.Close()

		p.Publish("test/publish", []byte("Test1234"))
		result := <-c
		if result != "Test1234" {
			t.Fatalf("bad message: %v\n", result)
		}

	})
}

func TestNewMqttValue(t *testing.T) {
	cases := []struct {
		value    interface{}
		expected MqttValue
	}{
		{"text", []byte("text")},
		{float32(2.0123), []byte("2.01")},
		{3.12345, []byte("3.12")},
		{12, []byte("12")},
		{true, []byte("ON")},
		{false, []byte("OFF")},
		{MqttValue("13"), []byte("13")},
		{[]byte("test bytes"), []byte("test bytes")},

		{struct {
			Content string
		}{"other"}, []byte(`{"Content":"other"}`)},
	}

	for _, c := range cases {
		val := NewMqttValue(c.value)
		if string(val) != string(c.expected) {
			t.Errorf("NewMqttValue(%v): %v, wants %v", c.value, string(val), string(c.expected))
		}
	}
}

func TestMqttValue_BoolValue(t *testing.T) {
	cases := []struct {
		value    MqttValue
		expected bool
	}{
		{NewMqttValue("ON"), true},
		{NewMqttValue("OFF"), false},
	}
	for _, c := range cases {
		val, err := c.value.BoolValue()
		if err != nil {
			t.Errorf("unexpected conversion error: %v", err)
		}
		if c.expected != val {
			t.Errorf("MqttValue.BoolValue(): %v, wants %v", val, c.expected)
		}
	}
}

func TestMqttValue_ByteSliceValue(t *testing.T) {
	cases := []struct {
		value    MqttValue
		expected []byte
	}{
		{NewMqttValue([]byte("content")), []byte("content")},
	}
	for _, c := range cases {
		val, err := c.value.ByteSliceValue()
		if err != nil {
			t.Errorf("unexpected conversion error: %v", err)
		}
		if string(c.expected) != string(val) {
			t.Errorf("MqttValue.BoolValue(): %v, wants %v", val, c.expected)
		}
	}
}

func TestMqttValue_Float32Value(t *testing.T) {
	cases := []struct {
		value    MqttValue
		expected float32
	}{
		{NewMqttValue("32.0123"), float32(32.0123)},
		{NewMqttValue("33"), float32(33.)},
	}
	for _, c := range cases {
		val, err := c.value.Float32Value()
		if err != nil {
			t.Errorf("unexpected conversion error: %v", err)
		}
		if c.expected != val {
			t.Errorf("MqttValue.BoolValue(): %v, wants %v", val, c.expected)
		}
	}
}

func TestMqttValue_Float64Value(t *testing.T) {
	cases := []struct {
		value    MqttValue
		expected float64
	}{
		{NewMqttValue("32.0123"), 32.0123},
		{NewMqttValue("33"), 33.},
	}
	for _, c := range cases {
		val, err := c.value.Float64Value()
		if err != nil {
			t.Errorf("unexpected conversion error: %v", err)
		}
		if c.expected != val {
			t.Errorf("MqttValue.BoolValue(): %v, wants %v", val, c.expected)
		}
	}
}
func TestMqttValue_IntValue(t *testing.T) {
	cases := []struct {
		value    MqttValue
		expected int
	}{
		{NewMqttValue("1"), 1},
		{NewMqttValue("-10"), -10},
	}
	for _, c := range cases {
		val, err := c.value.IntValue()
		if err != nil {
			t.Errorf("unexpected conversion error: %v", err)
		}
		if c.expected != val {
			t.Errorf("MqttValue.BoolValue(): %v, wants %v", val, c.expected)
		}
	}
}
func TestMqttValue_StringValue(t *testing.T) {
	cases := []struct {
		value    MqttValue
		expected string
	}{
		{NewMqttValue("ON"), "ON"},
		{NewMqttValue("OFF"), "OFF"},
	}
	for _, c := range cases {
		val, err := c.value.StringValue()
		if err != nil {
			t.Errorf("unexpected conversion error: %v", err)
		}
		if c.expected != val {
			t.Errorf("MqttValue.BoolValue(): %v, wants %v", val, c.expected)
		}
	}
}

func TestMqttValue_DriveModeValue(t *testing.T) {
	cases := []struct {
		value    MqttValue
		expected types.DriveMode
	}{
		{NewMqttValue(types.DriveModeUser), types.DriveModeUser},
		{NewMqttValue(types.DriveModePilot), types.DriveModePilot},
		{NewMqttValue(types.DriveModeInvalid), types.DriveModeInvalid},
	}
	for _, c := range cases {
		val, err := c.value.DriveModeValue()
		if err != nil {
			t.Errorf("unexpected conversion error: %v", err)
		}
		if c.expected != val {
			t.Errorf("MqttValue.DriveMode(): %v, wants %v", val, c.expected)
		}
	}
}
