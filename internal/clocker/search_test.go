package clocker

import (
	"deferredMessage/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_binarySearchMsgWithTime(t *testing.T) {
	type args struct {
		messages []models.Message
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
				messages: []models.Message{
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
				messages: []models.Message{
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
				messages: []models.Message{
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
				messages: []models.Message{},
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

func Test_findAndRemoveMessage(t *testing.T) {
	type args struct {
		messages []models.Message
		msg      models.Message
	}

	tests := []struct {
		name string
		args args
		want []models.Message
	}{
		{
			name: "emptyArray",
			args: args{
				messages: []models.Message{},
				msg:      models.Message{},
			},
			want: []models.Message{},
		},
		{
			name: "msgNotInArray",
			args: args{
				messages: []models.Message{
					{
						Message: "hello1",
						Time:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
						Id:      "100000",
					},
				},
				msg: models.Message{
					Message: "hello2",
					Time:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
					Id:      "100001",
				},
			},
			want: []models.Message{
				{
					Message: "hello1",
					Time:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
					Id:      "100000",
				},
			},
		},
		{
			name: "removeMsgFromArray",
			args: args{
				messages: []models.Message{
					{
						Message: "hello1",
						Time:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
						Id:      "100000",
					},
					{
						Message: "hello2",
						Time:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
						Id:      "100001",
					},
				},
				msg: models.Message{
					Message: "hello2",
					Time:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
					Id:      "100001",
				},
			},
			want: []models.Message{
				{
					Message: "hello1",
					Time:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
					Id:      "100000",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findAndRemoveMessage(tt.args.messages, tt.args.msg)
			if !assert.ElementsMatch(t, tt.want, got) {
				t.Errorf("findAndRemovemodels.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}
