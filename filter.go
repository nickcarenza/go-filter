package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nickcarenza/go-template"
	"github.com/robertkrimen/otto"
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
	Or       *Filter            `json:"or"`
	And      *Filter            `json:"and"`
	Script   *ScriptFilter      `json:"script"`
}

type ScriptFilter struct {
	Interpreter string                 `json:"interpreter"`
	Script      string                 `json:"script"`
	ScriptFile  string                 `json:"scriptFile"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Test evaluates if the filter or the Or clause passes
func (f *Filter) Test(msg interface{}) (bool, error) {
	pass, err := Test(f, msg)
	if err != nil {
		return false, err
	}
	if !pass && f.Or != nil {
		pass, err = f.Or.Test(msg)
		if err != nil {
			return false, err
		}
	}
	if pass && f.And != nil {
		pass, err = f.And.Test(msg)
		if err != nil {
			return false, err
		}
	}
	return pass, err
}

func Test(f *Filter, msg interface{}) (bool, error) {
	var val interface{}
	var err error
	var fVal interface{}
	if f.Script != nil {
		switch strings.ToLower(f.Script.Interpreter) {
		case "javascript", "js", "es5":
			vm := otto.New()
			vm.Set("input", msg)
			vm.Set("metadata", f.Script.Metadata)
			var res otto.Value
			if f.Script.ScriptFile != "" {
				dat, err := os.ReadFile(f.Script.ScriptFile)
				if err != nil {
					return false, err
				}
				res, err = vm.Run(dat)
				if err != nil {
					return false, err
				}
				b, err := res.ToBoolean()
				if err != nil {
					return false, err
				}
				return b, nil
			} else {
				res, err = vm.Run(f.Script.Script)
				if err != nil {
					return false, err
				}
				b, err := res.ToBoolean()
				if err != nil {
					return false, err
				}
				return b, nil
			}
		default:
			return false, fmt.Errorf("unsupported interpreter %s", f.Script.Interpreter)
		}
	}
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
		fVal, err = n.Float64()
		if err != nil {
			return false, fmt.Errorf("TypeAssertionError")
		}
	} else if str, ok := f.Value.(string); ok {
		fVal, err = template.Interpolate(msg, str)
		if err != nil {
			return false, err
		}
	} else {
		fVal = f.Value
	}
	switch f.Operator {
	case "!=", "<>", "ne", "doesn't equal", "not equal to":
		return fVal != val, nil
	case "=", "==", "eq", "equal", "equals":
		return fVal == val, nil
	case "in", "not in":
		rt := reflect.TypeOf(fVal)
		switch rt.Kind() {
		case reflect.Slice:
			s := fVal.([]interface{})
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
		fNum, err = interfaceToFloat64(fVal)
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
			return false, fmt.Errorf("impossible condition")
		}
	case "olderThan", "older than", "older",
		"newerThan", "newer than", "newer":
		var rVal string
		var ok bool
		var err error
		var dVal timeutils.ApproxBigDuration
		var tStr string
		var tVal time.Time
		if val == nil {
			return false, nil
		}
		rVal, ok = fVal.(string)
		if !ok {
			return false, fmt.Errorf("TypeAssertionError")
		}
		dVal, err = timeutils.ParseApproxBigDuration([]byte(rVal))
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
		case "olderThan", "older than", "older":
			return time.Since(tVal) > time.Duration(dVal), nil
		case "newerThan", "newer than", "newer":
			return time.Since(tVal) < time.Duration(dVal), nil
		default:
			return false, fmt.Errorf("impossible condition")
		}
	case "regexMatch", "regex match",
		"regexNoMatch", "regex no match":
		var re *regexp.Regexp
		var err error
		var rVal string
		var tStr string
		var ok bool
		rVal, ok = fVal.(string)
		if !ok {
			return false, fmt.Errorf("TypeAssertionError")
		}
		re, err = regexp.Compile(rVal)
		if err != nil {
			return false, err
		}
		tStr, ok = val.(string)
		if !ok {
			return false, fmt.Errorf("TypeAssertionError")
		}
		switch f.Operator {
		case "regexMatch", "regex match":
			return re.MatchString(tStr), nil
		case "regexNoMatch", "regex no match":
			return !re.MatchString(tStr), nil
		default:
			return false, fmt.Errorf("impossible condition")
		}
	default:
		return fVal == val, nil
	}
}

// func AndOr(bool, and *Filter, or *Filter) (bool, error) {

// }

func interfaceToFloat64(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
		return f, nil
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return 0, err
		}
		return f, nil
	default:
		return 0, fmt.Errorf("TypeAssertionError")
	}
}
