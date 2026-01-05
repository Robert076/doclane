export type UserRole = "PROFESSIONAL" | "CLIENT";

export interface User {
  id: string;
  email: string;
  role: UserRole;
  professional_id?: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ApiResponse<T> {
  data: T;
  error?: string;
}
