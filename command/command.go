package command

import (
	"errors"

	cafe "github.com/nevindra/community-bot/cafe"
)

type CommandDomain struct {
	cafeDomain *cafe.CafeDomain
}

func NewCommandDomain(cafeDomain *cafe.CafeDomain) (*CommandDomain, error) {
	if cafeDomain == nil {
		return nil, errors.New("cafeDomain is required")
	}

	return &CommandDomain{cafeDomain: cafeDomain}, nil
}
