import React from "react";
import "./RadioInput.css";

interface RadioInputProps {
  label: string;
  isChecked: boolean;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

const RadioInput: React.FC<RadioInputProps> = ({ label, isChecked, onChange }) => {
  return (
    <label className="checkbox-input">
      <input type="radio" checked={isChecked} onChange={onChange} />
      <span className="checkmark"></span>
      {label}
    </label>
  );
};

export default RadioInput;
