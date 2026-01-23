export interface User {
  id: string;
  email: string;
  role: UserRole;
  professional_id?: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface DocumentRequest {
  id: number;
  professional_id: number;
  client_id: number;
  client_email: string;
  title: string;
  description?: string | null;
  due_date?: string | null;
  status: RequestStatus;
  created_at: string;
  updated_at: string;
}

export interface DocumentFile {
  id: number;
  document_request_id: number;
  file_name: string;
  file_path: string;
  mime_type: string;
  file_size: number;
  uploaded_at: string;
  s3_version_id?: string;
}

export type RequestStatus = "pending" | "uploaded" | "overdue";
export type UserRole = "CLIENT" | "PROFESSIONAL";

export interface ApiResponse<T> {
  data: T;
  error?: string;
}
