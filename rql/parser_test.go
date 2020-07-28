package rql

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		args    args
		want    *Ast
		wantErr bool
	}{
		{
			"SET STATEMENT",
			args{`SET data "Hello World";`},
			&Ast{
				Statements: []*Statement{
					{
						SetStatement: &SetStatement{
							key: "data",
							val: "Hello World",
						},
						typ: SetType,
					},
				},
			},
			false,
		},
		{
			"SET STATEMENT WITH EXPIRY",
			args{`SET data "Hello World" 234;`},
			&Ast{
				Statements: []*Statement{
					{
						SetStatement: &SetStatement{
							key: "data",
							val: "Hello World",
							exp: 234,
						},
						typ: SetType,
					},
				},
			},
			false,
		},
		{
			"MULTI SET STATEMENTS",
			args{`SET data "Hello World" 234; SET data1 3454 565;`},
			&Ast{
				Statements: []*Statement{
					{
						SetStatement: &SetStatement{
							key: "data",
							val: "Hello World",
							exp: 234,
						},
						typ: SetType,
					},
					{
						SetStatement: &SetStatement{
							key: "data1",
							val: "3454",
							exp: 565,
						},
						typ: SetType,
					},
				},
			},
			false,
		},
		{
			"GET STATEMENT",
			args{`GET data data1 data2 data3;`},
			&Ast{
				Statements: []*Statement{
					{
						GetStatement: &GetStatement{
							keys: []string{"data", "data1", "data2", "data3"},
						},
						typ: GetType,
					},
				},
			},
			false,
		},
		{
			"DELETE STATEMENT",
			args{`DEL data data1 data2 data3;`},
			&Ast{
				Statements: []*Statement{
					{
						DeleteStatement: &DeleteStatement{
							keys: []string{"data", "data1", "data2", "data3"},
						},
						typ: DeleteType,
					},
				},
			},
			false,
		},
		{
			"MIX STATEMENTS",
			args{`SET data "Hello World"; GET data data1 data2 data3; DEL data data1 data2 data3; GET data; DEL data;`},
			&Ast{
				Statements: []*Statement{
					{
						SetStatement: &SetStatement{
							key: "data",
							val: "Hello World",
						},
						typ: SetType,
					},
					{
						GetStatement: &GetStatement{
							keys: []string{"data", "data1", "data2", "data3"},
						},
						typ: GetType,
					},
					{
						DeleteStatement: &DeleteStatement{
							keys: []string{"data", "data1", "data2", "data3"},
						},
						typ: DeleteType,
					},
					{
						GetStatement: &GetStatement{
							keys: []string{"data"},
						},
						typ: GetType,
					},
					{
						DeleteStatement: &DeleteStatement{
							keys: []string{"data"},
						},
						typ: DeleteType,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
