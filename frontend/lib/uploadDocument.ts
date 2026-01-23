export const ALLOWED_EXTENSIONS = [".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"];
export const MAX_FILE_SIZE = 20 * 1024 * 1024;

interface UploadResponse {
  success: boolean;
  msg: string;
  data?: any;
}

export async function uploadDocument(requestId: string, file: File): Promise<UploadResponse> {
  const fileExt = file.name.substring(file.name.lastIndexOf(".")).toLowerCase();
  if (!ALLOWED_EXTENSIONS.includes(fileExt)) {
    throw new Error(`Extensia ${fileExt} nu este permisă.`);
  }

  if (file.size > MAX_FILE_SIZE) {
    throw new Error("Fișierul depășește limita de 20MB.");
  }

  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch(
    `http://localhost:8080/api/document-requests/${requestId}/files`,
    {
      method: "POST",
      body: formData,
      credentials: "include",
    }
  );

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.msg || "Upload failed");
  }

  return result;
}
