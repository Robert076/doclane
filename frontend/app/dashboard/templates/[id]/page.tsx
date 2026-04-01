import { notFound } from "next/navigation";
import TemplateDetailsHeader from "@/components/Pages/TemplatesComponents/TemplateDetailsHeader";
import TemplateDetailsActions from "@/components/Pages/TemplatesComponents/TemplateDetailsActions";
import TemplateDetailsSummary from "@/components/Pages/TemplatesComponents/TemplateDetailsSummary";
import TemplateDetailsExpectedDocuments from "@/components/Pages/TemplatesComponents/TemplateDetailsExpectedDocuments";
import { getTemplateByID, getExpectedDocumentTemplatesByTemplate } from "@/lib/api/templates";

interface PageProps {
        params: Promise<{ id: string }>;
}

export default async function TemplateDetailsPage({ params }: PageProps) {
        const { id } = await params;

        const [templateResponse, expectedDocumentsResponse] = await Promise.all([
                getTemplateByID(+id),
                getExpectedDocumentTemplatesByTemplate(+id),
        ]);

        if (!templateResponse?.data || !expectedDocumentsResponse?.data) {
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
                                                templateID={data.id}
                                                documents={expectedDocumentsResponse.data}
                                        />
                                </div>
                                <TemplateDetailsActions id={+id} template={data} />
                        </div>
                </div>
        );
}
