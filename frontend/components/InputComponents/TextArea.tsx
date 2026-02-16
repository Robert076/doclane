import "./TextArea.css";

interface TextAreaProps {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
  label?: string;
  placeholder?: string;
}

const TextArea: React.FC<TextAreaProps> = ({ value, onChange, label, placeholder }) => {
  return (
    <div className="textarea-wrapper">
      {label && <label>{label}</label>}
      <textarea value={value} onChange={onChange} placeholder={placeholder} />
    </div>
  );
};

export default TextArea;
