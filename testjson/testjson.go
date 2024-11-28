package testjson

import (
	"testing"

	"github.com/Jeffail/gabs/v2"
	"github.com/stretchr/testify/require"
)

type J struct {
	c          *gabs.Container
	readFields map[string]struct{}
}

func Unmarshal(t *testing.T, data []byte) J {
	c, err := gabs.ParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	return J{
		c:          c,
		readFields: make(map[string]struct{}),
	}
}

func (j J) RequireEqual(t *testing.T, expected any) {
	t.Helper()
	require.Equal(t, expected, j.c.Data())
}

func (j J) RequireNumber(t *testing.T, expected float64) {
	t.Helper()
	require.Equal(t, expected, j.c.Data())
}

func (j J) RequireString(t *testing.T, expected string) {
	t.Helper()
	actual, ok := j.c.Data().(string)
	if !ok {
		t.Errorf("expected string, got %T", j.c.Data())
		return
	}
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
		return
	}
}

func (j J) Get(path string) J {
	j.readFields[path] = struct{}{}
	return J{j.c.Path(path), make(map[string]struct{})}
}

func (j J) RequireNoAdditionalFields(t *testing.T) {
	fieldNames := listFlattenedFields(j.c)

	fields := make(map[string]struct{})
	for _, field := range fieldNames {
		fields[field] = struct{}{}
	}

	for field := range j.readFields {
		delete(fields, field)
	}

	if len(fields) > 0 {
		var keys []string
		for field := range fields {
			keys = append(keys, field)
		}
		t.Errorf("fields not read: %v", keys)
		return
	}
}

func listFlattenedFields(c *gabs.Container) []string {
	var fields []string
	for field, child := range c.ChildrenMap() {
		subfields := listFlattenedFields(child)
		if len(subfields) == 0 {
			fields = append(fields, field)
		} else {
			prefix := field + "."
			for _, subfield := range subfields {
				fields = append(fields, prefix+subfield)
			}
		}
	}
	return fields
}
