package strategy

import (
	pb "github.com/golang/protobuf/proto"
	"github.com/mcorbin/riemann-relay/pkg/client"
	"github.com/riemann/riemann-go-client"
	"github.com/riemann/riemann-go-client/proto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBroadcastStrategyTest(t *testing.T) {
	sink1 := make(chan *proto.Msg)
	sink2 := make(chan *proto.Msg)
	strategy := BroadcastStrategy{
		Clients: client.NewFixtureClients([]chan *proto.Msg{sink1, sink2}, true),
	}
	events := []riemanngo.Event{
		{
			Host:    "baz",
			Service: "foobar",
			Metric:  10,
			Time:    time.Unix(100, 0),
		},
	}
	strategy.Send(&events)
	msg1 := <-sink1
	msg2 := <-sink2
	protoEvents := []*proto.Event{
		&proto.Event{
			Host:         pb.String("baz"),
			Time:         pb.Int64(100),
			MetricSint64: pb.Int64(10),
			TimeMicros:   pb.Int64(100000000),
			Service:      pb.String("foobar"),
		},
	}
	assert.Equal(t, len(msg1.Events), 1)
	assert.Equal(t, msg1, &proto.Msg{Events: protoEvents})
	assert.Equal(t, len(msg2.Events), 1)
	assert.Equal(t, msg2, &proto.Msg{Events: protoEvents})

	events = []riemanngo.Event{
		{
			Host:    "baz",
			Service: "foobar",
			Metric:  10,
			Time:    time.Unix(100, 0),
		},
		{
			Host:    "foo",
			Service: "bar",
			Metric:  100,
			Time:    time.Unix(100, 0),
		},
	}
	strategy.Send(&events)
	msg1 = <-sink1
	msg2 = <-sink2
	protoEvents = []*proto.Event{
		&proto.Event{
			Host:         pb.String("baz"),
			Time:         pb.Int64(100),
			MetricSint64: pb.Int64(10),
			TimeMicros:   pb.Int64(100000000),
			Service:      pb.String("foobar"),
		},
		&proto.Event{
			Host:         pb.String("foo"),
			Time:         pb.Int64(100),
			MetricSint64: pb.Int64(100),
			TimeMicros:   pb.Int64(100000000),
			Service:      pb.String("bar"),
		},
	}
	assert.Equal(t, len(msg1.Events), 2)
	assert.Equal(t, msg1, &proto.Msg{Events: protoEvents})
	assert.Equal(t, len(msg2.Events), 2)
	assert.Equal(t, msg2, &proto.Msg{Events: protoEvents})
}

func TestBroadcastStrategyNotConnected(t *testing.T) {
	sink1 := make(chan *proto.Msg)
	strategy := BroadcastStrategy{
		Clients: client.NewFixtureClients([]chan *proto.Msg{sink1}, false),
	}
	events := []riemanngo.Event{
		{
			Host:    "baz",
			Service: "foobar",
			Metric:  10,
			Time:    time.Unix(100, 0),
		},
	}
	strategy.Send(&events)
	select {
	case msg := <-sink1:
		t.Log("error : received ", msg)
		t.FailNow()
	default:
		// the channel should be empty
	}
}
