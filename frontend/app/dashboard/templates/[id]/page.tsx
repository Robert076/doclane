import { notFound } from "next/navigation";
import TemplateDetailsHeader from "./_components/DetailsHeader/TemplateDetailsHeader";
import {
        getClientsByProfessional,
        getDocumentRequestTemplateByID,
        getExpectedDocumentTemplatesByTemplate,
} from "@/lib/api/api";
import TemplateActions from "./_components/TemplateActions/TemplateDetailsActions";
import TemplateDetailsSummary from "./_components/DetailsCard/TemplateDetailsSummary";
import TemplateDetailsExpectedDocuments from "./_components/TemplateDetailsExpectedDocuments/TemplateDetailsExpectedDocuments";

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
        console.log(templateResponse);
        console.log(clientsResponse);
        console.log(expectedDocumentsResponse);

        if (
                !templateResponse ||
                !clientsResponse ||
                !templateResponse.data ||
                !clientsResponse.data
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
