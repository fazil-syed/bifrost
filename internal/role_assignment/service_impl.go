package roleassignment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fazil-syed/bifrost/internal/application"
	"github.com/fazil-syed/bifrost/internal/role"
	"github.com/fazil-syed/bifrost/internal/team"
	"github.com/fazil-syed/bifrost/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type roleAssignmentService struct {
	db                 *pgxpool.Pool
	roleService        role.RoleService
	applicationService application.ApplicationService
	userService        user.UserService
	teamService        team.TeamService
	membershipService  team.MembershipService
}

func NewRoleAssignmentService(db *pgxpool.Pool, roleService role.RoleService, applicationService application.ApplicationService, userService user.UserService, teamService team.TeamService, membershipService team.MembershipService) RoleAssignmentService {
	return &roleAssignmentService{
		db:                 db,
		roleService:        roleService,
		applicationService: applicationService,
		userService:        userService,
		teamService:        teamService,
		membershipService:  membershipService,
	}
}

func (s *roleAssignmentService) AssignToUser(ctx context.Context, roleID uuid.UUID, userID uuid.UUID) (*RoleAssignment, error) {
	role, err := s.roleService.GetByID(ctx, roleID)

	if err != nil {
		return nil, err
	}
	if _, err := s.userService.GetByID(ctx, userID); err != nil {
		return nil, err
	}
	if err := s.authorizeRoleManagement(ctx, role.ApplicationID); err != nil {
		return nil, err
	}
	now := time.Now().UTC()

	assignment, err := New(roleID, &userID, nil, now)

	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("begin assign role to user transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewRoleAssignmentRepository(tx)

	if err := repository.Create(ctx, assignment); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit assign role to user transaction: %w", err)
	}

	return assignment, nil
}
func (s *roleAssignmentService) AssignToTeam(ctx context.Context, roleID uuid.UUID, teamID uuid.UUID) (*RoleAssignment, error) {
	role, err := s.roleService.GetByID(ctx, roleID)

	if err != nil {
		return nil, err
	}
	if _, err := s.teamService.GetByID(ctx, teamID); err != nil {
		return nil, err
	}
	if err := s.authorizeRoleManagement(ctx, role.ApplicationID); err != nil {
		return nil, err
	}
	now := time.Now().UTC()

	assignment, err := New(roleID, nil, &teamID, now)

	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("begin assign role to team transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewRoleAssignmentRepository(tx)

	if err := repository.Create(ctx, assignment); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit assign role to team transaction: %w", err)
	}

	return assignment, nil
}

func (s *roleAssignmentService) GetByID(ctx context.Context, id uuid.UUID) (*RoleAssignment, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, fmt.Errorf("begin get role assignment by id transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewRoleAssignmentRepository(tx)

	return repository.GetByID(ctx, id)
}
func (s *roleAssignmentService) ListByRole(ctx context.Context, roleID uuid.UUID) ([]*RoleAssignment, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, fmt.Errorf("begin list role assignments by role_id transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewRoleAssignmentRepository(tx)

	return repository.ListByRole(ctx, roleID)
}
func (s *roleAssignmentService) ListByUser(ctx context.Context, userID uuid.UUID) ([]*RoleAssignment, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, fmt.Errorf("begin list role assignments by user_id transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewRoleAssignmentRepository(tx)

	return repository.ListByRole(ctx, userID)
}
func (s *roleAssignmentService) ListByTeam(ctx context.Context, teamID uuid.UUID) ([]*RoleAssignment, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, fmt.Errorf("begin list role assignments by team_id transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewRoleAssignmentRepository(tx)

	return repository.ListByRole(ctx, teamID)
}

func (s *roleAssignmentService) Remove(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)

	if err != nil {
		return fmt.Errorf("begin remove role assignment transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewRoleAssignmentRepository(tx)
	if err := repository.Delete(ctx, id); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit remove role assignment transaction: %w", err)
	}
	return nil
}

func (s *roleAssignmentService) authorizeRoleManagement(ctx context.Context, applicationID uuid.UUID) error {
	application, err := s.applicationService.GetByID(ctx, applicationID)
	if err != nil {
		return err
	}
	actorUserID, err := actorUserIDFromContext(ctx)

	if err != nil {
		return ErrRoleAssignmentUnauthorized
	}
	if application.OwnerUserID != nil {
		if *application.OwnerUserID == actorUserID {
			return nil
		}
		return ErrRoleAssignmentUnauthorized
	}
	if application.OwnerTeamID != nil {
		membership, err := s.membershipService.GetMember(ctx, *application.OwnerTeamID, actorUserID)
		if err != nil {
			if errors.Is(err, team.ErrTeamMembershipNotFound) {
				return ErrRoleAssignmentUnauthorized
			}
			return err
		}
		if membership.MembershipType == team.MembershipOwner {
			return nil
		}
	}
	return ErrRoleAssignmentUnauthorized
}

func actorUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	// Temporary key before implementing auth context

	value := ctx.Value(actorUserIDContextKey{})

	if value == nil {
		return uuid.Nil, fmt.Errorf("actor user ID missing from context")
	}

	userID, ok := value.(uuid.UUID)

	if !ok {
		return uuid.Nil, fmt.Errorf("invalid actor user ID in context")
	}

	return userID, nil
}

type actorUserIDContextKey struct{}
