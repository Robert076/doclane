"use client";
import { useState, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import { MdNotifications } from "react-icons/md";
import { AuditEvent } from "@/types";
import { getNotifications, markNotificationsSeen } from "@/lib/api/notifications";
import { formatDate } from "@/lib/client/formatDate";
import "./NotificationBell.css";

const ITEMS_PER_PAGE = 5;

const EVENT_LABELS: Record<string, string> = {
  "request.created": "Dosar creat",
  "request.updated": "Dosar actualizat",
  "request.claimed": "Dosar revendicat",
  "request.unclaimed": "Dosar nerevendicat",
  "request.closed": "Dosar finalizat",
  "request.reopened": "Dosar redeschis",
  "request.cancelled": "Dosar anulat",
  "document.uploaded": "Document încărcat",
  "document.approved": "Document aprobat",
  "document.rejected": "Document respins",
  "user.notified": "Utilizator notificat",
  "user.deactivated": "Cont dezactivat",
  "department.created": "Departament creat",
};

function getLabel(eventType: string): string {
  return EVENT_LABELS[eventType] ?? eventType;
}

function getNotificationText(event: AuditEvent): string {
  const actor =
    event.actor_first_name && event.actor_last_name
      ? `${event.actor_first_name} ${event.actor_last_name}`
      : null;

  const title = event.metadata?.title ? String(event.metadata.title) : null;
  const label = getLabel(event.event_type);

  if (actor && title) return `${actor}: ${label} — "${title}"`;
  if (actor) return `${actor}: ${label}`;
  if (title) return `${label} — "${title}"`;
  return label;
}

export default function NotificationBell() {
  const router = useRouter();
  const [isOpen, setIsOpen] = useState(false);
  const [notifications, setNotifications] = useState<AuditEvent[]>([]);
  const [seenAt, setSeenAt] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const panelRef = useRef<HTMLDivElement>(null);

  const fetchNotifications = async () => {
    setIsLoading(true);
    const res = await getNotifications(50);
    const data = res.data as any;
    setNotifications(data?.notifications ?? []);
    setSeenAt(data?.seen_at ?? null);
    setIsLoading(false);
  };

  useEffect(() => {
    fetchNotifications();
    const interval = setInterval(fetchNotifications, 30000);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (panelRef.current && !panelRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    }
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [isOpen]);

  const unreadCount = seenAt
    ? notifications.filter((n) => new Date(n.occurred_at) > new Date(seenAt)).length
    : notifications.length;

  const totalPages = Math.ceil(notifications.length / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const currentItems = notifications.slice(startIndex, startIndex + ITEMS_PER_PAGE);

  const handleToggle = () => {
    if (!isOpen) {
      fetchNotifications();
      markNotificationsSeen();
    }
    setIsOpen((prev) => !prev);
  };

  const handleClick = (event: AuditEvent) => {
    if (event.resource_type === "request") {
      router.push(`/dashboard/requests/${event.resource_id}`);
      setIsOpen(false);
    }
    if (event.resource_type === "document" && event.metadata?.request_id) {
      router.push(`/dashboard/requests/${event.metadata.request_id}`);
      setIsOpen(false);
    }
  };

  const isUnread = (event: AuditEvent): boolean => {
    if (!seenAt) return true;
    return new Date(event.occurred_at) > new Date(seenAt);
  };

  return (
    <div className="notif-container" ref={panelRef}>
      <button className="notif-bell" onClick={handleToggle} aria-label="Notificări">
        <MdNotifications size={22} />
        {unreadCount > 0 && (
          <span className="notif-badge">{unreadCount > 9 ? "9+" : unreadCount}</span>
        )}
      </button>

      {isOpen && (
        <div className="notif-panel">
          <div className="notif-panel-header">
            <span className="notif-panel-title">Notificări</span>
            {unreadCount > 0 && (
              <span className="notif-panel-count">
                {unreadCount} {unreadCount === 1 ? "nouă" : "noi"}
              </span>
            )}
          </div>

          <div className="notif-panel-list">
            {isLoading ? (
              <p className="notif-empty">Se încarcă...</p>
            ) : notifications.length === 0 ? (
              <p className="notif-empty">Nicio notificare.</p>
            ) : (
              currentItems.map((event) => (
                <button
                  key={event.id}
                  className={`notif-item ${isUnread(event) ? "notif-item--unread" : ""}`}
                  onClick={() => handleClick(event)}
                >
                  <span className="notif-text">{getNotificationText(event)}</span>
                  <span className="notif-time">{formatDate(event.occurred_at)}</span>
                </button>
              ))
            )}
          </div>

          {totalPages > 1 && (
            <div className="notif-pagination">
              <button
                className="notif-pag-btn"
                disabled={currentPage === 1}
                onClick={() => setCurrentPage((p) => p - 1)}
              >
                ‹
              </button>
              <span className="notif-pag-info">
                {currentPage} / {totalPages}
              </span>
              <button
                className="notif-pag-btn"
                disabled={currentPage === totalPages}
                onClick={() => setCurrentPage((p) => p + 1)}
              >
                ›
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}