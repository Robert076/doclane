"use client";
import { useState, useEffect } from "react";
import { DocumentRequest, User } from "@/types";
import Request from "../Request/Request";
import NotFound from "@/components/NotFound/NotFound";
import "./RequestsSection.css";
import SearchBar from "../SearchBar/SearchBar";

interface RequestsSectionProps {
  requests: DocumentRequest[];
  user: User;
}

const ITEMS_PER_PAGE = 8;

const RequestsSection: React.FC<RequestsSectionProps> = ({ requests, user }) => {
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [searchInput, setSearchInput] = useState<string>("");

  const filteredRequests = requests.filter((req) => {
    if (!searchInput) return true;
    const searchLower = searchInput.toLowerCase().trim();

    const fullName =
      `${req.client_first_name || ""} ${req.client_last_name || ""}`.toLowerCase();

    return (
      req.title?.toLowerCase().includes(searchLower) ||
      req.description?.toLowerCase().includes(searchLower) ||
      req.client_first_name?.toLowerCase().includes(searchLower) ||
      req.client_last_name?.toLowerCase().includes(searchLower) ||
      req.client_email?.toLowerCase().includes(searchLower) ||
      fullName.includes(searchLower) ||
      req.status?.toLowerCase().includes(searchLower)
    );
  });

  useEffect(() => {
    setCurrentPage(1);
  }, [searchInput]);

  if (requests.length === 0) {
    return user.role === "PROFESSIONAL" ? (
      <NotFound
        text="No document requests found."
        subtext="Create a request to get started."
        background="#fff"
      />
    ) : (
      <NotFound
        text="No document requests found."
        subtext="Seems like it's your lucky day."
        background="#fff"
      />
    );
  }

  const totalPages = Math.ceil(filteredRequests.length / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const endIndex = startIndex + ITEMS_PER_PAGE;
  const currentRequests = filteredRequests.slice(startIndex, endIndex);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const getPageNumbers = () => {
    const pages = [];
    const showEllipsis = totalPages > 7;

    if (!showEllipsis) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      if (currentPage <= 3) {
        pages.push(1, 2, 3, 4, "ellipsis", totalPages);
      } else if (currentPage >= totalPages - 2) {
        pages.push(1, "ellipsis", totalPages - 3, totalPages - 2, totalPages - 1, totalPages);
      } else {
        pages.push(
          1,
          "ellipsis",
          currentPage - 1,
          currentPage,
          currentPage + 1,
          "ellipsis",
          totalPages,
        );
      }
    }
    return pages;
  };

  return (
    <div className="requests-section">
      <SearchBar
        value={searchInput}
        onChange={setSearchInput}
        placeholder="Search requests..."
      />

      {filteredRequests.length === 0 ? (
        <NotFound text="No requests found." subtext="Try adjusting your search terms." />
      ) : (
        <>
          <div className="requests-grid">
            {currentRequests.map((req) => (
              <Request key={req.id} request={req} user={user} searchTerm={searchInput} />
            ))}
          </div>

          {totalPages > 1 && (
            <div className="pagination">
              <button
                className="pagination-btn pagination-arrow"
                onClick={() => handlePageChange(currentPage - 1)}
                disabled={currentPage === 1}
                aria-label="Previous page"
              >
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
                  <path
                    d="M12.5 15L7.5 10L12.5 5"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              </button>
              {getPageNumbers().map((page, index) =>
                page === "ellipsis" ? (
                  <span key={`ellipsis-${index}`} className="pagination-ellipsis">
                    ...
                  </span>
                ) : (
                  <button
                    key={page}
                    className={`pagination-btn ${currentPage === page ? "active" : ""}`}
                    onClick={() => handlePageChange(page as number)}
                  >
                    {page}
                  </button>
                ),
              )}
              <button
                className="pagination-btn pagination-arrow"
                onClick={() => handlePageChange(currentPage + 1)}
                disabled={currentPage === totalPages}
                aria-label="Next page"
              >
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
                  <path
                    d="M7.5 15L12.5 10L7.5 5"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default RequestsSection;
