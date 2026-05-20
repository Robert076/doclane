import jsPDF from "jspdf";
import { AuditEvent } from "@/types";
import { getEventLabel } from "@/components/Pages/RequestsComponents/RequestTimeline";

function formatDateForPDF(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleString("ro-RO", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function getActorName(event: AuditEvent): string {
  if (event.actor_first_name && event.actor_last_name) {
    return `${event.actor_first_name} ${event.actor_last_name}`;
  }
  return "Sistem";
}

function getMetaText(event: AuditEvent): string | null {
  const m = event.metadata;
  if (!m) return null;

  switch (event.event_type) {
    case "request.created":
    case "request.updated":
      return m.title ? `Titlu: ${String(m.title)}` : null;
    case "document.uploaded":
      return m.file_name ? `Fișier: ${String(m.file_name)}` : null;
    case "document.approved":
      return m.title ? `Document: ${String(m.title)}` : null;
    case "document.rejected":
      return m.rejection_reason
        ? `Motiv respingere: ${String(m.rejection_reason)}`
        : null;
    default:
      return null;
  }
}

export function exportAuditLogToPDF(
  events: AuditEvent[],
  requestTitle?: string,
) {
  const doc = new jsPDF();
  const pageWidth = doc.internal.pageSize.getWidth();
  const margin = 20;
  const contentWidth = pageWidth - margin * 2;
  let y = margin;

  // Header
  doc.setFontSize(18);
  doc.setFont("helvetica", "bold");
  doc.text("Istoric dosar", margin, y);
  y += 8;

  if (requestTitle) {
    doc.setFontSize(12);
    doc.setFont("helvetica", "normal");
    doc.text(requestTitle, margin, y);
    y += 6;
  }

  doc.setFontSize(9);
  doc.setTextColor(130);
  doc.text(
    `Generat la: ${new Date().toLocaleString("ro-RO")}`,
    margin,
    y,
  );
  doc.setTextColor(0);
  y += 10;

  // Separator line
  doc.setDrawColor(200);
  doc.line(margin, y, pageWidth - margin, y);
  y += 8;

  // Table header
  doc.setFontSize(9);
  doc.setFont("helvetica", "bold");
  doc.setFillColor(245, 245, 245);
  doc.rect(margin, y - 4, contentWidth, 7, "F");
  doc.text("Data", margin + 2, y);
  doc.text("Eveniment", margin + 40, y);
  doc.text("Actor", margin + 105, y);
  doc.text("Detalii", margin + 145, y);
  y += 6;

  doc.setFont("helvetica", "normal");
  doc.setFontSize(8.5);

  // Sort chronologically for PDF (oldest first)
  const sorted = [...events].sort(
    (a, b) =>
      new Date(a.occurred_at).getTime() - new Date(b.occurred_at).getTime(),
  );

  sorted.forEach((event, index) => {
    // Check if we need a new page
    if (y > 270) {
      doc.addPage();
      y = margin;

      // Repeat header on new page
      doc.setFontSize(9);
      doc.setFont("helvetica", "bold");
      doc.setFillColor(245, 245, 245);
      doc.rect(margin, y - 4, contentWidth, 7, "F");
      doc.text("Data", margin + 2, y);
      doc.text("Eveniment", margin + 40, y);
      doc.text("Actor", margin + 105, y);
      doc.text("Detalii", margin + 145, y);
      y += 6;
      doc.setFont("helvetica", "normal");
      doc.setFontSize(8.5);
    }

    // Alternating row background
    if (index % 2 === 0) {
      doc.setFillColor(250, 250, 252);
      doc.rect(margin, y - 4, contentWidth, 7, "F");
    }

    const dateStr = formatDateForPDF(event.occurred_at);
    const label = getEventLabel(event.event_type);
    const actor = getActorName(event);
    const meta = getMetaText(event);

    doc.text(dateStr, margin + 2, y);
    doc.text(label, margin + 40, y);
    doc.text(actor, margin + 105, y);

    if (meta) {
      // Truncate long meta text
      const maxLen = 30;
      const truncated =
        meta.length > maxLen ? meta.substring(0, maxLen) + "..." : meta;
      doc.text(truncated, margin + 145, y);
    }

    y += 7;
  });

  // Footer
  y += 4;
  doc.setDrawColor(200);
  doc.line(margin, y, pageWidth - margin, y);
  y += 6;
  doc.setFontSize(8);
  doc.setTextColor(130);
  doc.text(`Total evenimente: ${sorted.length}`, margin, y);
  doc.text("Doclane", pageWidth - margin - 15, y);

  // Download
  const fileName = requestTitle
    ? `istoric-${requestTitle.toLowerCase().replace(/\s+/g, "-")}.pdf`
    : "istoric-dosar.pdf";
  doc.save(fileName);
}