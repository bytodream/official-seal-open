package structx

import "reflect"

// DiffFields 比较两个结构体中指定字段的属性值是否发生改变，并返回改变的字段名
func DiffFields[T any](t1 T, t2 T, fieldsToCompare ...string) []string {
	var changedFields []string

	// 指针入参统一解引用，否则 FieldByName 对指针类型直接 panic
	oldValue := Indirect(reflect.ValueOf(t1))
	newValue := Indirect(reflect.ValueOf(t2))
	if oldValue.Kind() != reflect.Struct || newValue.Kind() != reflect.Struct {
		return changedFields
	}

	for _, fieldName := range fieldsToCompare {
		oldFieldValue := oldValue.FieldByName(fieldName)
		newFieldValue := newValue.FieldByName(fieldName)

		// 字段不存在时得到零值 Value，Interface() 会 panic，直接跳过
		if !oldFieldValue.IsValid() || !newFieldValue.IsValid() {
			continue
		}

		if !reflect.DeepEqual(oldFieldValue.Interface(), newFieldValue.Interface()) {
			changedFields = append(changedFields, fieldName)
		}
	}

	return changedFields
}
