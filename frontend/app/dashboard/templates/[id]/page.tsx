import { notFound } from "next/navigation";
import TemplateDetailsHeader from "../../../../components/Pages/TemplatesComponents/TemplateDetailsHeader";
import {
        getClientsByProfessional,
        getDocumentRequestTemplateByID,
        getExpectedDocumentTemplatesByTemplate,
} from "@/lib/api/api";
import TemplateActions from "../../../../components/Pages/TemplatesComponents/TemplateDetailsActions";
import TemplateDetailsSummary from "../../../../components/Pages/TemplatesComponents/TemplateDetailsSummary";
import TemplateDetailsExpectedDocuments from "../../../../components/Pages/TemplatesComponents/TemplateDetailsExpectedDocuments";

interface PageProps {
        params: Promise<{ id: string }>;
}

export default async function TemplateDetailsPage({ params }: PageProps) {
        const { id } = await params;

        const [templateResponse, clientsResponse, expectedDocumentsResponse] =
                await Promise.all([
                        getDocumentRequestTemplateByID(+id),
                        getClientsByProfessional(),
                        getExpectedDocumentTemplatesByTemplate(+id),
                ]);

        if (
                !templateResponse ||
                !clientsResponse ||
                !templateResponse.data ||
                !clientsResponse.data ||
                !expectedDocumentsResponse ||
                !expectedDocumentsResponse.data
        ) {
                notFound();
        }

        const data = templateResponse.data;

        return (
                <div className="details-container">
                        <TemplateDetailsHeader data={data} />
                        <div className="details-grid">
                                <div className="main-content">
                                        <TemplateDetailsSummary data={data} />
                                        <TemplateDetailsExpectedDocuments
                                                templateID={templateResponse.data.id}
                                                documents={expectedDocumentsResponse.data}
                                        />
                                </div>
                                <TemplateActions id={id} clients={clientsResponse.data} />
                        </div>
                </div>
        );
}
