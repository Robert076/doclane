import Link from "next/link";
import CreateTemplateForm from "@/components/Pages/TemplatesComponents/CreateTemplateForm";
import { getCurrentUser } from "@/lib/api/users";
import { getDepartments } from "@/lib/api/departments";
import { UI_TEXT } from "@/locales/ro";
import { redirect } from "next/navigation";
import "./style.css";

export default async function CreateTemplatePage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");
        if (userResponse.data.role !== "admin") redirect("/dashboard/templates");

        const departmentsResponse = await getDepartments();
        const departments = departmentsResponse.data ?? [];

        return (
                <div className="add-template">
                        <header className="create-template-header">
                                <Link className="back-link" href="/dashboard/templates">
                                        {UI_TEXT.common.back}
                                </Link>
                                <h1>Şablon nou</h1>
                                <p>Completează detaliile şablonului de dosar.</p>
                        </header>
                        <CreateTemplateForm departments={departments} />
                </div>
        );
}
