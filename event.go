package eventbus

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

func NewEmptyEvent() Event {
	return &event{}
}

func NewEvent(eventType EventType, data Data, created int64) Event {
	return &event{
		EventType_: eventType,
		Data_:      data,
		Created_:   created,
	}
}

// ErrEventNotRegistered is when no command factory was registered.
var ErrEventNotRegistered = errors.New("Event not registered")

var eventDatas = make(map[EventType]func() Data)

type EventType string

func (e EventType) String() string {
	return string(e)
}

type Event interface {
	// EventType returns the type of the event.
	EventType() EventType
	// The data attached to the event.
	Data() Data
	// Timestamp of when the event was created.
	Created() int64
	Bytes() []byte
	FromBytes(bs []byte)
}

type event struct {
	// Event的数据
	Data_ Data
	// 返回Event的名称
	EventType_ EventType
	// 时间戳
	Created_ int64
}

func (e *event) Data() Data {
	return e.Data_
}

func (e *event) Created() int64 {
	return e.Created_
}

func (e *event) EventType() EventType {
	return e.EventType_
}

func (e *event) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(e); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (e *event) FromBytes(bs []byte) {
	dec := gob.NewDecoder(bytes.NewReader(bs))
	if err := dec.Decode(e); err != nil {
		panic(err)
	}
}

type DataType string
type DataId string

func (t DataType) String() string {
	return string(t)
}
func (d DataId) String() string {
	return string(d)
}

// 以后可能需要把存储的数据和事件的数据分开
type Data interface {
	Id() DataId
	DataType() DataType
}

// 根据事件，创建消息数据
func RegisterData(factory func() Data, evTypes ...EventType) {
	// TODO: Explore the use of reflect/gob for creating concrete types without
	// a factory func.

	// Check that the created command matches the type registered.
	data := factory()
	if data == nil {
		panic("eventbus: created eventData is nil")
	}
	gob.Register(data)

	for _, evType := range evTypes {
		if evType == EventType("") {
			panic("eventbus: attempt to register empty event type")
		}

		if _, ok := eventDatas[evType]; ok {
			panic(fmt.Sprintf("eventbus: registering duplicate types for %q", evType))
		}
		eventDatas[evType] = factory
	}
}

func UnregisterData(evType EventType) {
	if evType == EventType("") {
		panic("eventbus: attempt to unregister empty event type")
	}

	if _, ok := eventDatas[evType]; !ok {
		panic(fmt.Sprintf("eventbus: unregister of non-registered type %q", evType))
	}
	delete(eventDatas, evType)
}

// 根据事件，创建数据
func CreateEventData(eventT EventType) (Data, error) {
	if factory, ok := eventDatas[eventT]; ok {
		return factory(), nil
	}
	return nil, ErrEventNotRegistered
}
