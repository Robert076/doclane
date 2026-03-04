import Link from "next/link";
import CreateTemplateForm from "../_components/CreateTemplateForm";
import "./style.css";
import { UI_TEXT } from "@/locales/ro";

export default function CreateTemplatePage() {
        return (
                <div className="create-template">
                        <header className="create-template-header">
                                <Link className="back-link" href="/dashboard/templates">
                                        {UI_TEXT.common.back}
                                </Link>
                                <h1 className="overview-h1">Şablon nou</h1>
                                <p className="overview-p">
                                        Completează detaliile şablonului de dosar.
                                </p>
                        </header>
                        <CreateTemplateForm />
                </div>
        );
}
