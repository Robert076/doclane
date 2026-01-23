import "./EmptyFilesSection.css";
const EmptyFileSection = () => {
  return (
    <section className="files-section details-card">
      <h2 className="section-title">Documents</h2>
      <div className="empty-files">
        <p>No documents have been uploaded yet.</p>
      </div>
    </section>
  );
};

export default EmptyFileSection;
