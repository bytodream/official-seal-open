package model

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// FlexString 前端常把数字与空字符串混传。
type FlexString string

// UnmarshalJSON 宽松解析 JSON 字段（数字/空串混传）。
func (s *FlexString) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*s = ""
		return nil
	}
	if b[0] == '"' {
		var str string
		if err := json.Unmarshal(b, &str); err != nil {
			return err
		}
		*s = FlexString(str)
		return nil
	}
	var num json.Number
	if err := json.Unmarshal(b, &num); err != nil {
		return err
	}
	*s = FlexString(num.String())
	return nil
}

// String 返回 FlexString 的字符串形式。
func (s FlexString) String() string { return string(s) }

// FlexBool 支持 bool 与 "" / null。
type FlexBool struct {
	Value *bool
}

// UnmarshalJSON 宽松解析 JSON 字段（支持 bool 与 "" / null）。
func (b *FlexBool) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) || bytes.Equal(data, []byte(`""`)) {
		b.Value = nil
		return nil
	}
	var v bool
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	b.Value = &v
	return nil
}

// BoolDefault 解析 FlexBool，缺省为 defaultVal。
func (b FlexBool) BoolDefault(def bool) bool {
	if b.Value == nil {
		return def
	}
	return *b.Value
}

// FlexFloat32 支持 float、整型、数字字符串、空字符串。
type FlexFloat32 struct {
	Value *float32
}

// UnmarshalJSON 宽松解析 JSON 字段（数字/空串混传）。
func (f *FlexFloat32) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) || bytes.Equal(b, []byte(`""`)) {
		f.Value = nil
		return nil
	}
	parse := func(s string) error {
		if s == "" {
			f.Value = nil
			return nil
		}
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		fv := float32(v)
		f.Value = &fv
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		return parse(s)
	}
	var v float64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	fv := float32(v)
	f.Value = &fv
	return nil
}

// FlexInt 支持 int、数字字符串、空字符串。
type FlexInt struct {
	Value *int
}

// UnmarshalJSON 宽松解析 JSON 字段（数字/空串混传）。
func (i *FlexInt) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) || bytes.Equal(b, []byte(`""`)) {
		i.Value = nil
		return nil
	}
	parse := func(s string) error {
		if s == "" {
			i.Value = nil
			return nil
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		i.Value = &v
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		return parse(s)
	}
	var v int
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Value = &v
	return nil
}

// FlexTransparencyString 接口里 transparency 可能是 30 或 "30"。
func FlexTransparencyString(b []byte) (string, error) {
	b = bytes.TrimSpace(b)
	if bytes.Equal(b, []byte("null")) {
		return "", nil
	}
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return "", err
		}
		return s, nil
	}
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return "", err
	}
	return n.String(), nil
}

// PtrFloat32 返回 float32 值的堆副本指针。
func PtrFloat32(v float32) *float32 { return &v }

// ParseFlexTransparencyField 解析 electronicModel.transparency 字段。
func ParseFlexTransparencyField(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	return FlexTransparencyString(raw)
}

// MustAtoi 将字符串转为 int，失败时返回 0。
func MustAtoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
