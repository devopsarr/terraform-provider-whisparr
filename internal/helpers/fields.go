package helpers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/slices"
)

type fieldException struct {
	apiName string
	tfName  string
}

func getFieldExceptions() []fieldException {
	return []fieldException{
		{
			apiName: "tags",
			tfName:  "fieldTags",
		},
		{
			apiName: "seedCriteria.seedTime",
			tfName:  "seedTime",
		},
		{
			apiName: "seedCriteria.seedRatio",
			tfName:  "seedRatio",
		},
		{
			apiName: "seedCriteria.seasonPackSeedTime",
			tfName:  "seasonPackSeedTime",
		},
		{
			apiName: "filterCriteria.minVoteAverage",
			tfName:  "minVoteAverage",
		},
		{
			apiName: "filterCriteria.minVotes",
			tfName:  "minVotes",
		},
		{
			apiName: "filterCriteria.includeGenreIds",
			tfName:  "includeGenreIds",
		},
		{
			apiName: "filterCriteria.excludeGenreIds",
			tfName:  "excludeGenreIds",
		},
		{
			apiName: "filterCriteria.languageCode",
			tfName:  "languageCode",
		},
		{
			apiName: "filterCriteria.certification",
			tfName:  "tmdbCertification",
		},
		{
			apiName: "listType",
			tfName:  "userListType",
		},
	}
}

// selectTFName identifies the TF name starting from API name.
func selectTFName(name string) string {
	for _, f := range getFieldExceptions() {
		if f.apiName == name {
			name = f.tfName
		}
	}

	return name
}

// selectAPIName identifies the API name starting from TF name.
func selectAPIName(name string) string {
	for _, f := range getFieldExceptions() {
		if f.tfName == name {
			name = f.apiName
		}
	}

	return name
}

// selectWriteField identifies which struct field should be written.
func selectWriteField(fieldOutput *whisparr.Field, fieldCase interface{}) reflect.Value {
	fieldName := selectTFName(fieldOutput.GetName())
	value := reflect.ValueOf(fieldCase).Elem()

	return value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
}

// selectReadField identifies which struct field should be read.
func selectReadField(name string, fieldCase interface{}) reflect.Value {
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()

	return value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
}

// setField sets the whisparr field value.
func setField(name string, value interface{}) *whisparr.Field {
	field := whisparr.NewField()
	field.SetName(name)
	field.SetValue(value)

	return field
}

