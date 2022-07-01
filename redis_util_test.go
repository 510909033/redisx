package redisx

import "testing"

func TestFloatToString(t *testing.T) {
	type args struct {
		val float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "name1", args: args{val: 1}, want: "1"},
		{name: "name12", args: args{val: 0}, want: "0"},
		{name: "name13", args: args{val: 0.1}, want: "0.1"},
		{name: "name14", args: args{val: 1.234}, want: "1.234"},
		{name: "name14-1", args: args{val: 1.2340}, want: "1.234"},
		{name: "name15", args: args{val: 10.2345}, want: "10.2345"},
		{name: "name16", args: args{val: 100.23456}, want: "100.23456"},
		{name: "name17", args: args{val: 100.234567}, want: "100.234567"},
		{name: "name18", args: args{val: 100.2345678}, want: "100.2345678"},
		{name: "name19", args: args{val: 100.23456789}, want: "100.23456789"},
		{name: "name20", args: args{val: 100.234567890}, want: "100.23456789"},
		{name: "name21", args: args{val: 100.23456789901}, want: "100.23456789901"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FloatToString(tt.args.val); got != tt.want {
				t.Errorf("FloatToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
