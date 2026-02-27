package services

import (
	"context"
	"fmt"

	"github.com/richardbowden/degrees/internal/accesscontrol"
	"github.com/richardbowden/degrees/internal/problems"
)

const (
	SystemObject = "system:main" // The singleton system object for global roles
)

type AuthzSvc struct {
	ac accesscontrol.AC
}

func NewAuthz(ac accesscontrol.AC) *AuthzSvc {
	return &AuthzSvc{ac: ac}
}

// SetUserAsSysop grants a user the sysop role (highest system privileges)
func (az *AuthzSvc) SetUserAsSysop(ctx context.Context, userID int64) error {
	user := fmt.Sprintf("user:%d", userID)
	err := az.ac.WriteRelationship(ctx, user, "sysop", SystemObject)
	if err != nil {
		return problems.New(problems.Internal, "failed to grant sysop role", err)
	}
	return nil
}

// RevokeUserSysop removes sysop role from a user
func (az *AuthzSvc) RevokeUserSysop(ctx context.Context, userID int64) error {
	user := fmt.Sprintf("user:%d", userID)
	err := az.ac.DeleteRelationship(ctx, user, "sysop", SystemObject)
	if err != nil {
		return problems.New(problems.Internal, "failed to revoke sysop role", err)
	}
	return nil
}

// SetUserAsSystemAdmin grants a user the admin role (can manage system but not other sysops)
func (az *AuthzSvc) SetUserAsSystemAdmin(ctx context.Context, userID int64) error {
	user := fmt.Sprintf("user:%d", userID)
	err := az.ac.WriteRelationship(ctx, user, "admin", SystemObject)
	if err != nil {
		return problems.New(problems.Internal, "failed to grant admin role", err)
	}
	return nil
}

// RevokeUserSystemAdmin removes admin role from a user
func (az *AuthzSvc) RevokeUserSystemAdmin(ctx context.Context, userID int64) error {
	user := fmt.Sprintf("user:%d", userID)
	err := az.ac.DeleteRelationship(ctx, user, "admin", SystemObject)
	if err != nil {
		return problems.New(problems.Internal, "failed to revoke admin role", err)
	}
	return nil
}

// IsSysop checks if a user has sysop privileges
func (az *AuthzSvc) IsSysop(ctx context.Context, userID int64) (bool, error) {
	user := fmt.Sprintf("user:%d", userID)
	allowed, err := az.ac.Check(ctx, user, "sysop", SystemObject)
	if err != nil {
		return false, problems.New(problems.Internal, "failed to check sysop permission", err)
	}
	return allowed, nil
}

// ListSysopUserIDs returns a set of user IDs that have sysop privileges, keyed from FGA
func (az *AuthzSvc) ListSysopUserIDs(ctx context.Context) (map[int64]bool, error) {
	users, err := az.ac.ListSysopUsers(ctx)
	if err != nil {
		return nil, problems.New(problems.Internal, "failed to list sysop users", err)
	}
	ids := make(map[int64]bool, len(users))
	for _, u := range users {
		var id int64
		fmt.Sscanf(u, "user:%d", &id)
		if id > 0 {
			ids[id] = true
		}
	}
	return ids, nil
}

// IsSystemAdmin checks if a user has system admin privileges (sysop or admin)
func (az *AuthzSvc) IsSystemAdmin(ctx context.Context, userID int64) (bool, error) {
	user := fmt.Sprintf("user:%d", userID)
	allowed, err := az.ac.Check(ctx, user, "admin", SystemObject)
	if err != nil {
		return false, problems.New(problems.Internal, "failed to check admin permission", err)
	}
	return allowed, nil
}