// writeStringField writes a whisparr string field into struct field.
func writeStringField(fieldOutput *whisparr.Field, fieldCase interface{}) {
	stringValue := fmt.Sprint(fieldOutput.GetValue())

	v := reflect.ValueOf(types.StringValue(stringValue))
	if fieldOutput.GetValue() == nil {
		v = reflect.ValueOf(types.StringNull())
	}

	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// writeBoolField writes a whisparr bool field into struct field.
func writeBoolField(fieldOutput *whisparr.Field, fieldCase interface{}) {
	boolValue, _ := fieldOutput.GetValue().(bool)

	v := reflect.ValueOf(types.BoolValue(boolValue))
	if fieldOutput.GetValue() == nil {
		v = reflect.ValueOf(types.BoolNull())
	}

	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// writeIntField writes a whisparr int field into struct field.
func writeIntField(fieldOutput *whisparr.Field, fieldCase interface{}) {
	intValue, _ := fieldOutput.GetValue().(float64)

	v := reflect.ValueOf(types.Int64Value(int64(intValue)))
	if fieldOutput.GetValue() == nil {
		v = reflect.ValueOf(types.Int64Null())
	}

	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// writeFloatField writes a whisparr float field into struct field.
func writeFloatField(fieldOutput *whisparr.Field, fieldCase interface{}) {
	floatValue, _ := fieldOutput.GetValue().(float64)

	v := reflect.ValueOf(types.Float64Value(floatValue))
	if fieldOutput.GetValue() == nil {
		v = reflect.ValueOf(types.Float64Null())
	}

	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// writeStringSliceField writes a whisparr string slice field into struct field.
func writeStringSliceField(ctx context.Context, fieldOutput *whisparr.Field, fieldCase interface{}) {
	sliceValue, _ := fieldOutput.GetValue().([]interface{})
	setValue := types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, sliceValue, setValue.Type(ctx), &setValue)
	v := reflect.ValueOf(setValue)
	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// writeIntSliceField writes a whisparr int slice field into struct field.
func writeIntSliceField(ctx context.Context, fieldOutput *whisparr.Field, fieldCase interface{}) {
	sliceValue, _ := fieldOutput.GetValue().([]interface{})
	setValue := types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, sliceValue, setValue.Type(ctx), &setValue)
	v := reflect.ValueOf(setValue)
	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// readStringField reads from a string struct field and return a whisparr field.
func readStringField(name string, fieldCase interface{}) *whisparr.Field {
	fieldName := selectAPIName(name)
	stringField := (*types.String)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if !stringField.IsNull() && !stringField.IsUnknown() {
		return setField(fieldName, stringField.ValueString())
	}

	return nil
}

// readBoolField reads from a bool struct field and return a whisparr field.
func readBoolField(name string, fieldCase interface{}) *whisparr.Field {
	fieldName := selectAPIName(name)
	boolField := (*types.Bool)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if !boolField.IsNull() && !boolField.IsUnknown() {
		return setField(fieldName, boolField.ValueBool())
	}

	return nil
}

// readIntField reads from a int struct field and return a whisparr field.
func readIntField(name string, fieldCase interface{}) *whisparr.Field {
	fieldName := selectAPIName(name)
	intField := (*types.Int64)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if !intField.IsNull() && !intField.IsUnknown() {
		return setField(fieldName, intField.ValueInt64())
	}

	return nil
}

// readFloatField reads from a float struct field and return a whisparr field.
func readFloatField(name string, fieldCase interface{}) *whisparr.Field {
	fieldName := selectAPIName(name)
	floatField := (*types.Float64)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if !floatField.IsNull() && !floatField.IsUnknown() {
		return setField(fieldName, floatField.ValueFloat64())
	}

	return nil
}

// readStringSliceField reads from a string slice struct field and return a whisparr field.
func readStringSliceField(ctx context.Context, name string, fieldCase interface{}) *whisparr.Field {
	fieldName := selectAPIName(name)
	sliceField := (*types.Set)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if len(sliceField.Elements()) != 0 {
		slice := make([]string, len(sliceField.Elements()))
		tfsdk.ValueAs(ctx, sliceField, &slice)

		return setField(fieldName, slice)
	}

	return nil
}

// readIntSliceField reads from a int slice struct field and return a whisparr field.
func readIntSliceField(ctx context.Context, name string, fieldCase interface{}) *whisparr.Field {
	fieldName := selectAPIName(name)
	sliceField := (*types.Set)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if len(sliceField.Elements()) != 0 {
		slice := make([]int64, len(sliceField.Elements()))
		tfsdk.ValueAs(ctx, sliceField, &slice)

		return setField(fieldName, slice)
	}

	return nil
}

// Fields contains all the field lists of a specific resource per type.
type Fields struct {
	Bools                  []string
	BoolsExceptions        []string
	Ints                   []string
	IntsExceptions         []string
	Strings                []string
	StringsExceptions      []string
	Floats                 []string
	FloatsExceptions       []string
	IntSlices              []string
	IntSlicesExceptions    []string
	StringSlices           []string
	StringSlicesExceptions []string
}

// getList return a specific list of fields.
func (f Fields) getList(list string) []string {
	r := reflect.ValueOf(f)
	output, _ := reflect.Indirect(r).FieldByName(list).Interface().([]string)

	return output
}

// ReadFields takes in input a field container and populates a whisparr.Field slice.
func ReadFields(ctx context.Context, fieldContainer interface{}, fieldLists Fields) []*whisparr.Field {
	var output []*whisparr.Field

	// Map each list to its read function.
	readFuncs := map[string]func(string, interface{}) *whisparr.Field{
		"Bools":   readBoolField,
		"Ints":    readIntField,
		"Floats":  readFloatField,
		"Strings": readStringField,
		"StringSlices": func(name string, fieldContainer interface{}) *whisparr.Field {
			return readStringSliceField(ctx, name, fieldContainer)
		},
		"IntSlices": func(name string, fieldContainer interface{}) *whisparr.Field {
			return readIntSliceField(ctx, name, fieldContainer)
		},
	}

	// Loop over the map to populate the whisparr.Field slice.
	for fieldType, readFunc := range readFuncs {
		for _, f := range fieldLists.getList(fieldType) {
			if field := readFunc(f, fieldContainer); field != nil {
				output = append(output, field)
			}
		}
	}

	return output
}

// WriteFields takes in input a whisparr.Field slice and populate the relevant container fields.
func WriteFields(ctx context.Context, fieldContainer interface{}, fields []*whisparr.Field, fieldLists Fields) {
	// Map each list to its write function.
	writeFuncs := map[string]func(*whisparr.Field, interface{}){
		"Bools":             writeBoolField,
		"BoolsExceptions":   writeBoolField,
		"Ints":              writeIntField,
		"IntsExceptions":    writeIntField,
		"Strings":           writeStringField,
		"StringsExceptions": writeStringField,
		"Floats":            writeFloatField,
		"FloatsExceptions":  writeFloatField,
		"IntSlices": func(fieldOutput *whisparr.Field, fieldContainer interface{}) {
			writeIntSliceField(ctx, fieldOutput, fieldContainer)
		},
		"IntSlicesExceptions": func(fieldOutput *whisparr.Field, fieldContainer interface{}) {
			writeIntSliceField(ctx, fieldOutput, fieldContainer)
		},
		"StringSlices": func(fieldOutput *whisparr.Field, fieldContainer interface{}) {
			writeStringSliceField(ctx, fieldOutput, fieldContainer)
		},
		"StringSlicesExceptions": func(fieldOutput *whisparr.Field, fieldContainer interface{}) {
			writeStringSliceField(ctx, fieldOutput, fieldContainer)
		},
	}

	// Loop over each field and populate the related container field with the corresponding write function.
	for _, f := range fields {
		fieldName := f.GetName()
		for listName, writeFunc := range writeFuncs {
			if slices.Contains(fieldLists.getList(listName), fieldName) {
				writeFunc(f, fieldContainer)

				break
			}
		}
	}
}
