import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import UploadDocumentButton from "@/components/Buttons/UploadDocumentButton/UploadDocumentButton";
import "./DetailCardsActionSidebar.css";
import SectionTitle from "@/components/SectionTitle/SectionTitle";

export default function DetailCardsActionSidebar({ id }: { id: string }) {
  return (
    <aside className="details-card actions-sidebar">
      <SectionTitle text="Actions" />
      <div className="action-buttons">
        <UploadDocumentButton requestId={id} />
        <ButtonPrimary text="Mark as Completed" variant="secondary" fullWidth />
      </div>
    </aside>
  );
}
