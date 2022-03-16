package validator

import (
	"reflect"
	"testing"
)

func TestValidateStruct(t *testing.T) {
	type testStruct struct {
		Name string `validate:"required,min=5"`
	}

	testCases := []struct {
		name    string
		testObj interface{}
		want    []ValidationError
	}{
		{
			name: "Valid struct",
			testObj: testStruct{
				Name: "Success Test",
			},
			want: nil,
		},
		{
			name: "Valid pointer struct",
			testObj: &testStruct{
				Name: "Success Test",
			},
			want: nil,
		},
		{
			name: "Invalid struct",
			testObj: testStruct{
				Name: "Err",
			},
			want: []ValidationError{
				{
					Detail: "Name must be at least 5 characters in length",
					Source: "Name",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ValidateStruct(tc.testObj); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ValidateStruct() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestValidateVar(t *testing.T) {
	type args struct {
		field interface{}
		tag   string
	}
	testCases := []struct {
		name string
		args args
		want []ValidationError
	}{
		{
			name: "Valid var",
			args: args{
				field: "Test Success",
				tag:   "required,min=5",
			},
			want: nil,
		},
		{
			name: "Invalid var",
			args: args{
				field: "Test",
				tag:   "required,min=5",
			},
			want: []ValidationError{
				{
					Detail: " must be at least 5 characters in length",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ValidateVar(tc.args.field, tc.args.tag); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ValidateVar() = %v, want %v", got, tc.want)
			}
		})
	}
}
