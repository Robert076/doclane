package repositories

import (
	"context"
	"database/sql"

	"github.com/Robert076/doclane/backend/models"
)

type StatsRepo struct {
	db *sql.DB
}

func NewStatsRepo(db *sql.DB) *StatsRepo {
	return &StatsRepo{db: db}
}

func (r *StatsRepo) GetStats(ctx context.Context) (*models.Stats, error) {
	stats := &models.Stats{}

	if err := r.scanRequestCounts(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanAvgCompletionTime(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanWeeklyStats(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanMonthlyStats(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanDailyStats(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanRequestsPerDepartment(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanRequestsPerLocality(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanUserStats(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanDepartmentCount(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanTemplateStats(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanMostUsedTemplates(ctx, stats); err != nil {
		return nil, err
	}
	if err := r.scanMemberStats(ctx, stats); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *StatsRepo) scanRequestCounts(ctx context.Context, stats *models.Stats) error {
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE is_closed = false AND is_cancelled = false) AS open,
			COUNT(*) FILTER (WHERE is_closed = true AND is_cancelled = false) AS archived,
			COUNT(*) FILTER (WHERE is_cancelled = true) AS cancelled
		FROM document_requests
	`).Scan(&stats.TotalOpenRequests, &stats.TotalArchivedRequests, &stats.TotalCancelledRequests)
	if err != nil {
		return err
	}

	total := stats.TotalOpenRequests + stats.TotalArchivedRequests + stats.TotalCancelledRequests
	if total > 0 {
		stats.CompletionRate = float64(stats.TotalArchivedRequests) / float64(total) * 100
		stats.CancellationRate = float64(stats.TotalCancelledRequests) / float64(total) * 100
	}
	return nil
}

func (r *StatsRepo) scanAvgCompletionTime(ctx context.Context, stats *models.Stats) error {
	var avgHours sql.NullFloat64
	err := r.db.QueryRowContext(ctx, `
		SELECT EXTRACT(EPOCH FROM AVG(closed_at - created_at)) / 3600
		FROM document_requests
		WHERE is_closed = true
		  AND is_cancelled = false
		  AND closed_at IS NOT NULL
	`).Scan(&avgHours)
	if err != nil {
		return err
	}
	if avgHours.Valid {
		stats.AvgCompletionHours = avgHours.Float64
	}
	return nil
}

func (r *StatsRepo) scanWeeklyStats(ctx context.Context, stats *models.Stats) error {
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') AS this_week,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '14 days' AND created_at < NOW() - INTERVAL '7 days') AS last_week
		FROM document_requests
	`).Scan(&stats.RequestsThisWeek, &stats.RequestsLastWeek)
	if err != nil {
		return err
	}
	if stats.RequestsLastWeek > 0 {
		stats.WeeklyChangePercent = float64(stats.RequestsThisWeek-stats.RequestsLastWeek) / float64(stats.RequestsLastWeek) * 100
	}
	return nil
}

func (r *StatsRepo) scanMonthlyStats(ctx context.Context, stats *models.Stats) error {
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE created_at >= DATE_TRUNC('month', NOW())) AS this_month,
			COUNT(*) FILTER (WHERE created_at >= DATE_TRUNC('month', NOW()) - INTERVAL '1 month' AND created_at < DATE_TRUNC('month', NOW())) AS last_month
		FROM document_requests
	`).Scan(&stats.RequestsThisMonth, &stats.RequestsLastMonth)
	if err != nil {
		return err
	}
	if stats.RequestsLastMonth > 0 {
		stats.MonthlyChangePercent = float64(stats.RequestsThisMonth-stats.RequestsLastMonth) / float64(stats.RequestsLastMonth) * 100
	}
	return nil
}

func (r *StatsRepo) scanDailyStats(ctx context.Context, stats *models.Stats) error {
	rows, err := r.db.QueryContext(ctx, `
		SELECT TO_CHAR(DATE_TRUNC('day', created_at), 'YYYY-MM-DD') AS date, COUNT(*) AS count
		FROM document_requests
		WHERE created_at >= NOW() - INTERVAL '7 days'
		GROUP BY DATE_TRUNC('day', created_at)
		ORDER BY DATE_TRUNC('day', created_at) ASC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var d models.DailyRequestStat
		if err := rows.Scan(&d.Date, &d.Count); err != nil {
			return err
		}
		stats.RequestsLast7Days = append(stats.RequestsLast7Days, d)
	}
	return rows.Err()
}

func (r *StatsRepo) scanRequestsPerDepartment(ctx context.Context, stats *models.Stats) error {
	rows, err := r.db.QueryContext(ctx, `
		SELECT d.name, COUNT(dr.id) AS count
		FROM departments d
		LEFT JOIN document_requests dr ON dr.department_id = d.id
			AND dr.is_closed = false
			AND dr.is_cancelled = false
		GROUP BY d.name
		ORDER BY count DESC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var d models.DepartmentStat
		if err := rows.Scan(&d.DepartmentName, &d.RequestCount); err != nil {
			return err
		}
		stats.RequestsPerDepartment = append(stats.RequestsPerDepartment, d)
	}
	return rows.Err()
}

func (r *StatsRepo) scanRequestsPerLocality(ctx context.Context, stats *models.Stats) error {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.locality, COUNT(dr.id) AS count
		FROM document_requests dr
		JOIN users u ON u.id = dr.assignee
		WHERE u.locality IS NOT NULL
		GROUP BY u.locality
		ORDER BY count DESC
		LIMIT 10
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var l models.LocalityStat
		if err := rows.Scan(&l.Locality, &l.RequestCount); err != nil {
			return err
		}
		stats.RequestsPerLocality = append(stats.RequestsPerLocality, l)
	}
	return rows.Err()
}

func (r *StatsRepo) scanUserStats(ctx context.Context, stats *models.Stats) error {
	return r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE department_id IS NOT NULL AND role = 'member') AS members,
			COUNT(*) FILTER (WHERE department_id IS NULL AND role != 'admin') AS citizens,
			COUNT(*) FILTER (WHERE is_active = true) AS active,
			COUNT(*) FILTER (WHERE is_active = false) AS deactivated
		FROM users
		WHERE role != 'admin'
	`).Scan(
		&stats.TotalUsers,
		&stats.TotalDepartmentMembers,
		&stats.TotalCitizens,
		&stats.TotalActiveUsers,
		&stats.TotalDeactivatedUsers,
	)
}

func (r *StatsRepo) scanDepartmentCount(ctx context.Context, stats *models.Stats) error {
	return r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM departments`).Scan(&stats.TotalDepartments)
}

func (r *StatsRepo) scanTemplateStats(ctx context.Context, stats *models.Stats) error {
	return r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE is_closed = false) AS active,
			COUNT(*) FILTER (WHERE is_closed = true) AS archived
		FROM document_request_templates
	`).Scan(&stats.TotalActiveTemplates, &stats.TotalArchivedTemplates)
}

func (r *StatsRepo) scanMostUsedTemplates(ctx context.Context, stats *models.Stats) error {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.title, COUNT(dr.id) AS count
		FROM document_request_templates t
		LEFT JOIN document_requests dr ON dr.template_id = t.id
		GROUP BY t.id, t.title
		ORDER BY count DESC
		LIMIT 5
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var t models.TemplateStat
		if err := rows.Scan(&t.TemplateID, &t.TemplateTitle, &t.RequestCount); err != nil {
			return err
		}
		stats.MostUsedTemplates = append(stats.MostUsedTemplates, t)
	}
	return rows.Err()
}

func (r *StatsRepo) scanMemberStats(ctx context.Context, stats *models.Stats) error {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			d.name AS department_name,
			COUNT(dr.id) FILTER (WHERE dr.claimed_by = u.id) AS total_claimed,
			COUNT(dr.id) FILTER (WHERE dr.claimed_by = u.id AND dr.is_closed = true) AS total_closed,
			COUNT(dr.id) FILTER (WHERE dr.claimed_by = u.id AND dr.is_closed = false AND dr.is_cancelled = false) AS total_pending,
			COALESCE(
				EXTRACT(EPOCH FROM AVG(
					CASE WHEN dr.claimed_by = u.id AND dr.is_closed = true AND dr.closed_at IS NOT NULL AND dr.claimed_at IS NOT NULL
					THEN dr.closed_at - dr.claimed_at
					END
				)) / 3600,
				0
			) AS avg_close_time_hours
		FROM users u
		JOIN departments d ON d.id = u.department_id
		LEFT JOIN document_requests dr ON dr.claimed_by = u.id
		WHERE u.role = 'member'
		  AND u.department_id IS NOT NULL
		GROUP BY u.id, u.first_name, u.last_name, d.name
		ORDER BY total_closed DESC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var m models.MemberStat
		if err := rows.Scan(
			&m.UserID,
			&m.FirstName,
			&m.LastName,
			&m.DepartmentName,
			&m.TotalClaimed,
			&m.TotalClosed,
			&m.TotalPending,
			&m.AvgCloseTimeHours,
		); err != nil {
			return err
		}
		stats.MemberStats = append(stats.MemberStats, m)
	}
	return rows.Err()
}
