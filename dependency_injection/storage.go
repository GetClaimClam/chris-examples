package storage

import (
	"context"

	"github.com/chr1sbest/api.mobl.ai/internal/storage/model"
)

// Storage defines the interface for stateful operations.
type Storage interface {
	// User stores core information about the user.
	CreateUser(ctx context.Context, name, email string) (*model.User, error)
	ReadUser(ctx context.Context, email string) (*model.User, error)
	UpdateLastLoggedIn(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, email string) error
	GetUserByUserID(ctx context.Context, userID string) (*model.User, error)

	// Routines are stored sets of exercises.
	CreateRoutine(ctx context.Context, email string, routine *model.Routine) (*model.Routine, error)
	ReadRoutine(ctx context.Context, email, id string) (*model.Routine, error)
	UpdateRoutine(ctx context.Context, email, id string, routine *model.Routine) error
	DeleteRoutine(ctx context.Context, email, id string) error
	ListRoutines(ctx context.Context, email string) ([]model.Routine, error)

	// ExerciseStreakRecords are records of what routines a user did and when.
	CreateExerciseStreakRecord(ctx context.Context, email string, routine *model.Routine) (*model.ExerciseStreakRecord, error)

	// Friends
	CreateFriendRequest(ctx context.Context, frq *model.FriendRequest) error
	ReadFriendRequest(ctx context.Context, email, friendRequestID string) (*model.FriendRequest, error)
	AcceptFriendRequest(ctx context.Context, email, friendRequestID string) error
	RejectFriendRequest(ctx context.Context, email, friendRequestID string) error
	DeleteFriendRequest(ctx context.Context, email, friendRequestID string) error
	ListFriendRequests(ctx context.Context, email string) ([]*model.FriendRequest, error)
	ListOutgoingFriendRequests(ctx context.Context, email string) ([]*model.FriendRequest, error)

	CreateFriend(ctx context.Context, email string, friend *model.Friend) (*model.Friend, error)
	ReadFriend(ctx context.Context, email string, friendID string) (*model.Friend, error)
	UpdateFriend(ctx context.Context, email string, friend *model.Friend) error
	DeleteFriend(ctx context.Context, email string, friendID string) error
	ListFriends(ctx context.Context, email string) ([]*model.Friend, error)

	// Miscellaneous
	CreateEmailRecord(ctx context.Context, email, userAgent string) error
	ReadRateLimit(ctx context.Context, key string) (int, int64, error)
	WriteRateLimit(ctx context.Context, key string, requestCount int, ttl int64) error

	// Purchases
	CreatePurchase(ctx context.Context, purchase *model.Purchase) error
	UpdatePurchase(ctx context.Context, purchase *model.Purchase) error
	GetMostRecentPurchaseByUserID(ctx context.Context, userID string) (*model.Purchase, error)

	// RoutineShares
	CreateRoutineShare(ctx context.Context, frq *model.RoutineShare) (*model.RoutineShare, error)
	ReadRoutineShare(ctx context.Context, email, routineShareID string) (*model.RoutineShare, error)
	AcceptRoutineShare(ctx context.Context, email, routineShareID string) error
	RejectRoutineShare(ctx context.Context, email, routineShareID string) error
	DeleteRoutineShare(ctx context.Context, email, routineShareID string) error
	ListRoutineShares(ctx context.Context, email string) ([]*model.RoutineShare, error)
	ListOutgoingRoutineShares(ctx context.Context, email string) ([]*model.RoutineShare, error)

	// Badges
	IncrementBadge(ctx context.Context, email, badge string) (*model.Badge, error)
	ListBadges(ctx context.Context, email string) ([]*model.Badge, error)
	DeleteBadges(ctx context.Context, email string) error
}
