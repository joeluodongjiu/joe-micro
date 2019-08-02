package toolfunc

import (
	"github.com/bwmarrin/snowflake"
	"log"
	"math/rand"
)

var node  *snowflake.Node

func init()  {
	var  err error
	node, err = snowflake.NewNode(int64(rand.Intn(1022)+1))
	if err != nil {
		log.Println(err)
		return
	}
}


func   GenerateUUID()  string{
	// Generate a snowflake ID.
	id := node.Generate()
	return id.String()
}