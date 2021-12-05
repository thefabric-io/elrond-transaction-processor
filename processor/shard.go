package processor

import (
	"fmt"
	"strings"
)

const (
	ShardMetachain = 4294967295
)

type Shard int

func (s Shard) Equals(s2 Shard) bool {
	return s == s2
}

func (s Shard) Name() string {
	if s.Equals(ShardMetachain) {
		return "Metachain"
	}

	return fmt.Sprintf("Shard %d", s)
}

type Shards []Shard

func (ss Shards) String() string {
	b := strings.Builder{}

	for i, shard := range ss {
		b.WriteString(shard.Name())
		if i != len(ss)-1 {
			b.WriteString(", ")
		}
	}

	return b.String()
}
