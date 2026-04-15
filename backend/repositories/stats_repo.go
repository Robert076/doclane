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

	// request counts
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE is_closed = false AND is_cancelled = false) AS open,
			COUNT(*) FILTER (WHERE is_closed = true AND is_cancelled = false) AS archived,
			COUNT(*) FILTER (WHERE is_cancelled = true) AS cancelled
		FROM document_requests
	`).Scan(&stats.TotalOpenRequests, &stats.TotalArchivedRequests, &stats.TotalCancelledRequests)
	if err != nil {
		return nil, err
	}

	// rates
	total := stats.TotalOpenRequests + stats.TotalArchivedRequests + stats.TotalCancelledRequests
	if total > 0 {
		stats.CompletionRate = float64(stats.TotalArchivedRequests) / float64(total) * 100
		stats.CancellationRate = float64(stats.TotalCancelledRequests) / float64(total) * 100
	}

	// avg completion time
	var avgHours sql.NullFloat64
	err = r.db.QueryRowContext(ctx, `
		SELECT EXTRACT(EPOCH FROM AVG(closed_at - created_at)) / 3600
		FROM document_requests
		WHERE is_closed = true
		  AND is_cancelled = false
		  AND closed_at IS NOT NULL
	`).Scan(&avgHours)
	if err != nil {
		return nil, err
	}
	if avgHours.Valid {
		stats.AvgCompletionHours = avgHours.Float64
	}

	// weekly requests
	err = r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') AS this_week,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '14 days' AND created_at < NOW() - INTERVAL '7 days') AS last_week
		FROM document_requests
	`).Scan(&stats.RequestsThisWeek, &stats.RequestsLastWeek)
	if err != nil {
		return nil, err
	}
	if stats.RequestsLastWeek > 0 {
		stats.WeeklyChangePercent = float64(stats.RequestsThisWeek-stats.RequestsLastWeek) / float64(stats.RequestsLastWeek) * 100
	}

	// monthly requests
	err = r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE created_at >= DATE_TRUNC('month', NOW())) AS this_month,
			COUNT(*) FILTER (WHERE created_at >= DATE_TRUNC('month', NOW()) - INTERVAL '1 month' AND created_at < DATE_TRUNC('month', NOW())) AS last_month
		FROM document_requests
	`).Scan(&stats.RequestsThisMonth, &stats.RequestsLastMonth)
	if err != nil {
		return nil, err
	}
	if stats.RequestsLastMonth > 0 {
		stats.MonthlyChangePercent = float64(stats.RequestsThisMonth-stats.RequestsLastMonth) / float64(stats.RequestsLastMonth) * 100
	}

	// requests last 7 days
	rows, err := r.db.QueryContext(ctx, `
		SELECT TO_CHAR(DATE_TRUNC('day', created_at), 'YYYY-MM-DD') AS date, COUNT(*) AS count
		FROM document_requests
		WHERE created_at >= NOW() - INTERVAL '7 days'
		GROUP BY DATE_TRUNC('day', created_at)
		ORDER BY DATE_TRUNC('day', created_at) ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d models.DailyRequestStat
		if err := rows.Scan(&d.Date, &d.Count); err != nil {
			return nil, err
		}
		stats.RequestsLast7Days = append(stats.RequestsLast7Days, d)
	}

	// requests per department
	rows, err = r.db.QueryContext(ctx, `
		SELECT d.name, COUNT(dr.id) AS count
		FROM departments d
		LEFT JOIN document_requests dr ON dr.department_id = d.id
			AND dr.is_closed = false
			AND dr.is_cancelled = false
		GROUP BY d.name
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d models.DepartmentStat
		if err := rows.Scan(&d.DepartmentName, &d.RequestCount); err != nil {
			return nil, err
		}
		stats.RequestsPerDepartment = append(stats.RequestsPerDepartment, d)
	}

	// user stats
	err = r.db.QueryRowContext(ctx, `
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
	if err != nil {
		return nil, err
	}

	// department count
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM departments`).Scan(&stats.TotalDepartments)
	if err != nil {
		return nil, err
	}

	// template stats
	err = r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE is_closed = false) AS active,
			COUNT(*) FILTER (WHERE is_closed = true) AS archived
		FROM document_request_templates
	`).Scan(&stats.TotalActiveTemplates, &stats.TotalArchivedTemplates)
	if err != nil {
		return nil, err
	}

	// most used templates top 5
	rows, err = r.db.QueryContext(ctx, `
		SELECT t.id, t.title, COUNT(dr.id) AS count
		FROM document_request_templates t
		LEFT JOIN document_requests dr ON dr.template_id = t.id
		GROUP BY t.id, t.title
		ORDER BY count DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t models.TemplateStat
		if err := rows.Scan(&t.TemplateID, &t.TemplateTitle, &t.RequestCount); err != nil {
			return nil, err
		}
		stats.MostUsedTemplates = append(stats.MostUsedTemplates, t)
	}

	return stats, nil
}
