import React from "react";
import "./SeparatorWithText.css";

interface SeparatorWithTextProps {
  text: string;
}

const SeparatorWithText: React.FC<SeparatorWithTextProps> = ({ text }) => {
  return (
    <div className="separator-with-text">
      <span className="separator-line" />
      <span className="separator-text">{text}</span>
      <span className="separator-line" />
    </div>
  );
};

export default SeparatorWithText;
