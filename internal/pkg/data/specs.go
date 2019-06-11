package data

import (
	"strings"
)

// Sort contains sorting info
type Sort struct {
	Prop      string `json:"-"`
	Direction string `json:"-"`
}

// Order direction
const (
	OrderAsc  = "ASC"
	OrderDesc = "DESC"
	OrderNone = ""
)

// NewSort instantiates new sort instance
func NewSort(prop, dir string) *Sort {
	switch dir = strings.ToUpper(dir); dir {
	case OrderAsc:
		fallthrough
	case OrderDesc:
		return &Sort{
			Prop:      prop,
			Direction: dir,
		}
	case OrderNone:
		fallthrough
	default:
		return nil
	}
}

// Filter cuma berisi filter querystring
type Filter struct {
	Prop     string
	Operator string
	Value    interface{}
}

// NewFilter constructs new filter instance
func NewFilter(prop, op string, val interface{}) *Filter {
	return &Filter{
		Prop:     prop,
		Operator: op,
		Value:    val,
	}
}

// Spec isinya sorts dan filters
type Spec struct {
	Search  string
	Sorts   []*Sort
	Filters []*Filter
}

// NewSpec instantiates new Spec instance
/*func NewSpec(sorts []*Sort, filters []*Filter) *Spec {
	if sorts == nil {
		sorts = make([]*Sort, 0)
	}

	if filters == nil {
		filters = make([]*Filter, 0)
	}

	return &Spec{
		Sorts:   sorts,
		Filters: filters,
	}
}*/
