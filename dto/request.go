package dto

type RequestUri struct {
	Id int64 `uri:"id" binding:"required,numeric"`
}
