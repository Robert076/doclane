import { getTemplates } from "@/lib/api/api";
import TemplatesSection from "../../../components/Pages/TemplatesComponents/TemplatesSection";
import PageHeader from "@/components/PageHeader/PageHeader";
import { notFound } from "next/navigation";
import TemplatesActions from "../../../components/Pages/TemplatesComponents/TemplatesActions";

export default async function TemplatesPage() {
        const templatesResponse = await getTemplates();
        if (!templatesResponse.success || !templatesResponse.data) {
                notFound();
        }

        const templates = templatesResponse.data;

        return (
                <div>
                        <PageHeader
                                title="Şabloanele tale"
                                subtitle="Administrează şi gestionează şabloanele tale."
                        />
                        <TemplatesActions />
                        <TemplatesSection templates={templates} />
                </div>
        );
}
