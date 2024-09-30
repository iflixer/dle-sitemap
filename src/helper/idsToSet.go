package helper

import (
	"fmt"
)

func IDsToSet(ids []int) (set string) {
	if len(ids) == 0 {
		return
	}
	for _, id := range ids {
		set = fmt.Sprintf("%s/%d", set, id)
	}
	set = fmt.Sprintf("%s/", set)
	return
}
func IDsToSetString(ids []string) (set string) {
	if len(ids) == 0 {
		return
	}
	for _, id := range ids {
		set = fmt.Sprintf("%s/%s", set, id)
	}
	set = fmt.Sprintf("%s/", set)
	return
}
