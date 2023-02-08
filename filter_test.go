package filter

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/the-control-group/go-timeutils"
)

func TestFilterOlderThan(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"path":"$.value","operator":"olderThan","value":"5m"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var _4mAgo = time.Now().Add(-4 * time.Minute).Format(time.RFC3339)
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":"`, _4mAgo, `"}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("4m ago should not pass")
		return
	}
	var _6mAgo = time.Now().Add(-6 * time.Minute).Format(timeutils.ISO8601_DATETIME)
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":"`, _6mAgo, `"}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("6m ago should pass")
		return
	}
}

func TestFilterIn(t *testing.T) {
	var filter = Filter{}
	dec := json.NewDecoder(bytes.NewBuffer([]byte(`{"path":"$.value","operator":"in","value":[1,2,3]}`)))
	dec.UseNumber()
	err := dec.Decode(&filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":1}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("int 1 should pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":"1"}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("string 1 should not pass")
		return
	}
	var msg3 interface{}
	msg3, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":4}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	pass, err = filter.Test(msg3)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("int 4 should not pass")
		return
	}
}

func TestFilterNotIn(t *testing.T) {
	var filter = Filter{}
	dec := json.NewDecoder(bytes.NewBuffer([]byte(`{"path":"$.value","operator":"not in","value":[1,2,3]}`)))
	dec.UseNumber()
	err := dec.Decode(&filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":1}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("int 1 should not pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":"1"}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("string 1 should pass")
		return
	}
	var msg3 interface{}
	msg3, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":4}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	pass, err = filter.Test(msg3)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("int 4 should pass")
		return
	}
}

func TestFilterNewerThan(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"path":"$.value","operator":"newerThan","value":"5m"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var _4mAgo = time.Now().Add(-4 * time.Minute).Format(time.RFC3339)
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":"`, _4mAgo, `"}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("4m ago should pass")
		return
	}
	var _6mAgo = time.Now().Add(-6 * time.Minute).Format(timeutils.ISO8601_DATETIME)
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":"`, _6mAgo, `"}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("6m ago should not pass")
		return
	}
}

func TestNotEqualNullFilter(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"path":"$.value","operator":"!=","value":null}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(`{"value":null}`))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("null should not pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(`{"value":"stuff"}`))
	if err != nil {
		t.Error("Failed to parse message 2", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("string value should pass")
		return
	}
	var msg3 interface{}
	msg3, err = decodeJSONMessage([]byte(`{}`))
	if err != nil {
		t.Error("Failed to parse message 3", err)
		return
	}
	pass, err = filter.Test(msg3)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("missing key should not pass")
		return
	}
}

func TestEqualStringFilter(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"path":"$.value","value":"test"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(`{"value":null}`))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("null should not pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(`{"value":"test"}`))
	if err != nil {
		t.Error("Failed to parse message 2", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("string value should pass")
		return
	}
	var msg3 interface{}
	msg3, err = decodeJSONMessage([]byte(`{}`))
	if err != nil {
		t.Error("Failed to parse message 3", err)
		return
	}
	pass, err = filter.Test(msg3)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("missing key should not pass")
		return
	}
}

func TestEqualStringFilterTemplate(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"template":"{{.value}}","value":"test"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(`{"value":null}`))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("null should not pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(`{"value":"test"}`))
	if err != nil {
		t.Error("Failed to parse message 2", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("string value should pass")
		return
	}
	var msg3 interface{}
	msg3, err = decodeJSONMessage([]byte(`{}`))
	if err != nil {
		t.Error("Failed to parse message 3", err)
		return
	}
	pass, err = filter.Test(msg3)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("missing key should not pass")
		return
	}
}

func TestEqualNullFilter(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"path":"$.value","value":null}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(`{"value":null}`))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("null should pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(`{"value":"test"}`))
	if err != nil {
		t.Error("Failed to parse message 2", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("string value should not pass")
		return
	}
	var msg3 interface{}
	msg3, err = decodeJSONMessage([]byte(`{}`))
	if err != nil {
		t.Error("Failed to parse message 3", err)
		return
	}
	pass, err = filter.Test(msg3)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("missing key should pass")
		return
	}
}

