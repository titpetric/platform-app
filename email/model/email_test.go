package model

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/titpetric/platform/pkg/require"
)

// TestEmail ensures that the email tables match.
// We want to swap the table name as we move emails to different queues, but use the `Email` type.
func TestEmail(t *testing.T) {
	a := Email{}
	b := EmailFailed{}
	c := EmailSent{}

	diff1, err1 := CompareStructShape(a, b)
	diff2, err2 := CompareStructShape(b, c)
	diff3, err3 := CompareStructShape(c, a)

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)

	require.Empty(t, diff1)
	require.Empty(t, diff2)
	require.Empty(t, diff3)
}

type FieldInfo struct {
	Name string
	Type string
	Tag  string
}

func ExtractStructShape(v interface{}) ([]FieldInfo, error) {
	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, nil
	}

	fields := make([]FieldInfo, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fields = append(fields, FieldInfo{
			Name: f.Name,
			Type: f.Type.String(),
			Tag:  string(f.Tag),
		})
	}

	return fields, nil
}

func CompareStructShape(a, b interface{}) (string, error) {
	shapeA, err := ExtractStructShape(a)
	if err != nil {
		return "", err
	}
	shapeB, err := ExtractStructShape(b)
	if err != nil {
		return "", err
	}

	// go-cmp does the diff for us
	diff := cmp.Diff(shapeA, shapeB)
	return diff, nil
}
