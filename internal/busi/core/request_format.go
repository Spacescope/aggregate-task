package core

import (
	"errors"
)

type Walk struct {
	MinHeight      uint64   `form:"from" json:"from"`
	MaxHeight      uint64   `form:"to" json:"to"`
	Force          bool     `form:"force" json:"force" desc:"force to walk"`
	Task           string   `form:"-" json:"-"`
	DependentTasks []string `form:"-" json:"-"`
}

func (r *Walk) Validate() error {
	if r.MinHeight > r.MaxHeight {
		return errors.New("'from' should less or equal than 'to'")
	}

	return nil
}
