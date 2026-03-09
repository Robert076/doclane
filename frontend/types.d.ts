export interface APIResponse<T = unknown> {
        success: boolean;
        message: string;
        error?: string;
        data?: T;
}

export interface User {
        id: string;
        email: string;
        role: UserRole;
        professional_id?: string | null;
        is_active: boolean;
        created_at: string;
        updated_at: string;
        first_name: string;
        last_name: string;
}

export interface DocumentRequest {
        id: number;
        professional_id: number;
        client_id: number;
        client_email: string;
        client_first_name: string;
        client_last_name: string;
        title: string;
        description?: string | null;
        due_date?: string | null;
        is_recurring: boolean;
        recurrence_cron?: string | null;
        is_scheduled: boolean;
        scheduled_for?: string | null;
        next_due_at?: string | null;
        last_uploaded_at?: string | null;
        is_closed: boolean;
        template_id?: number | null;
        status: RequestStatus;
        created_at: string;
        updated_at: string;
        expected_documents: ExpectedDocument[];
}

export interface DocumentRequestTemplate {
        id: number;
        title: string;
        description?: string | null;
        is_recurring: boolean;
        recurrence_cron?: string | null;
        created_by: number;
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

export interface DocumentFile {
        id: number;
        document_request_id: number;
        file_name: string;
        file_path: string;
        mime_type: string;
        file_size: number;
        expected_document_id: number;
        uploaded_at: string;
        s3_version_id?: string;
        uploaded_by: number;
        uploaded_by_first_name: string;
        uploaded_by_last_name: string;
}

export interface ExpectedDocument {
        id: number;
        document_request_id: number;
        title: string;
        description: string;
        status: ExpectedDocumentStatus;
        rejection_reason: string;
        example_file_path?: string | null;
        example_mime_type?: string | null;
}

export type RequestStatus = "pending" | "uploaded" | "overdue";
export type UserRole = "CLIENT" | "PROFESSIONAL";

export interface ApiResponse<T> {
        data: T;
        error?: string;
}

export const ALLOWED_EXTENSIONS = [".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"];

export type RecurrenceUnit = "day" | "week" | "month" | "year";

export type ExpectedDocumentStatus = "approved" | "rejected" | "uploaded" | "pending";
