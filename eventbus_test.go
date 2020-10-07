package eventbus

import (
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	bus := NewEventBus()
	m := &Model1{}
	// Register(bus, m)
	ev := NewEvent("create", &Data1{
		Name: "Jeek",
		Age:  1,
	}, time.Now().Unix())
	//AddHandler = func(string, h func(context.Context, Event)) error
	// m.Event1  = func(string, h func(context.Context, Event1)) error
	// Event是接口，Event1 已经实现了Event的接口
	bus.AddHandler(m, "create")
	bus.HandleEvent(ev)
}

type Model1 struct {
}

func (m *Model1) HandleEvent(ev Event) error {
	return nil
}

type Data1 struct {
	Name string
	Age  int
}

// var _ = Event(&Event1{})

func init() {
	RegisterData(func() Data {
		return &Data1{}
	}, "create", "udpate")
}

func (ev *Data1) DataType() DataType {
	return "data1"
}

func (ev *Data1) Id() DataId {
	return DataId(ev.Name)
}
