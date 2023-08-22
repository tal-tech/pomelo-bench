package main

import "testing"

func Test_findMaxNumber(t *testing.T) {
	type args struct {
		target    int
		maxFactor int
	}
	tests := []struct {
		name       string
		args       args
		wantFactor int
		wantNumber int
	}{
		{
			name: "t1",
			args: args{
				target:    1600,
				maxFactor: 1000,
			},
			wantFactor: 800,
			wantNumber: 2,
		},
		{
			name: "t2",
			args: args{
				target:    1603,
				maxFactor: 1000,
			},
			wantFactor: 229,
			wantNumber: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFactor, gotNumber := findMaxNumber(tt.args.target, tt.args.maxFactor)
			if gotFactor != tt.wantFactor {
				t.Errorf("findMaxNumber() gotFactor = %v, want %v", gotFactor, tt.wantFactor)
			}
			if gotNumber != tt.wantNumber {
				t.Errorf("findMaxNumber() gotNumber = %v, want %v", gotNumber, tt.wantNumber)
			}
		})
	}
}
