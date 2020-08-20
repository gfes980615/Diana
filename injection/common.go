package injection

import (
	"errors"
	"reflect"
	"strings"
)

func ParseSuffix(suffix string) ([]reflect.Value, error) {
	if injector == nil {
		return nil, errors.New("fail with nil injector")
	}

	injector.rwLock.Lock()
	defer injector.rwLock.Unlock()

	ret := make([]reflect.Value, 0)
	for k, v := range injector.objs {
		if strings.HasSuffix(k, suffix) {
			ret = append(ret, v)
		}
	}

	return ret, nil
}
