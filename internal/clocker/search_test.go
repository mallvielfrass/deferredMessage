package clocker

import (
	"testing"
	"time"
)

func Test_binarySearchMsgWithTime(t *testing.T) {
	type args struct {
		messages []Message
		time     time.Time
	}
	tests := []struct {
		name        string
		args        args
		wantLength  int
		afterLength int
	}{
		{
			name: "binarySearchMsgWithTime",
			args: args{
				messages: []Message{
					{
						Message: "hello",
						Time:    time.Now().Add(-1 * time.Second),
					},
					{
						Message: "hello",
						Time:    time.Now().Add(5 * time.Second),
					},
					{
						Message: "hello",
						Time:    time.Now().Add(5 * time.Second),
					},
				},
				time: time.Now(),
			},
			wantLength:  1,
			afterLength: 2,
		},
		{
			name: "binarySearchMsgWithTimeAfter",
			args: args{
				messages: []Message{
					{
						Message: "hello",
						Time:    time.Now().Add(5 * time.Second),
					},
					{
						Message: "hello",
						Time:    time.Now().Add(5 * time.Second),
					},
				},
				time: time.Now(),
			},
			wantLength:  0,
			afterLength: 2,
		},
		{
			name: "binarySearchMsgWithTimeBefore",
			args: args{
				messages: []Message{
					{
						Message: "hello",
						Time:    time.Now().Add(-1 * time.Second),
					},
					{
						Message: "hello",
						Time:    time.Now().Add(-1 * time.Second),
					},
				},
				time: time.Now(),
			},
			wantLength:  2,
			afterLength: 0,
		},
		{
			name: "binarySearchMsgEmpty",
			args: args{
				messages: []Message{},
				time:     time.Now(),
			},
			wantLength:  0,
			afterLength: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, index := binarySearchMsgWithTime(tt.args.messages, tt.args.time)
			if len(got) != tt.wantLength {
				t.Errorf("binarySearchMsgWithTime() = %v, want %v", got, tt.wantLength)
			}
			if index != tt.wantLength {
				t.Errorf("index = %v, want %v", index, tt.wantLength)
			}
			after := tt.args.messages[index:]
			if len(after) != tt.afterLength {
				t.Errorf("after = %v, want %v", after, tt.afterLength)
			}
		})
	}
}
