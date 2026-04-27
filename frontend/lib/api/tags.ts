import { APIResponse, Tag } from "@/types";
import { doclaneHTTPHelper } from "./core";

export async function getTags(): Promise<APIResponse<Tag[]>> {
        return doclaneHTTPHelper("/tags", { method: "GET" });
}

export async function createTag(dto: {
        name: string;
        color: string;
}): Promise<APIResponse<Tag>> {
        return doclaneHTTPHelper("/tags", { method: "POST", body: dto });
}

export async function updateTag(
        id: number,
        dto: { name: string; color: string },
): Promise<APIResponse<Tag>> {
        return doclaneHTTPHelper(`/tags/${id}`, { method: "PATCH", body: dto });
}

export async function deleteTag(id: number): Promise<APIResponse<void>> {
        return doclaneHTTPHelper(`/tags/${id}`, { method: "DELETE" });
}

export async function setTemplateTags(
        templateID: number,
        tagIDs: number[],
): Promise<APIResponse<void>> {
        return doclaneHTTPHelper(`/templates/${templateID}/tags`, {
                method: "PUT",
                body: { tag_ids: tagIDs },
        });
}
