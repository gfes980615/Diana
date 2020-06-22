package line

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

func GetGoogleExcelValueById(id int64) string {
	url := "https://script.google.com/macros/s/AKfycbzDtZfQHmr0YJF7F_m2ZfatU7Hu-FwTpBTwQfYXqZAv7P1JnHQ/exec?msg=" + fmt.Sprintf("%d", id)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("err:\n" + err.Error())
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read error", err)
		return ""
	}

	type Tmp struct {
		Msg interface{}
	}

	test := Tmp{}
	if err := json.Unmarshal(body, &test); err != nil {
		log.Print(err.Error())
		return ""
	}

	switch reflect.TypeOf(test.Msg).Kind() {
	case reflect.Int:
		return fmt.Sprintf("%d", test.Msg.(int))
	case reflect.Int8:
		return fmt.Sprintf("%d", test.Msg.(int8))
	case reflect.Int16:
		return fmt.Sprintf("%d", test.Msg.(int16))
	case reflect.Int32:
		return fmt.Sprintf("%d", test.Msg.(int32))
	case reflect.Int64:
		return fmt.Sprintf("%d", test.Msg.(int64))
	case reflect.String:
		return test.Msg.(string)
	case reflect.Float64:
		return fmt.Sprintf("%.f", test.Msg.(float64))
	case reflect.Float32:
		return fmt.Sprintf("%.f", test.Msg.(float32))
	default:
		fmt.Println(reflect.TypeOf(test.Msg).Kind())
		return "unknow type"
	}

	return "unexcept error"
}
