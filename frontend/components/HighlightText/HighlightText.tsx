import React from "react";
import "./HighlightText.css";

interface HighlightTextProps {
  text: string;
  search?: string;
  className?: string;
}

const HighlightText: React.FC<HighlightTextProps> = ({ text, search, className = "" }) => {
  if (!search || search.trim() === "") {
    return <span className={className}>{text}</span>;
  }

  // Escape special regex characters
  const escapeRegex = (str: string) => {
    return str.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  };

  try {
    // Create regex pattern for case-insensitive matching
    const pattern = new RegExp(`(${escapeRegex(search)})`, "gi");
    const parts = text.split(pattern);

    return (
      <span className={className}>
        {parts.map((part, index) => {
          // Check if this part matches the search term (case-insensitive)
          if (part.toLowerCase() === search.toLowerCase()) {
            return (
              <mark key={index} className="highlight-match">
                {part}
              </mark>
            );
          }
          return <React.Fragment key={index}>{part}</React.Fragment>;
        })}
      </span>
    );
  } catch (error) {
    // If regex fails, return original text
    return <span className={className}>{text}</span>;
  }
};

export default HighlightText;
