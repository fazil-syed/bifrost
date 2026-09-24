package role

import (
	"context"
	"fmt"
	"time"

	"github.com/fazil-syed/bifrost/internal/application"
	"github.com/fazil-syed/bifrost/internal/team"
	"github.com/fazil-syed/bifrost/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type roleService struct {
	db                 *pgxpool.Pool
	applicationService application.ApplicationService
	userService        user.UserService
	teamService        team.TeamService
	membershipService  team.MembershipService
}

func NewRoleService(db *pgxpool.Pool, applicationService application.ApplicationService, userService user.UserService, teamService team.TeamService, membershipService team.MembershipService) RoleService {
	return &roleService{db: db, applicationService: applicationService, userService: userService, teamService: teamService, membershipService: membershipService}
}

func (s *roleService) authorizeManagement(ctx context.Context, actorUserID uuid.UUID, applicationID uuid.UUID) error {
	application, err := s.applicationService.GetByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("get application: %w", err)
	}

	if application.OwnerUserID != nil {
		if *application.OwnerUserID != actorUserID {
			return ErrRoleUnauthorized
		}
		return nil
	}

	if application.OwnerTeamID != nil {
		if _, err := s.teamService.GetByID(ctx, *application.OwnerTeamID); err != nil {
			return fmt.Errorf("get application owner team id: %w", err)
		}
		membership, err := s.membershipService.GetMember(ctx, *application.OwnerTeamID, actorUserID)
		if err != nil {
			return fmt.Errorf("get application owner team membership: %w", err)
		}

		if membership.MembershipType != team.MembershipOwner {
			return ErrRoleUnauthorized
		}
		return nil
	}

	return ErrRoleUnauthorized

}

func (s *roleService) Create(ctx context.Context, actorUserID uuid.UUID, applicationID uuid.UUID, name string) (*Role, error) {
	if err := s.authorizeManagement(ctx, actorUserID, applicationID); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	role, err := New(applicationID, name, now)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create role transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	repository := NewRoleRepository(tx)

	if err := repository.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create role transaction: %w", err)
	}

	return role, nil
}

func (s *roleService) GetByID(ctx context.Context, id uuid.UUID) (*Role, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})

	if err != nil {
		return nil, fmt.Errorf("begin get role transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	repository := NewRoleRepository(tx)
	return repository.GetByID(ctx, id)
}
func (s *roleService) ListByApplication(ctx context.Context, applicationID uuid.UUID) ([]*Role, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})

	if err != nil {
		return nil, fmt.Errorf("begin list roles by application transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	repository := NewRoleRepository(tx)
	return repository.ListByApplication(ctx, applicationID)
}

func (s *roleService) Update(ctx context.Context, actorUserID uuid.UUID, id uuid.UUID, name string) error {
	role, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.authorizeManagement(ctx, actorUserID, role.ApplicationID); err != nil {
		return err
	}

	now := time.Now().UTC()

	updatedRole, err := New(role.ApplicationID, name, now)

	if err != nil {
		return err
	}

	updatedRole.ID = role.ID
	updatedRole.CreatedAt = role.CreatedAt
	updatedRole.UpdatedAt = now

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin update role transaction: %w", err)
	}

	defer tx.Rollback(ctx)
	repository := NewRoleRepository(tx)
	if err := repository.Update(ctx, updatedRole); err != nil {
		return fmt.Errorf("update role: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit update role transaction: %w", err)
	}

	return nil
}

func (s *roleService) Delete(ctx context.Context, actorUserID uuid.UUID, id uuid.UUID) error {
	role, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.authorizeManagement(ctx, actorUserID, role.ApplicationID); err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin update role transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	repository := NewRoleRepository(tx)

	if err := repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit delete transaction: %w", err)
	}
	return nil
}
