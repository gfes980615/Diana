package injection

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

func init() {
	injector = &Injector{
		objs: make(map[string]reflect.Value, 0),
	}
}

var injector *Injector

const INJECTION = "injection"

type Injector struct {
	rwLock sync.RWMutex
	objs   map[string]reflect.Value
}

func Register(name string, v interface{}) error {
	if injector == nil {
		return errors.New("fail with nil injector")
	}

	ok, err := injector.has(name)
	if err != nil {
		return err
	}

	if ok {
		panic("Obj " + name + "is exist.")
		return errors.New("AutoRegister fail, the entry key object already exist")
	}

	return injector.put(name, v)
}

func AutoRegister(value interface{}) error {
	if injector == nil {
		return errors.New("fail with nil injector")
	}

	name := getStructName(value)

	return Register(name, value)
}

func AutoRegisterMultiple(value interface{}, keys []string) error {
	name := getStructName(value)
	var err error
	for _, key := range keys {
		err = Register(key+"-"+name, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitInject() error {
	for _, v := range injector.objs {
		value := v
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		if value.Kind() == reflect.Struct {
			for i := 0; i < value.NumField(); i++ {
				name := value.Type().Field(i).Tag.Get(INJECTION)
				temp, ok := injector.objs[name]
				if ok && reflect.TypeOf(temp).Kind() == reflect.Struct {
					field := value.Field(i)
					if field.CanSet() {
						field.Set(temp)
					} else {
						field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
						field.Set(temp)
					}
				} else if name != "" {
					return errors.New("injection error not find " + name)
				}
			}
		}
	}
	return nil
}

func getStructName(value interface{}) string {
	typ := reflect.TypeOf(value).String()
	sub := "."
	tSlice := strings.Split(typ, sub)
	name := tSlice[len(tSlice)-1]
	return name
}

func (i *Injector) has(key string) (bool, error) {
	i.rwLock.Lock()
	defer i.rwLock.Unlock()

	if injector.objs == nil {
		return false, errors.New("has fail, the injector with nil objs")
	}

	_, ok := i.objs[key]

	return ok, nil
}

func (i *Injector) put(key string, value interface{}) error {
	i.rwLock.Lock()
	defer i.rwLock.Unlock()

	if injector.objs == nil {
		return errors.New("put fail, the injector with nil objs")
	}

	i.objs[key] = reflect.ValueOf(value)

	return nil
}

func (i *Injector) delete(key string) error {
	i.rwLock.Lock()
	defer i.rwLock.Unlock()

	if injector.objs == nil {
		return errors.New("delete fail, the injector with nil objs")
	}

	delete(i.objs, key)

	return nil
}

func (i *Injector) get(key string) (reflect.Value, error) {
	i.rwLock.Lock()
	defer i.rwLock.Unlock()

	if injector.objs == nil {
		return reflect.Value{}, errors.New("get fail, the injector with nil objs")
	}

	val, _ := i.objs[key]

	return val, nil
}
