package helpers

import (
	"context"
	"testing"

	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Set      types.Set
	Str      types.String
	In       types.Int64
	Fl       types.Float64
	Boo      types.Bool
	SeedTime types.Int64
}

func TestWriteStringField(t *testing.T) {
	t.Parallel()

	field := whisparr.NewField()
	field.SetName("str")
	field.SetValue("string")

	tests := map[string]struct {
		fieldOutput whisparr.Field
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: *field,
			written:     Test{},
			expected:    Test{Str: types.StringValue("string")},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			WriteStringField(&test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteBoolField(t *testing.T) {
	t.Parallel()

	field := whisparr.NewField()
	field.SetName("boo")
	field.SetValue(true)

	tests := map[string]struct {
		fieldOutput whisparr.Field
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: *field,
			written:     Test{},
			expected:    Test{Boo: types.BoolValue(true)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			WriteBoolField(&test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteIntField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name     string
		written  Test
		expected Test
		// use float to simulate unmarshal response
		value float64
	}{
		"working": {
			name:     "in",
			value:    float64(50),
			written:  Test{},
			expected: Test{In: types.Int64Value(50)},
		},
		"seedtime": {
			name:     "seedCriteria.seedTime",
			value:    float64(50),
			written:  Test{},
			expected: Test{SeedTime: types.Int64Value(50)},
		},
	}
	for name, test := range tests {
		test := test

		field := whisparr.NewField()
		field.SetName(test.name)
		field.SetValue(test.value)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			WriteIntField(field, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteFloatField(t *testing.T) {
	t.Parallel()

	field := whisparr.NewField()
	field.SetName("fl")
	field.SetValue(float64(3.5))

	tests := map[string]struct {
		fieldOutput whisparr.Field
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: *field,
			written:     Test{},
			expected:    Test{Fl: types.Float64Value(3.5)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			WriteFloatField(&test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteIntSliceField(t *testing.T) {
	t.Parallel()

	field := whisparr.NewField()
	field.SetName("set")
	// use interface to simulate unmarshal response
	field.SetValue(append(make([]interface{}, 0), 1, 2))

	tests := map[string]struct {
		fieldOutput whisparr.Field
		set         []int64
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: *field,
			written:     Test{},
			set:         []int64{1, 2},
			expected:    Test{Set: types.SetValueMust(types.Int64Type, nil)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.expected.Set.Type(context.Background()), &test.expected.Set)
			WriteIntSliceField(context.Background(), &test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteStringSliceField(t *testing.T) {
	t.Parallel()

	field := whisparr.NewField()
	field.SetName("set")
	// use interface to simulate unmarshal response
	field.SetValue(append(make([]interface{}, 0), "test1", "test2"))

	tests := map[string]struct {
		fieldOutput whisparr.Field
		set         []string
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: *field,
			written:     Test{},
			set:         []string{"test1", "test2"},
			expected:    Test{Set: types.SetValueMust(types.StringType, nil)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.expected.Set.Type(context.Background()), &test.expected.Set)
			WriteStringSliceField(context.Background(), &test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestReadStringField(t *testing.T) {
	t.Parallel()

	field := whisparr.NewField()
	field.SetName("str")
	field.SetValue("string")

	tests := map[string]struct {
		expected  *whisparr.Field
		name      string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Str: types.StringValue("string"),
			},
			name:     "str",
			expected: field,
		},
		"nil": {
			fieldCase: Test{},
			name:      "str",
			expected:  nil,
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := ReadStringField(test.name, &test.fieldCase)
			assert.Equal(t, test.expected, field)
		})
	}
}

func TestReadIntField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name      string
		tfName    string
		fieldCase Test
		value     int
	}{
		"working": {
			fieldCase: Test{
				In: types.Int64Value(10),
			},
			name:   "in",
			tfName: "in",
			value:  10,
		},
		"nil": {
			fieldCase: Test{},
			name:      "in",
			tfName:    "in",
			value:     0,
		},
		"seedtime": {
			fieldCase: Test{
				SeedTime: types.Int64Value(10),
			},
			name:   "seedCriteria.seedTime",
			tfName: "seedTime",
			value:  10,
		},
	}
	for name, test := range tests {
		test := test

		expected := whisparr.NewField()
		expected.SetName(test.name)
		expected.SetValue(int64(test.value))

		if test.value == 0 {
			expected = nil
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := ReadIntField(test.tfName, &test.fieldCase)
			assert.Equal(t, expected, field)
		})
	}
}

func TestReadBoolField(t *testing.T) {
	t.Parallel()

	field := whisparr.NewField()
	field.SetName("boo")
	field.SetValue(true)

	tests := map[string]struct {
		expected  *whisparr.Field
		name      string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Boo: types.BoolValue(true),
			},
			name:     "boo",
			expected: field,
		},
		"nil": {
			fieldCase: Test{},
			name:      "boo",
			expected:  nil,
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := ReadBoolField(test.name, &test.fieldCase)
			assert.Equal(t, test.expected, field)
		})
	}
}

func TestReadFloatField(t *testing.T) {
	t.Parallel()

	field := whisparr.NewField()
	field.SetName("fl")
	field.SetValue(3.5)

	tests := map[string]struct {
		expected  *whisparr.Field
		name      string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Fl: types.Float64Value(3.5),
			},
			name:     "fl",
			expected: field,
		},
		"nil": {
			fieldCase: Test{},
			name:      "fl",
			expected:  nil,
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := ReadFloatField(test.name, &test.fieldCase)
			assert.Equal(t, test.expected, field)
		})
	}
}

func TestReadStringSliceField(t *testing.T) {
	t.Parallel()

	field := whisparr.NewField()
	field.SetName("set")
	field.SetValue([]string{"test1", "test2"})

	tests := map[string]struct {
		expected  *whisparr.Field
		name      string
		set       []string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Set: types.SetValueMust(types.StringType, nil),
			},
			name:     "set",
			expected: field,
			set:      []string{"test1", "test2"},
		},
		"nil": {
			fieldCase: Test{
				Set: types.SetValueMust(types.StringType, nil),
			},
			name:     "set",
			expected: nil,
			set:      []string{},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.fieldCase.Set.Type(context.Background()), &test.fieldCase.Set)
			field := ReadStringSliceField(context.Background(), test.name, &test.fieldCase)
			assert.Equal(t, test.expected, field)
		})
	}
}

func TestReadIntSliceField(t *testing.T) {
	t.Parallel()

	field := whisparr.NewField()
	field.SetName("set")
	field.SetValue([]int64{1, 2})

	tests := map[string]struct {
		expected  *whisparr.Field
		name      string
		set       []int64
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Set: types.SetValueMust(types.Int64Type, nil),
			},
			name:     "set",
			expected: field,
			set:      []int64{1, 2},
		},
		"nil": {
			fieldCase: Test{
				Set: types.SetValueMust(types.Int64Type, nil),
			},
			name:     "set",
			expected: nil,
			set:      []int64{},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.fieldCase.Set.Type(context.Background()), &test.fieldCase.Set)
			field := ReadIntSliceField(context.Background(), test.name, &test.fieldCase)
			assert.Equal(t, test.expected, field)
		})
	}
}
