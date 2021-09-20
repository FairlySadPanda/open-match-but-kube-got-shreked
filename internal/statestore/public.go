// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package statestore

import (
	"context"

	"github.com/FairlySadPanda/open-match-but-kube-got-shreked/internal/config"
	"github.com/FairlySadPanda/open-match-but-kube-got-shreked/internal/telemetry"
	"github.com/FairlySadPanda/open-match-but-kube-got-shreked/pkg/pb"
)

// Service is a generic interface for talking to a storage backend.
type Service interface {
	// HealthCheck indicates if the database is reachable.
	HealthCheck(ctx context.Context) error

	// Closes the connection to the underlying storage.
	Close() error

	// Ticket

	// CreateTicket creates a new Ticket in the state storage. If the id already exists, it will be overwritten.
	CreateTicket(ctx context.Context, ticket *pb.Ticket) error

	// GetTicket gets the Ticket with the specified id from state storage.
	// This method fails if the Ticket does not exist.
	GetTicket(ctx context.Context, id string) (*pb.Ticket, error)

	// DeleteTicket removes the Ticket with the specified id from state storage.
	// This method succeeds if the Ticket does not exist.
	DeleteTicket(ctx context.Context, id string) error

	// IndexTicket adds the ticket to the index.
	IndexTicket(ctx context.Context, ticket *pb.Ticket) error

	// DeindexTicket removes specified ticket from the index. The Ticket continues to exist.
	DeindexTicket(ctx context.Context, id string) error

	// GetIndexedIDSet returns the ids of all tickets currently indexed.
	GetIndexedIDSet(ctx context.Context) (map[string]struct{}, error)

	// GetTickets returns multiple tickets from storage.
	// Missing tickets are silently ignored.
	GetTickets(ctx context.Context, ids []string) ([]*pb.Ticket, error)

	// UpdateAssignments update using the request's specified tickets with assignments.
	UpdateAssignments(ctx context.Context, req *pb.AssignTicketsRequest) (*pb.AssignTicketsResponse, []*pb.Ticket, error)

	// GetAssignments returns the assignment associated with the input ticket id.
	GetAssignments(ctx context.Context, id string, callback func(*pb.Assignment) error) error

	// AddTicketsToPendingRelease appends new proposed tickets to the proposed sorted set with current timestamp.
	AddTicketsToPendingRelease(ctx context.Context, ids []string) error

	// DeleteTicketsFromPendingRelease deletes tickets from the proposed sorted set.
	DeleteTicketsFromPendingRelease(ctx context.Context, ids []string) error

	// ReleaseAllTickets releases all pending tickets back to active.
	ReleaseAllTickets(ctx context.Context) error

	// Backfill

	// CreateBackfill creates a new Backfill in the state storage if one doesn't exist.
	// The xids algorithm used to create the ids ensures that they are unique with no system wide synchronization.
	// Calling clients are forbidden from choosing an id during create. So no conflicts will occur.
	CreateBackfill(ctx context.Context, backfill *pb.Backfill, ticketIDs []string) error

	// GetBackfill gets the Backfill with the specified id from state storage.
	// This method fails if the Backfill does not exist.
	// Returns the Backfill and associated ticketIDs if they exist.
	GetBackfill(ctx context.Context, id string) (*pb.Backfill, []string, error)

	// GetBackfills returns multiple backfills from storage
	GetBackfills(ctx context.Context, ids []string) ([]*pb.Backfill, error)

	// DeleteBackfill removes the Backfill with the specified id from state storage.
	// This method succeeds if the Backfill does not exist.
	DeleteBackfill(ctx context.Context, id string) error

	// DeleteBackfillCompletely performs a set of operations to remove backfill and all related entities.
	DeleteBackfillCompletely(ctx context.Context, id string) error

	// UpdateBackfill updates an existing Backfill with a new data. ticketIDs can be nil.
	UpdateBackfill(ctx context.Context, backfill *pb.Backfill, ticketIDs []string) error

	// NewMutex returns an interface of a new distributed mutex with given name
	NewMutex(key string) RedisLocker

	// CleanupBackfills removes expired backfills
	CleanupBackfills(ctx context.Context) error

	// UpdateAcknowledgmentTimestamp updates Backfill's last acknowledged time
	UpdateAcknowledgmentTimestamp(ctx context.Context, id string) error

	// GetExpiredBackfillIDs gets all backfill IDs which are expired
	GetExpiredBackfillIDs(ctx context.Context) ([]string, error)

	// IndexBackfill adds the backfill to the index.
	IndexBackfill(ctx context.Context, backfill *pb.Backfill) error

	// DeindexBackfill removes specified Backfill ID from the index. The Backfill continues to exist.
	DeindexBackfill(ctx context.Context, id string) error

	// GetIndexedBackfills returns a map containing the IDs and
	// the Generation number of the backfills currently indexed.
	GetIndexedBackfills(ctx context.Context) (map[string]int, error)
}

// New creates a Service based on the configuration.
func New(cfg config.View) Service {
	s := newRedis(cfg)
	if cfg.GetBool(telemetry.ConfigNameEnableMetrics) {
		return &instrumentedService{
			s: s,
		}
	}
	return s
}

// RedisLocker provides methods to use distributed locks against redis
type RedisLocker interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) (bool, error)
}
