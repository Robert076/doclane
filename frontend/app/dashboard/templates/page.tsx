import { getTemplates } from "@/lib/api/api";
import "./style.css";
import Link from "next/link";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { MdAdd } from "react-icons/md";
import TemplatesSection from "./_components/TemplatesSection";

export default async function TemplatesPage() {
        const res = await getTemplates();
        const templates = res.success ? res.data : [];

        return (
                <div className="templates-container">
                        <header className="templates-header">
                                <div>
                                        <h1 className="overview-h1">Şabloane de dosar</h1>
                                        <p className="overview-p">
                                                Creează şi gestionează şabloanele tale de
                                                dosar.
                                        </p>
                                </div>
                                <div className="header-actions">
                                        <Link href="/dashboard/templates/create">
                                                <ButtonPrimary
                                                        text="Şablon nou"
                                                        icon={MdAdd}
                                                        type="button"
                                                />
                                        </Link>
                                </div>
                        </header>
                        <TemplatesSection templates={templates} />
                </div>
        );
}
