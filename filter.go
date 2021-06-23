package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/nickcarenza/go-template"
	"github.com/the-control-group/go-jsonpath"
	"github.com/the-control-group/go-timeutils"
)

// Filter ...
type Filter struct {
	Template *template.Template `json:"template"`
	Path     jsonpath.JsonPath  `json:"path"`
	Value    interface{}        `json:"value"`
	Operator string             `json:"operator"`
	Requeue  bool               `json:"requeue"`
}

// Test ...
func (f *Filter) Test(msg interface{}) (bool, error) {
	var val interface{}
	var err error
	if f.Template != nil {
		var b bytes.Buffer
		err = f.Template.Execute(&b, msg)
		if err != nil {
			return false, err
		}
		val = b.String()
	} else {
		val, _ = jsonpath.GetPathValue(msg, f.Path)
	}
	if n, ok := val.(json.Number); ok {
		val, err = n.Float64()
		if err != nil {
			return false, fmt.Errorf("TypeAssertionError")
		}
	} else if n, ok := val.(int); ok {
		val = float64(n)
	} else if n, ok := val.(int64); ok {
		val = float64(n)
	}
	if n, ok := f.Value.(json.Number); ok {
		var err error
		f.Value, err = n.Float64()
		if err != nil {
			return false, fmt.Errorf("TypeAssertionError")
		}
	}
	switch f.Operator {
	case "!=", "<>", "ne", "doesn't equal", "not equal to":
		return f.Value != val, nil
	case "==", "eq", "equal", "equals":
		return f.Value == val, nil
	case "in", "not in":
		rt := reflect.TypeOf(f.Value)
		switch rt.Kind() {
		case reflect.Slice:
			s := f.Value.([]interface{})
			for _, v := range s {
				f2 := Filter{
					Path:     f.Path,
					Value:    v,
					Operator: "==",
					Requeue:  f.Requeue,
				}
				ok, err := f2.Test(msg)
				if ok || err != nil {
					if f.Operator == "not in" {
						return !ok, err
					}
					return ok, err
				}
			}
			if f.Operator == "not in" {
				return true, nil
			}
			return false, nil
		default:
			return false, fmt.Errorf("TypeAssertionError")
		}
	case "<", "lt", "less than",
		">", "gt", "greater than",
		">=", "ge", "gte", "greater than or equal to",
		"<=", "le", "lte", "less than or equal to":
		var fNum, vNum float64
		var err error
		if val == nil {
			return false, nil
		}
		fNum, err = interfaceToFloat64(f.Value)
		if err != nil {
			return false, err
		}
		vNum, err = interfaceToFloat64(val)
		if err != nil {
			return false, err
		}
		switch f.Operator {
		case "<", "lt", "less than":
			return vNum < fNum, nil
		case ">", "gt", "greater than":
			return vNum > fNum, nil
		case ">=", "ge", "gte", "greater than or equal to":
			return vNum >= fNum, nil
		case "<=", "le", "lte", "less than or equal to":
			return vNum <= fNum, nil
		default:
			return false, fmt.Errorf("Impossible condition")
		}
	case "olderThan", "newerThan":
		var fVal string
		var ok bool
		var err error
		var dVal timeutils.ApproxBigDuration
		var tStr string
		var tVal time.Time
		if val == nil {
			return false, nil
		}
		fVal, ok = f.Value.(string)
		if !ok {
			return false, fmt.Errorf("TypeAssertionError")
		}
		dVal, err = timeutils.ParseApproxBigDuration([]byte(fVal))
		if err != nil {
			return false, err
		}
		tStr, ok = val.(string)
		if !ok {
			return false, fmt.Errorf("TypeAssertionError")
		}
		tVal, err = timeutils.ParseAny(tStr)
		if err != nil {
			return false, err
		}
		switch f.Operator {
		case "olderThan":
			return time.Since(tVal) > time.Duration(dVal), nil
		case "newerThan":
			return time.Since(tVal) < time.Duration(dVal), nil
		default:
			return false, fmt.Errorf("Impossible condition")
		}
	default:
		return f.Value == val, nil
	}
}
