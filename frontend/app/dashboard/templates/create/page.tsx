import Link from "next/link";
import CreateTemplateForm from "../../../../components/Pages/TemplatesComponents/CreateTemplateForm";
import { UI_TEXT } from "@/locales/ro";
import "./style.css";

export default function CreateTemplatePage() {
        return (
                <div className="add-template">
                        <header className="create-template-header">
                                <Link className="back-link" href="/dashboard/templates">
                                        {UI_TEXT.common.back}
                                </Link>
                                <h1>Şablon nou</h1>
                                <p>Completează detaliile şablonului de dosar.</p>
                        </header>
                        <CreateTemplateForm />
                </div>
        );
}
