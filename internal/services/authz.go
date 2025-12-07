package services

import "github.com/typewriterco/p402/internal/accesscontrol"

type AuthzSvc struct {
	ac accesscontrol.AC
}

func NewAuthz(ac accesscontrol.AC) *AuthzSvc {
	return &AuthzSvc{ac: ac}
}
