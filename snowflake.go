package snowflake

import (
	"errors"
	"time"
)

type Snowflake struct {
	NodeId         int64
	Epoch          int64
	NodeBits       uint
	SequenceBits   uint
	NodeMax        int64
	SequenceMax    int64
	NodeShift      uint
	TimestampShift uint
	Sequence       int64
	LastTimestamp  int64
}

type Config struct {
	NodeId       int64
	Epoch        int64
	NodeBits     uint
	SequenceBits uint
}

var defaultConfig = Config{
	NodeId:       0,
	Epoch:        1288834974657,
	NodeBits:     10,
	SequenceBits: 12,
}

func Node(config Config) (*Snowflake, error) {
	if config.Epoch == 0 {
		config.Epoch = defaultConfig.Epoch
	}
	if config.NodeBits == 0 {
		config.NodeBits = defaultConfig.NodeBits
	}
	if config.SequenceBits == 0 {
		config.SequenceBits = defaultConfig.SequenceBits
	}

	if config.NodeBits+config.SequenceBits != 22 {
		return nil, errors.New("sum of nodeBits and sequenceBits must be 22")
	}

	nodeMax := int64(-1) ^ (int64(-1) << config.NodeBits)
	sequenceMax := int64(-1) ^ (int64(-1) << config.SequenceBits)

	if config.NodeId > nodeMax {
		return nil, errors.New("node ID can't be greater than nodeMax")
	}

	return &Snowflake{
		NodeId:         config.NodeId,
		Epoch:          config.Epoch,
		NodeBits:       config.NodeBits,
		SequenceBits:   config.SequenceBits,
		NodeMax:        nodeMax,
		SequenceMax:    sequenceMax,
		NodeShift:      config.SequenceBits,
		TimestampShift: config.SequenceBits + config.NodeBits,
		LastTimestamp:  -1,
	}, nil
}

func (s *Snowflake) Now() int64 {
	return time.Now().UnixMilli()
}

func (s *Snowflake) Generate() (int64, error) {
	now := s.Now()
	if now < s.LastTimestamp {
		return 0, errors.New("invalid system clock")
	}

	if s.LastTimestamp == now {
		s.Sequence = (s.Sequence + 1) & s.SequenceMax
		if s.Sequence == 0 {
			for now <= s.LastTimestamp {
				now = s.Now()
			}
		}
	} else {
		s.Sequence = 0
	}

	s.LastTimestamp = now

	if now < s.Epoch {
		return 0, errors.New("time is moving backwards. Refusing to generate id")
	}

	return ((now - s.Epoch) << s.TimestampShift) | (s.NodeId << s.NodeShift) | s.Sequence, nil
}

func (s *Snowflake) GetTimeFromId(id int64) int64 {
	return (id >> s.TimestampShift) + s.Epoch
}

func (s *Snowflake) GetNodeFromId(id int64) int64 {
	return (id & ((int64(-1) ^ (int64(-1) << s.NodeShift)) << s.SequenceBits)) >> s.SequenceBits
}

func (s *Snowflake) GetSequenceFromId(id int64) int64 {
	return id & (int64(-1) ^ (int64(-1) << s.SequenceBits))
}
