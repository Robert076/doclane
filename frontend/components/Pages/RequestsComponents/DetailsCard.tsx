import { Request, RequestStatus } from "@/types";
import "./DetailsCard.css";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import StatusBadge from "@/components/Pages/RequestsComponents/StatusBadge";
import { formatDate } from "@/lib/client/formatDate";
import CardRow from "@/components/CardRow/CardRow";

export default function DetailsCard({ data }: { data: Request }) {
        const rows = [
                {
                        label: "Solicitant:",
                        value: `${data.assignee_first_name} ${data.assignee_last_name} (${data.assignee_email})`,
                        isDescription: false,
                },
                {
                        label: "Creat la:",
                        value: formatDate(data.created_at),
                        isDescription: false,
                },
                {
                        label: "Status:",
                        value: <StatusBadge status={data.status as RequestStatus} />,
                },
        ];

        if (data.due_date) {
                rows.push({
                        label: "Termen limită:",
                        value: formatDate(data.due_date),
                        isDescription: false,
                });
        } else if (data.next_due_at) {
                rows.push({
                        label: "Următorul termen:",
                        value: formatDate(data.next_due_at),
                        isDescription: false,
                });
        }

        if (data.description) {
                rows.push({
                        label: "Descriere:",
                        value: data.description,
                        isDescription: true,
                });
        }

        return (
                <section className="details-card">
                        <SectionTitle text="Detalii dosar" />
                        {rows.map((row, index) => (
                                <CardRow
                                        key={index}
                                        label={row.label}
                                        value={row.value}
                                        isDescription={row.isDescription}
                                />
                        ))}
                </section>
        );
}
