import React, { useState, Dispatch, SetStateAction } from "react";
import "./DeadlineInput.css";

interface DeadlineInputProps {
  dueDate: string;
  setDueDate: Dispatch<SetStateAction<string>>;
}

const DeadlineInput: React.FC<DeadlineInputProps> = ({ dueDate, setDueDate }) => {
  return (
    <div className="deadline-ui">
      <label>
        Due date
        <input type="date" value={dueDate} onChange={(e) => setDueDate(e.target.value)} />
      </label>
    </div>
  );
};

export default DeadlineInput;
