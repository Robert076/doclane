import React from "react";
import "./CheckboxInput.css";

interface CheckboxInputProps {
  label: string;
  isChecked: boolean;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

const CheckboxInput: React.FC<CheckboxInputProps> = ({ label, isChecked, onChange }) => {
  return (
    <label className="checkbox-input">
      <input type="checkbox" checked={isChecked} onChange={onChange} />
      <span className="checkmark"></span>
      {label}
    </label>
  );
};

export default CheckboxInput;
