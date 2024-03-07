package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

type UnknownPlant struct {
	FlowerType string
	LeafType   string
	Color      int `color_scheme:"rgb"`
}

type AnotherUnknownPlant struct {
	FlowerColor int
	LeafType    string
	Height      int `unit:"inches"`
}

func describePlant(anyPlant interface{}) {
	t := reflect.TypeOf(anyPlant)
	if t.Kind() == reflect.Struct {
		v := reflect.ValueOf(anyPlant)
		printPlant(&t, &v)
	} else {
		log.Fatal(errors.New("неверный тип для describePlant"))
	}
}

func printPlant(t *reflect.Type, v *reflect.Value) {
	for i := 0; i < (*t).NumField(); i++ {
		field := (*t).Field(i)
		value := (*v).Field(i)

		// Получаем имя поля
		fieldName := field.Name

		// Получаем значение поля как строку
		var fieldValue string
		switch value.Interface().(type) {
		case int, int8, int16, int32, int64:
			fieldValue = fmt.Sprintf("%d", value.Interface())
		default:
			fieldValue = fmt.Sprintf("%v", value.Interface())
		}

		// Получаем теги для поля (если есть)
		tag := field.Tag.Get("unit")
		if tag == "" {
			tag = field.Tag.Get("color_scheme")
		}
		if tag != "" {
			Tag := strings.Split(string(field.Tag), ":")[0]
			fieldName = fmt.Sprintf("%s(%s=%s)", fieldName, Tag, tag)
		}

		fmt.Printf("%s:%s\n", fieldName, fieldValue)
	}
}

func main() {
	UP := UnknownPlant{
		"pion",
		"lanceolate",
		333,
	}
	AP := AnotherUnknownPlant{
		10,
		"lanceolate",
		15,
	}

	fmt.Println(color.HiMagentaString("UnknownPlant"))
	describePlant(UP)
	fmt.Println(color.HiGreenString("\nAnotherUnknownPlant"))
	describePlant(AP)
}
