package utils

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

//GenerateID -
func GenerateID() int64 {
	if node == nil {
		node, _ = snowflake.NewNode(1)
	}
	id := node.Generate()

	return id.Int64()
}

//GeerateIDBase64 -
func GeerateIDBase64() string {
	if node == nil {
		node, _ = snowflake.NewNode(1)
	}
	id := node.Generate()

	return id.Base64()
}

//GeerateIDBase36 -
func GeerateIDBase36() string {
	if node == nil {
		node, _ = snowflake.NewNode(1)
	}
	id := node.Generate()

	return id.Base36()
}

// GenerateTraceID -
func GenerateTraceID() string {
	if node == nil {
		node, _ = snowflake.NewNode(1)
	}
	id := node.Generate()

	return id.Base36()
}
