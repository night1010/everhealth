package valueobject

type AssociationType int

type Association struct {
	Type   AssociationType
	Entity string
}

type Order string

type Operator string

const (
	AssociationTypeJoin AssociationType = iota + 1
	AssociationTypePreload
)

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

const (
	Equal            Operator = "="
	NotEqual                  = "!="
	LessThan                  = "<"
	GreaterThan               = ">"
	LessThanEqual             = "<="
	GreaterThanEqual          = ">="
	Is                        = "IS"
	In                        = "IN"
	Not                       = "NOT"
	Like                      = "LIKE"
	NotLike                   = "NOT LIKE"
	ILike                     = "ILIKE"
	NotILike                  = "NOT ILIKE"
)
