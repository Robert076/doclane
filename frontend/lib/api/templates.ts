import { doclaneHTTPHelper } from "@/lib/api/core";
import { APIResponse, Template, ExpectedDocumentTemplate } from "@/types";

export async function createTemplate(payload: {
        title: string;
        description?: string;
        department_id: number;
        is_recurring: boolean;
        recurrence_cron?: string;
        expected_documents: Array<{
                title: string;
                description: string;
                example_file?: File;
        }>;
}): Promise<APIResponse<number>> {
        const formData = new FormData();
        formData.append("title", payload.title);
        if (payload.description) formData.append("description", payload.description);
        formData.append("department_id", payload.department_id.toString());
        formData.append("is_recurring", payload.is_recurring.toString());
        if (payload.recurrence_cron)
                formData.append("recurrence_cron", payload.recurrence_cron);

        payload.expected_documents.forEach((doc, i) => {
                formData.append(`expected_documents[${i}][title]`, doc.title);
                formData.append(`expected_documents[${i}][description]`, doc.description);
                if (doc.example_file) {
                        formData.append(
                                `expected_documents[${i}][example_file]`,
                                doc.example_file,
                        );
                }
        });

        return doclaneHTTPHelper("/templates", {
                method: "POST",
                formData,
                revalidate: "/dashboard/templates",
        });
}

export async function getTemplates(): Promise<APIResponse<Template[]>> {
        return doclaneHTTPHelper("/templates", { method: "GET" });
}

export async function getTemplateByID(id: number): Promise<APIResponse<Template>> {
        return doclaneHTTPHelper(`/templates/${id}`, { method: "GET" });
}

export async function patchTemplate(
        id: number,
        payload: {
                title?: string;
                description?: string;
                is_recurring?: boolean;
                recurrence_cron?: string;
        },
): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${id}`, {
                method: "PATCH",
                body: payload,
                revalidate: `/dashboard/templates/${id}`,
        });
}

export async function archiveTemplate(id: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${id}/archive`, {
                method: "POST",
                revalidate: "/dashboard/templates",
        });
}

export async function unarchiveTemplate(id: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${id}/unarchive`, {
                method: "POST",
                revalidate: "/dashboard/templates",
        });
}

export async function deleteTemplate(id: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${id}`, {
                method: "DELETE",
                revalidate: "/dashboard/templates",
        });
}

export async function getExpectedDocumentTemplatesByTemplate(
        templateID: number,
): Promise<APIResponse<ExpectedDocumentTemplate[]>> {
        return doclaneHTTPHelper(`/templates/${templateID}/expected-documents`, {
                method: "GET",
        });
}

export async function deleteExpectedDocumentTemplate(
        templateID: number,
        expectedDocTemplateID: number,
): Promise<APIResponse> {
        return doclaneHTTPHelper(
                `/templates/${templateID}/expected-documents/${expectedDocTemplateID}`,
                {
                        method: "DELETE",
                        revalidate: `/dashboard/templates/${templateID}`,
                },
        );
}

export async function presignTemplateExample(
        templateID: number,
        docID: number,
): Promise<APIResponse<string>> {
        return doclaneHTTPHelper(
                `/templates/${templateID}/expected-documents/${docID}/presign-example`,
                { method: "GET" },
        );
}
