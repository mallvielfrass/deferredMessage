package clocker

import "time"

func binarySearchMsgWithTime(messages []Message, time time.Time) ([]Message, int) {
	lowIndex := 0
	high := len(messages) - 1
	for lowIndex <= high {
		mid := (lowIndex + high) / 2
		if messages[mid].Time.After(time) {
			high = mid - 1
		} else {
			lowIndex = mid + 1
		}
	}
	//fmt.Printf("lowIndex: %v, high: %v\n", lowIndex, high)
	return messages[:lowIndex], lowIndex
}
func findAndRemoveMessage(messages []Message, msg Message) []Message {
	var newArray []Message
	for i := 0; i < len(messages); i++ {
		if messages[i].Id == msg.Id {
			newArray = append(messages[:i], messages[i+1:]...)
			return newArray
		}
	}
	return messages
}
