package models

import "strconv"

type BaseModel struct {
	ID uint64
}

func (a BaseModel) GetStringId() string {
	return strconv.FormatUint(a.ID, 10)
}
