export interface APIResponse<T = unknown> {
        success: boolean;
        message: string;
        error?: string;
        data?: T;
}

export interface User {
        id: number;
        email: string;
        role: UserRole;
        department_id?: number | null;
        is_active: boolean;
        created_at: string;
        updated_at: string;
        phone?: string;
        locality?: string;
        street?: string;
        last_notified?: string | null;
        first_name: string;
        last_name: string;
}

export interface Department {
        id: number;
        name: string;
        created_at: string;
        updated_at: string;
}

export interface Request {
        id: number;
        assignee: number;
        assignee_email: string;
        assignee_first_name: string;
        assignee_last_name: string;
        department_id: number;
        department_name: string;
        title: string;
        description?: string | null;
        due_date?: string | null;
        is_recurring: boolean;
        recurrence_cron?: string | null;
        is_scheduled: boolean;
        scheduled_for?: string | null;
        next_due_at?: string | null;
        last_uploaded_at?: string | null;
        is_cancelled: boolean;
        is_closed: boolean;
        template_id: number;
        status: RequestStatus;
        created_at: string;
        updated_at: string;
        expected_documents: ExpectedDocument[];
}

export interface Template {
        id: number;
        title: string;
        description?: string | null;
        department_id: number;
        department_name: string;
        is_recurring: boolean;
        recurrence_cron?: string | null;
        created_by: number;
        author_first_name?: string | null;
        author_last_name?: string | null;
        created_at: string;
        updated_at: string;
        is_closed: boolean;
}

export interface ExpectedDocumentTemplate {
        id: number;
        document_request_template_id: number;
        title: string;
        description: string;
        example_file_path?: string | null;
        example_mime_type?: string | null;
}

export interface ExpectedDocument {
        id: number;
        document_request_id: number;
        title: string;
        description: string;
        status: ExpectedDocumentStatus;
        rejection_reason?: string | null;
        example_file_path?: string | null;
        example_mime_type?: string | null;
}

export interface DocumentFile {
        id: number;
        document_request_id: number;
        file_name: string;
        file_path: string;
        mime_type: string;
        file_size: number;
        expected_document_id: number;
        uploaded_at: string;
        s3_version_id?: string | null;
        uploaded_by: number;
        uploaded_by_first_name: string;
        uploaded_by_last_name: string;
}

export interface RequestComment {
        id: number;
        request_id: number;
        user_id: number;
        comment: string;
        created_at: string;
        updated_at: string;
        user_first_name: string;
        user_last_name: string;
}

export interface InvitationCode {
        id: number;
        code: string;
        created_by: number;
        department_id: number;
        department_name: string;
        used_at?: string | null;
        expires_at?: string | null;
        created_at: string;
}

export type RequestStatus = "pending" | "uploaded" | "overdue";
export type UserRole = "admin" | "member";
export type ExpectedDocumentStatus = "accepted" | "rejected" | "uploaded" | "pending";
export type RecurrenceUnit = "day" | "week" | "month" | "year";

export const ALLOWED_EXTENSIONS = [".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"];