func TestGreaterThanFilter(t *testing.T) {
	var filter = Filter{}
	dec := json.NewDecoder(bytes.NewBuffer([]byte(`{"path":"$.value","operator":">","value":5}`)))
	dec.UseNumber()
	err := dec.Decode(&filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(`{"value":6}`))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("6 should pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(`{"value":null}`))
	if err != nil {
		t.Error("Failed to parse message 2", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("null should fail")
		return
	}
}

func TestFilterHttpStatus(t *testing.T) {
	t.SkipNow()
	var pass bool
	var filter = Filter{}
	dec := json.NewDecoder(bytes.NewBuffer([]byte(`{"template":"{{ $headers := dict \"Authorization\" (env \"AUTHX_TOKEN\") }}{{ (http \"GET\" (print \"https://bouncer.tcg.live/lists/ip_address_blacklist/items/\" .ip_address) $headers).StatusCode }}","operator":"eq","value":"200"}`)))
	dec.UseNumber()
	err := dec.Decode(&filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(`{"ip_address":"92.249.32.11"}`))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("92.249.32.11 should pass")
		return
	}
}

func TestFilterHttpStatusWithCache(t *testing.T) {
	t.SkipNow()
	var pass bool
	var filter = Filter{}
	dec := json.NewDecoder(bytes.NewBuffer([]byte(`{"template":"{{ $cachedVal := (cacheGet (print \"ip_address_blacklist:\" .ip_address)) }}{{ if $cachedVal }}{{ print $cachedVal }}{{ else }}{{ $headers := dict \"Authorization\" (env \"AUTHX_TOKEN\") }}{{ $res := ((http \"GET\" (print \"https://bouncer.tcg.live/lists/ip_address_blacklist/items/\" .ip_address) $headers).StatusCode) }}{{ cacheSet (print \"ip_address_blacklist:\" .ip_address) $res \"10m\" }}{{ end }}","operator":"eq","value":"200"}`)))
	dec.UseNumber()
	err := dec.Decode(&filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(`{"ip_address":"92.249.32.11"}`))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("92.249.32.11 should pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(`{"ip_address":"92.249.32.11"}`))
	if err != nil {
		t.Error("Failed to parse message 2", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("92.249.32.11 should pass again")
		return
	}
}

func TestFilterRegexMatch(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"path":"$.value","operator":"regex match","value":".*x{4}.*"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":"-xxxx"}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("-xxxx should pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":"-xxx-"}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 2", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("-xxx- should not pass")
		return
	}
}

func TestFilterRegexNoMatch(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"path":"$.value","operator":"regex no match","value":".*x{4}.*"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":"-xxxx"}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass {
		t.Error("-xxxx should not pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(strings.Join([]string{`{"value":"-xxx-"}`}, "")))
	if err != nil {
		t.Error("Failed to parse message 2", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("-xxx- should pass")
		return
	}
}

func decodeJSONMessage(data []byte) (interface{}, error) {
	var jsonData interface{}
	var msgReader = bytes.NewBuffer(data)
	var msgDecoder = json.NewDecoder(msgReader)
	msgDecoder.UseNumber()
	var decodeErr = msgDecoder.Decode(&jsonData)
	return jsonData, decodeErr
}

func TestValueAsTemplate(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"template":"{{.key}}","value":"{{.key}}"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg1 interface{}
	msg1, err = decodeJSONMessage([]byte(`{"key":null}`))
	if err != nil {
		t.Error("Failed to parse message 1", err)
		return
	}
	var pass bool
	pass, err = filter.Test(msg1)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("null should pass")
		return
	}
	var msg2 interface{}
	msg2, err = decodeJSONMessage([]byte(`{"key":"test"}`))
	if err != nil {
		t.Error("Failed to parse message 2", err)
		return
	}
	pass, err = filter.Test(msg2)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("string value should pass")
		return
	}
	var msg3 interface{}
	msg3, err = decodeJSONMessage([]byte(`{}`))
	if err != nil {
		t.Error("Failed to parse message 3", err)
		return
	}
	pass, err = filter.Test(msg3)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass {
		t.Error("missing key should pass")
		return
	}
}

func TestRandomBoolean(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"template":"{{ randomInt 0 1 }}","value":"0"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg interface{}
	msg, err = decodeJSONMessage([]byte(`{}`))
	if err != nil {
		t.Error("Failed to parse message", err)
		return
	}
	rand.Seed(0)
	var pass1 bool
	pass1, err = filter.Test(msg)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if !pass1 {
		t.Error("should pass")
		return
	}
	err = json.Unmarshal([]byte(`{"template":"{{ randomInt 0 1 }}","value":"0"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	rand.Seed(1)
	var pass2 bool
	pass2, err = filter.Test(msg)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	if pass2 {
		t.Error("should not pass")
		return
	}
}

func TestRandom(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"template":"{{ randomInt 0 1 }}","value":"0"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	var msg interface{}
	msg, err = decodeJSONMessage([]byte(`{}`))
	if err != nil {
		t.Error("Failed to parse message", err)
		return
	}
	rand.Seed(0)
	var pass bool
	pass, err = filter.Test(msg)
	if err != nil {
		t.Error("Filter test failed", err)
		return
	}
	t.Logf("%t", pass)
}

func TestRandomToJSON(t *testing.T) {
	var filter = Filter{}
	err := json.Unmarshal([]byte(`{"template":"{{ randomInt 0 1 }}","value":"0"}`), &filter)
	if err != nil {
		t.Error("Failed to parse filter", err)
		return
	}
	b, err := json.Marshal(filter)
	if err != nil {
		t.Error("Error marshaling to json", err)
		return
	}
	if string(b) != `{"template":"{{randomInt 0 1}}","path":"","value":"0","operator":"","requeue":false,"or":null,"and":null}` {
		t.Fail()
	}
}
