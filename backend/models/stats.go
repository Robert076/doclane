package models

type DepartmentStat struct {
	DepartmentName string `json:"department_name"`
	RequestCount   int    `json:"request_count"`
}

type DailyRequestStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type TemplateStat struct {
	TemplateID    int    `json:"template_id"`
	TemplateTitle string `json:"template_title"`
	RequestCount  int    `json:"request_count"`
}

type MemberStat struct {
	UserID            int     `json:"user_id"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	DepartmentName    string  `json:"department_name"`
	TotalClaimed      int     `json:"total_claimed"`
	TotalClosed       int     `json:"total_closed"`
	TotalPending      int     `json:"total_pending"`
	AvgCloseTimeHours float64 `json:"avg_close_time_hours"`
}

type LocalityStat struct {
	Locality     string `json:"locality"`
	RequestCount int    `json:"request_count"`
}

type Stats struct {
	TotalOpenRequests      int     `json:"total_open_requests"`
	TotalArchivedRequests  int     `json:"total_archived_requests"`
	TotalCancelledRequests int     `json:"total_cancelled_requests"`
	TotalOverdueRequests   int     `json:"total_overdue_requests"`
	CompletionRate         float64 `json:"completion_rate"`
	CancellationRate       float64 `json:"cancellation_rate"`
	AvgCompletionHours     float64 `json:"avg_completion_hours"`

	RequestsThisWeek     int     `json:"requests_this_week"`
	RequestsLastWeek     int     `json:"requests_last_week"`
	WeeklyChangePercent  float64 `json:"weekly_change_percent"`
	RequestsThisMonth    int     `json:"requests_this_month"`
	RequestsLastMonth    int     `json:"requests_last_month"`
	MonthlyChangePercent float64 `json:"monthly_change_percent"`

	RequestsLast7Days     []DailyRequestStat `json:"requests_last_7_days"`
	RequestsPerDepartment []DepartmentStat   `json:"requests_per_department"`
	RequestsPerLocality   []LocalityStat     `json:"requests_per_locality"`

	TotalDepartments       int `json:"total_departments"`
	TotalDepartmentMembers int `json:"total_department_members"`
	TotalCitizens          int `json:"total_citizens"`
	TotalUsers             int `json:"total_users"`
	TotalActiveUsers       int `json:"total_active_users"`
	TotalDeactivatedUsers  int `json:"total_deactivated_users"`

	TotalActiveTemplates   int            `json:"total_active_templates"`
	TotalArchivedTemplates int            `json:"total_archived_templates"`
	MostUsedTemplates      []TemplateStat `json:"most_used_templates"`

	MemberStats []MemberStat `json:"member_stats"`
}
