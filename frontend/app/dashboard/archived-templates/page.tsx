import { getTemplates } from "@/lib/api/api";
import { notFound } from "next/navigation";
import { DocumentRequestTemplate } from "@/types";
import ArchivedTemplatesSection from "@/components/Pages/ArchivedTemplatesComponents/ArchivedTemplatesSection";
import PageHeader from "@/components/PageHeader/PageHeader";

const ArchivedTemplates = async () => {
        const templateResponse = await getTemplates();

        if (!templateResponse?.data) {
                notFound();
        }

        const templates = templateResponse.data as DocumentRequestTemplate[];

        return (
                <div>
                        <PageHeader
                                title="Şabloane arhivate"
                                subtitle="Restaurează şi gestionează şabloanele arhivate."
                        />
                        <ArchivedTemplatesSection templates={templates} />
                </div>
        );
};

export default ArchivedTemplates;
