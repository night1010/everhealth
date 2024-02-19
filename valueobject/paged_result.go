package valueobject

type PagedResult struct {
	Data         any
	CurrentPage  int
	CurrentItems int
	TotalItem    int
	TotalPage    int
}
