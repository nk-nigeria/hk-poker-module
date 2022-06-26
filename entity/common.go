package entity

import (
	"github.com/bwmarrin/snowflake"
)

const (
	ModuleName      = "chinese-poker"
	MaxPresenceCard = 13
)

var SnowlakeNode, _ = snowflake.NewNode(1)
