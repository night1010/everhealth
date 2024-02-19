package valueobject

import "fmt"

type Query struct {
	page         int
	limit        *int
	sortBy       *string
	order        Order
	conditions   []*Condition
	associations []*Association
	locked       bool
}

func NewQuery() *Query {
	query := &Query{}
	query.page = 1
	query.limit = nil
	query.sortBy = nil
	query.order = OrderAsc
	query.conditions = make([]*Condition, 0)
	query.associations = make([]*Association, 0)
	return query
}

func (q *Query) Condition(field string, operator Operator, value any) *Query {
	condition := NewCondition(field, operator, value)
	if condition == nil {
		return q
	}
	q.conditions = append(q.conditions, condition)
	return q
}

func (q *Query) GetConditions() []*Condition {
	return q.conditions
}

func (q *Query) GetConditionValue(field string) any {
	for _, condition := range q.conditions {
		if condition.Field == field {
			return condition.Value
		}
	}
	return nil
}

func (q *Query) WithPage(page int) *Query {
	q.page = page
	return q
}

func (q *Query) GetPage() int {
	return q.page
}

func (q *Query) WithLimit(limit int) *Query {
	q.limit = &limit
	return q
}

func (q *Query) GetLimit() *int {
	return q.limit
}

func (q *Query) WithSortBy(field string) *Query {
	q.sortBy = &field
	return q
}

func (q *Query) WithOrder(order Order) *Query {
	q.order = order
	return q
}

func (q *Query) GetOrder() string {
	if q.sortBy == nil {
		return ""
	}
	return fmt.Sprintf("%s %s", *q.sortBy, q.order)
}

func (q *Query) WithJoin(entity string) *Query {
	q.associations = append(q.associations, &Association{
		Type:   AssociationTypeJoin,
		Entity: entity,
	})
	return q
}

func (q *Query) WithPreload(entity string) *Query {
	q.associations = append(q.associations, &Association{
		Type:   AssociationTypePreload,
		Entity: entity,
	})
	return q
}

func (q *Query) GetAssociations() []*Association {
	return q.associations
}

func (q *Query) Lock() *Query {
	q.locked = true
	return q
}

func (q *Query) IsLocked() bool {
	return q.locked
}
