import jsPDF from "jspdf";
import { Stats } from "@/types/stats";

function formatHours(hours: number): string {
        if (hours === 0) return "N/A";
        const days = Math.floor(hours / 24);
        const remainingHours = Math.round(hours % 24);
        if (days === 0) return `${remainingHours}h`;
        if (remainingHours === 0) return `${days}z`;
        return `${days}z ${remainingHours}h`;
}

function formatPercent(value: number): string {
        return `${value.toFixed(1)}%`;
}

function formatChange(value: number): string {
        const sign = value >= 0 ? "+" : "";
        return `${sign}${value.toFixed(1)}%`;
}

export function generateStatsPDF(stats: Stats) {
        const doc = new jsPDF();
        const pageWidth = doc.internal.pageSize.getWidth();
        let y = 20;

        const PRIMARY = [255, 87, 34] as [number, number, number];
        const DARK = [26, 32, 44] as [number, number, number];
        const GRAY = [71, 85, 105] as [number, number, number];
        const LIGHT_GRAY = [241, 245, 249] as [number, number, number];

        const addTitle = (text: string) => {
                doc.setFontSize(20);
                doc.setTextColor(...DARK);
                doc.setFont("helvetica", "bold");
                doc.text(text, pageWidth / 2, y, { align: "center" });
                y += 8;
        };

        const addSubtitle = (text: string) => {
                doc.setFontSize(10);
                doc.setTextColor(...GRAY);
                doc.setFont("helvetica", "normal");
                doc.text(text, pageWidth / 2, y, { align: "center" });
                y += 12;
        };

        const addSectionHeader = (text: string) => {
                if (y > 250) {
                        doc.addPage();
                        y = 20;
                }
                doc.setFillColor(...PRIMARY);
                doc.rect(14, y - 4, pageWidth - 28, 8, "F");
                doc.setFontSize(10);
                doc.setTextColor(255, 255, 255);
                doc.setFont("helvetica", "bold");
                doc.text(text.toUpperCase(), 18, y + 1);
                y += 10;
        };

        const addRow = (label: string, value: string, highlight = false) => {
                if (y > 270) {
                        doc.addPage();
                        y = 20;
                }
                if (highlight) {
                        doc.setFillColor(...LIGHT_GRAY);
                        doc.rect(14, y - 4, pageWidth - 28, 7, "F");
                }
                doc.setFontSize(10);
                doc.setTextColor(...GRAY);
                doc.setFont("helvetica", "normal");
                doc.text(label, 18, y);
                doc.setTextColor(...DARK);
                doc.setFont("helvetica", "bold");
                doc.text(value, pageWidth - 18, y, { align: "right" });
                y += 8;
        };

        const addDivider = () => {
                doc.setDrawColor(226, 232, 240);
                doc.line(14, y, pageWidth - 14, y);
                y += 6;
        };

        // header
        addTitle("Raport Statistici Doclane");
        addSubtitle(
                `Generat pe ${new Date().toLocaleDateString("ro-RO", { day: "2-digit", month: "long", year: "numeric" })}`,
        );
        addDivider();

        // requests
        addSectionHeader("Cereri");
        addRow("Cereri deschise", stats.total_open_requests.toString(), false);
        addRow("Cereri finalizate", stats.total_archived_requests.toString(), true);
        addRow("Cereri retrase", stats.total_cancelled_requests.toString(), false);
        addRow("Rată finalizare", formatPercent(stats.completion_rate), true);
        addRow("Rată retragere", formatPercent(stats.cancellation_rate), false);
        addRow("Timp mediu finalizare", formatHours(stats.avg_completion_hours), true);
        y += 4;

        // activity
        addSectionHeader("Activitate");
        addRow(
                "Cereri săptămâna aceasta",
                `${stats.requests_this_week} (${formatChange(stats.weekly_change_percent)})`,
                false,
        );
        addRow("Cereri săptămâna trecută", stats.requests_last_week.toString(), true);
        addRow(
                "Cereri luna aceasta",
                `${stats.requests_this_month} (${formatChange(stats.monthly_change_percent)})`,
                false,
        );
        addRow("Cereri luna trecută", stats.requests_last_month.toString(), true);
        y += 4;

        // departments
        addSectionHeader("Departamente");
        addRow("Total departamente", stats.total_departments.toString(), false);
        addRow("Membri departamente", stats.total_department_members.toString(), true);
        y += 4;

        addSectionHeader("Cereri per departament");
        stats.requests_per_department.forEach((d, i) => {
                addRow(d.department_name, `${d.request_count} cereri`, i % 2 === 0);
        });
        y += 4;

        // users
        addSectionHeader("Utilizatori");
        addRow("Total utilizatori", stats.total_users.toString(), false);
        addRow("Cetățeni", stats.total_citizens.toString(), true);
        addRow("Membri departamente", stats.total_department_members.toString(), false);
        addRow("Utilizatori activi", stats.total_active_users.toString(), true);
        addRow("Utilizatori dezactivați", stats.total_deactivated_users.toString(), false);
        y += 4;

        // templates
        addSectionHeader("Șabloane");
        addRow("Șabloane active", stats.total_active_templates.toString(), false);
        addRow("Șabloane arhivate", stats.total_archived_templates.toString(), true);
        y += 4;

        addSectionHeader("Top 5 șabloane folosite");
        stats.most_used_templates.forEach((t, i) => {
                addRow(t.template_title, `${t.request_count} cereri`, i % 2 === 0);
        });

        // footer on each page
        const totalPages = doc.getNumberOfPages();
        for (let i = 1; i <= totalPages; i++) {
                doc.setPage(i);
                doc.setFontSize(8);
                doc.setTextColor(...GRAY);
                doc.setFont("helvetica", "normal");
                doc.text(
                        `Doclane — Pagina ${i} din ${totalPages}`,
                        pageWidth / 2,
                        doc.internal.pageSize.getHeight() - 10,
                        { align: "center" },
                );
        }

        doc.save(`statistici-doclane-${new Date().toISOString().split("T")[0]}.pdf`);
}
