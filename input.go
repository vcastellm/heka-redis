package heka_redis

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/mozilla-services/heka/pipeline"
	"time"
)

type RedisPubSubInputConfig struct {
	Address     string `toml:"address"`
	Channel     string `toml:"channel"`
	DecoderName string `toml:"decoder"`
}

type RedisPubSubInput struct {
	conf *RedisPubSubInputConfig
	conn redis.Conn
}

func (rpsi *RedisPubSubInput) ConfigStruct() interface{} {
	return &RedisPubSubInputConfig{":6379", "*", ""}
}

func (rpsi *RedisPubSubInput) Init(config interface{}) error {
	rpsi.conf = config.(*RedisPubSubInputConfig)

	var err error
	rpsi.conn, err = redis.Dial("tcp", rpsi.conf.Address)
	if err != nil {
		return fmt.Errorf("connecting to - %s", err.Error())
	}

	return nil
}

func (rpsi *RedisPubSubInput) Run(ir pipeline.InputRunner, h pipeline.PluginHelper) error {
	var (
		dRunner pipeline.DecoderRunner
		decoder pipeline.Decoder
		pack    *pipeline.PipelinePack
		e       error
		ok      bool
	)
	// Get the InputRunner's chan to receive empty PipelinePacks
	packSupply := ir.InChan()

	if rpsi.conf.DecoderName != "" {
		if dRunner, ok = h.DecoderRunner(rpsi.conf.DecoderName, fmt.Sprintf("%s-%s", ir.Name(), rpsi.conf.DecoderName)); !ok {
			return fmt.Errorf("Decoder not found: %s", rpsi.conf.DecoderName)
		}
		decoder = dRunner.Decoder()
	}

	//Connect to the channel
	psc := redis.PubSubConn{Conn: rpsi.conn}
	psc.PSubscribe(rpsi.conf.Channel)

	for {
		switch n := psc.Receive().(type) {
		case redis.PMessage:
			// Grab an empty PipelinePack from the InputRunner
			pack = <-packSupply
			pack.Message.SetType("redis_pub_sub")
			pack.Message.SetLogger(n.Channel)
			pack.Message.SetPayload(string(n.Data))
			pack.Message.SetTimestamp(time.Now().UnixNano())
			var packs []*pipeline.PipelinePack
			if decoder == nil {
				packs = []*pipeline.PipelinePack{pack}
			} else {
				packs, e = decoder.Decode(pack)
			}
			if packs != nil {
				for _, p := range packs {
					ir.Inject(p)
				}
			} else {
				if e != nil {
					ir.LogError(fmt.Errorf("Couldn't parse Redis message: %s", n.Data))
				}
				pack.Recycle(nil)
			}
		case redis.Subscription:
			ir.LogMessage(fmt.Sprintf("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count))
			if n.Count == 0 {
				return errors.New("No channel to subscribe")
			}
		case error:
			ir.LogError(fmt.Errorf("error: %v\n", n))
			return n
		}
	}

	return nil
}

func (rpsi *RedisPubSubInput) Stop() {
	rpsi.conn.Close()
}

func init() {
	pipeline.RegisterPlugin("RedisPubSubInput", func() interface{} {
		return new(RedisPubSubInput)
	})
}
