import { notFound } from "next/navigation";
import RequestDetailsActions from "@/components/Pages/RequestsComponents/RequestDetailsActions";
import RequestDetailsHeader from "@/components/Pages/RequestsComponents/RequestDetailsHeader";
import RequestTabs from "@/components/Pages/RequestsComponents/RequestTabs";
import { getRequestById, getFilesByRequestId, getCommentsByRequest } from "@/lib/api/requests";

interface PageProps {
        params: Promise<{ id: string }>;
}

export default async function RequestDetailsPage({ params }: PageProps) {
        const { id } = await params;
        const requestId = parseInt(id);

        const [requestResponse, filesResponse, commentsResponse] = await Promise.all([
                getRequestById(requestId),
                getFilesByRequestId(requestId),
                getCommentsByRequest(requestId),
        ]);

        if (!requestResponse.success || !requestResponse.data) notFound();

        const request = requestResponse.data;
        const files = filesResponse?.data ?? [];
        const comments = commentsResponse?.data ?? [];

        return (
                <div className="details-container">
                        <RequestDetailsHeader data={request} />
                        <div className="details-grid">
                                <div className="main-content">
                                        <RequestTabs
                                                data={request}
                                                files={files}
                                                comments={comments}
                                                requestId={requestId}
                                        />
                                </div>
                                <RequestDetailsActions assignee={request.assignee} />
                        </div>
                </div>
        );
}
