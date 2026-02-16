import "./PaginationFooter.css";

interface PaginationFooterProps {
  currentPage: number;
  totalPages: number;
  setCurrentPage: React.Dispatch<React.SetStateAction<number>>;
}

const PaginationFooter: React.FC<PaginationFooterProps> = ({
  currentPage,
  totalPages,
  setCurrentPage,
}) => {
  return (
    <div className="pagination-footer">
      <button
        className="pag-button"
        onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
        disabled={currentPage === 1}
      >
        Previous
      </button>

      <div className="page-dots">
        {Array.from({ length: totalPages }).map((_, i) => (
          <div key={i} className={`dot ${currentPage === i + 1 ? "active" : ""}`} />
        ))}
      </div>

      <button
        className="pag-button"
        onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
        disabled={currentPage === totalPages}
      >
        Next
      </button>
    </div>
  );
};

export default PaginationFooter;
