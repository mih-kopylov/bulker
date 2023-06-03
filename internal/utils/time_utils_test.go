package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestAgeToTime(t *testing.T) {
	type args struct {
		time  time.Time
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "3h",
			args: args{
				time:  time.Date(2023, 1, 1, 6, 3, 3, 0, time.UTC),
				value: "3h",
			},
			want:    time.Date(2023, 1, 1, 3, 3, 3, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "3d",
			args: args{
				time:  time.Date(2023, 1, 6, 3, 3, 3, 0, time.UTC),
				value: "3d",
			},
			want:    time.Date(2023, 1, 3, 3, 3, 3, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "3w",
			args: args{
				time:  time.Date(2023, 1, 26, 3, 3, 3, 0, time.UTC),
				value: "3w",
			},
			want:    time.Date(2023, 1, 5, 3, 3, 3, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "3h3d3w",
			args: args{
				time:  time.Date(2023, 1, 26, 6, 3, 3, 0, time.UTC),
				value: "3h3d3w",
			},
			want:    time.Date(2023, 1, 2, 3, 3, 3, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "3w3h3d",
			args: args{
				time:  time.Date(2023, 1, 26, 6, 3, 3, 0, time.UTC),
				value: "3w3h3d",
			},
			want:    time.Date(2023, 1, 2, 3, 3, 3, 0, time.UTC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				clock := &FixedClock{tt.args.time}
				got, err := AgeToTime(clock, tt.args.value)
				if (err != nil) != tt.wantErr {
					t.Errorf("AgeToTime() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AgeToTime() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
