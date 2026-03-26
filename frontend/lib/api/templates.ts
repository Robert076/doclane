import { doclaneHTTPHelper } from "@/lib/api/core";
import { APIResponse, Template, ExpectedDocumentTemplate, PresignedURL } from "@/types";

export async function createTemplate(payload: object): Promise<APIResponse<number>> {
        return doclaneHTTPHelper("/templates", {
                method: "POST",
                body: payload,
                revalidate: "/dashboard/templates",
        });
}

export async function getTemplates(): Promise<APIResponse<Template[]>> {
        return doclaneHTTPHelper("/templates", {
                method: "GET",
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
                revalidate: "dashboard/archived-templates",
        });
}

export async function instantiateTemplate(
        templateID: number,
        payload: {
                client_id: number;
                is_scheduled: boolean;
                scheduled_for?: string;
                due_date?: string;
        },
): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${templateID}/instantiate`, {
                method: "POST",
                body: payload,
                revalidate: "/dashboard/requests",
        });
}

export async function addExpectedDocumentTemplate(
        templateID: number,
        title: string,
        description: string,
        exampleFile?: File,
): Promise<APIResponse> {
        const formData = new FormData();
        formData.append("title", title);
        formData.append("description", description);
        if (exampleFile) {
                formData.append("example_file", exampleFile);
        }

        return doclaneHTTPHelper(`/templates/${templateID}/expected-documents`, {
                method: "POST",
                formData,
                revalidate: `/dashboard/templates/${templateID}`,
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
): Promise<APIResponse<PresignedURL>> {
        return doclaneHTTPHelper(
                `/templates/${templateID}/expected-documents/${docID}/presign-example`,
                { method: "GET" },
        );
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
