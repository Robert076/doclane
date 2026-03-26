import { Request, RequestStatus } from "@/types";
import "./DetailsCard.css";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import StatusBadge from "@/components/Pages/RequestsComponents/StatusBadge";
import { formatDate } from "@/lib/client/formatDate";
import { UI_TEXT } from "@/locales/ro";
import CardRow from "@/components/CardRow/CardRow";

export default function DetailsCard({ data }: { data: Request }) {
        const rows = [
                {
                        label: UI_TEXT.roles.clientSingular,
                        value: data.client_email,
                        isDescription: false,
                },
                {
                        label: `${UI_TEXT.common.createdAt}:`,
                        value: formatDate(data.created_at),
                        isDescription: false,
                },
                {
                        label: UI_TEXT.common.status,
                        value: <StatusBadge status={data.status as RequestStatus} />,
                },
        ];

        if (data.due_date) {
                rows.push({
                        label: UI_TEXT.common.dueDate,
                        value: formatDate(data.due_date),
                        isDescription: false,
                });
        } else if (data.next_due_at) {
                rows.push({
                        label: UI_TEXT.common.nextDueAt,
                        value: formatDate(data.next_due_at),
                        isDescription: false,
                });
        }

        if (data.description) {
                rows.push({
                        label: UI_TEXT.common.description,
                        value: data.description,
                        isDescription: true,
                });
        }

        return (
                <section className="details-card">
                        <SectionTitle text={UI_TEXT.request.details.title} />

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
