package services

import "github.com/typewriterco/p402/internal/accesscontrol"

type Authz struct {
	ac accesscontrol.AC
}

func NewAuthz(ac accesscontrol.AC) *Authz {
	return &Authz{ac: ac}
}
