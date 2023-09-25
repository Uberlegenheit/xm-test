package smodels

import (
	"fmt"
	"strings"
)

type Company struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"            binding:"required"`
	Description string `json:"description"`
	Employees   uint64 `json:"employees_count"`
	Registered  bool   `json:"registered"`
	Type        string `json:"type"            binding:"required"`
}

func (c *Company) Validate() error {
	c.ID = ""

	c.Name = strings.Trim(c.Name, " ")
	if len(c.Name) < 5 || len(c.Name) > 15 {
		return fmt.Errorf("incorrect company name")
	}

	c.Description = strings.Trim(c.Description, " ")
	if len(c.Description) > 3000 {
		return fmt.Errorf("too big description (should be less than 3000)")
	}

	if c.Type == "" {
		return fmt.Errorf("type should be specified")
	}

	return nil
}
