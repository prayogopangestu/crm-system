package search

import (
	"github.com/prayogopangestu/crm-system/backend/internal/modules/contact"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/deal"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/task"
)

type Result struct {
	Contacts []contact.Contact `json:"contacts"`
	Tasks    []task.Task       `json:"tasks"`
	Deals    []deal.Deal       `json:"deals"`
}
