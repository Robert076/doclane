import { AuditEvent } from "@/types";
import { formatDate } from "@/lib/client/formatDate";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import "./RequestTimeline.css";

interface RequestTimelineProps {
  events: AuditEvent[];
}

interface EventConfig {
  label: string;
  badgeClass: string;
  dotClass: string;
  renderMeta?: (metadata: Record<string, unknown> | null) => React.ReactNode;
}

const EVENT_CONFIG: Record<string, EventConfig> = {
  "request.created": {
    label: "Dosar creat",
    badgeClass: "tl-badge--created",
    dotClass: "tl-dot--created",
    renderMeta: (m) =>
      m?.title ? <span>Titlu: <strong>{String(m.title)}</strong></span> : null,
  },
  "request.updated": {
    label: "Dosar actualizat",
    badgeClass: "tl-badge--updated",
    dotClass: "tl-dot--updated",
    renderMeta: (m) =>
      m?.title ? <span>Titlu nou: <strong>{String(m.title)}</strong></span> : null,
  },
  "request.claimed": {
    label: "Dosar revendicat",
    badgeClass: "tl-badge--claimed",
    dotClass: "tl-dot--claimed",
    renderMeta: () => null,
  },
  "request.unclaimed": {
    label: "Dosar nerevendicat",
    badgeClass: "tl-badge--unclaimed",
    dotClass: "tl-dot--unclaimed",
    renderMeta: () => null,
  },
  "request.closed": {
    label: "Dosar finalizat",
    badgeClass: "tl-badge--closed",
    dotClass: "tl-dot--closed",
    renderMeta: () => null,
  },
  "request.reopened": {
    label: "Dosar redeschis",
    badgeClass: "tl-badge--reopened",
    dotClass: "tl-dot--reopened",
    renderMeta: () => null,
  },
  "request.cancelled": {
    label: "Dosar anulat",
    badgeClass: "tl-badge--cancelled",
    dotClass: "tl-dot--cancelled",
    renderMeta: () => null,
  },
  "document.uploaded": {
    label: "Document încărcat",
    badgeClass: "tl-badge--uploaded",
    dotClass: "tl-dot--uploaded",
    renderMeta: (m) =>
      m?.file_name ? <span>{String(m.file_name)}</span> : null,
  },
  "document.approved": {
    label: "Document aprobat",
    badgeClass: "tl-badge--approved",
    dotClass: "tl-dot--approved",
    renderMeta: (m) =>
      m?.file_name ? <span>{String(m.file_name)}</span> : null,
  },
  "document.rejected": {
    label: "Document respins",
    badgeClass: "tl-badge--rejected",
    dotClass: "tl-dot--rejected",
    renderMeta: (m) =>
      m?.rejection_reason ? (
        <span>Motiv: {String(m.rejection_reason)}</span>
      ) : null,
  },
  "user.notified": {
    label: "Utilizator notificat",
    badgeClass: "tl-badge--notified",
    dotClass: "tl-dot--notified",
    renderMeta: () => null,
  },
  "user.deactivated": {
    label: "Cont dezactivat",
    badgeClass: "tl-badge--deactivated",
    dotClass: "tl-dot--deactivated",
    renderMeta: () => null,
  },
};

function getConfig(eventType: string): EventConfig {
  return (
    EVENT_CONFIG[eventType] ?? {
      label: eventType,
      badgeClass: "tl-badge--updated",
      dotClass: "tl-dot--updated",
      renderMeta: () => null,
    }
  );
}

export default function RequestTimeline({ events }: RequestTimelineProps) {
  const sorted = [...events ?? []].sort(
    (a, b) =>
      new Date(a.occurred_at).getTime() - new Date(b.occurred_at).getTime(),
  );

  return (
    <div className="timeline-section">
      <SectionTitle text="Istoric dosar" />
      {sorted.length === 0 ? (
        <NotFound
          text="Nu există activitate înregistrată"
          subtext="Acțiunile asupra dosarului vor apărea aici."
          background="#fff"
        />
      ) : (
        <ol className="tl-list">
          {sorted.map((event, index) => {
            const config = getConfig(event.event_type);
            const meta = config.renderMeta?.(event.metadata);
            const isLast = index === sorted.length - 1;

            return (
              <li key={event.id} className="tl-item">
                <div className="tl-indicator">
                  <div className={`tl-dot ${config.dotClass}`} />
                  {!isLast && <div className="tl-line" />}
                </div>
                <div className="tl-body">
                  <div className="tl-row">
                    <span className="tl-label">{config.label}</span>
                    <span className="tl-time">
                      {formatDate(event.occurred_at)}
                    </span>
                  </div>
                  {meta && <div className="tl-meta">{meta}</div>}
                  {(event.actor_first_name || event.actor_last_name) && (
                        <div className="tl-actor">
                        {event.actor_first_name} {event.actor_last_name}
                        </div>
                        )}
                  <span className={`tl-badge ${config.badgeClass}`}>
                    {config.label}
                  </span>
                </div>
              </li>
            );
          })}
        </ol>
      )}
    </div>
  );
}