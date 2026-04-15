export interface DailyRequestStat {
        date: string;
        count: number;
}

export interface DepartmentStat {
        department_name: string;
        request_count: number;
}

export interface TemplateStat {
        template_id: number;
        template_title: string;
        request_count: number;
}

export interface Stats {
        total_open_requests: number;
        total_archived_requests: number;
        total_cancelled_requests: number;
        completion_rate: number;
        cancellation_rate: number;
        avg_completion_hours: number;

        requests_this_week: number;
        requests_last_week: number;
        weekly_change_percent: number;
        requests_this_month: number;
        requests_last_month: number;
        monthly_change_percent: number;

        requests_last_7_days: DailyRequestStat[];
        requests_per_department: DepartmentStat[];
        most_used_templates: TemplateStat[];

        total_departments: number;
        total_department_members: number;
        total_users: number;
        total_citizens: number;
        total_active_users: number;
        total_deactivated_users: number;
        total_active_templates: number;
        total_archived_templates: number;
}
