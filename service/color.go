package service

import (
	"fmt"
	"github.com/tidwall/pretty"
)

var ColorService colorService

type colorService struct {
	Enable bool
}

func (c colorService) Print(json []byte) {
	if c.Enable {
		result := pretty.Color(json, nil)
		fmt.Println(string(result))
	} else {
		fmt.Println(string(json))
	}
}
