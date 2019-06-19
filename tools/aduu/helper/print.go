package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

func PrintAny(i interface{}) {
	d, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s:\n%s\n", reflect.TypeOf(i).String(), string(d))
}
