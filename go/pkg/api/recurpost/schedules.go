package recurpost

import "context"

// SchedulesService manages scheduling endpoints.
type SchedulesService struct{ c *Client }

// CreateSchedule schedules an existing post ID for a specific time across accounts.
func (s *SchedulesService) CreateSchedule(ctx context.Context, postID string, at string) (*Post, error) {
	// POST /schedules
	return nil, ErrNotImplemented
}

// DeleteSchedule removes a schedule by id.
func (s *SchedulesService) DeleteSchedule(ctx context.Context, scheduleID string) error {
	// DELETE /schedules/{id}
	return ErrNotImplemented
}

// ListSchedules lists schedules.
func (s *SchedulesService) ListSchedules(ctx context.Context, p Pagination) ([]Post, error) {
	// GET /schedules
	return nil, ErrNotImplemented
}
