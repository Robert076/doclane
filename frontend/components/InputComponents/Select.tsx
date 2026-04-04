import "./Select.css";

interface SelectOption {
        value: string | number;
        label: string;
}

interface SelectProps {
        value: string | number;
        onChange: (value: string) => void;
        options: SelectOption[];
        placeholder?: string;
        label?: string;
}

export default function Select({ value, onChange, options, placeholder, label }: SelectProps) {
        return (
                <div className="select-wrapper">
                        {label && <label className="select-label">{label}</label>}
                        <select
                                className="select"
                                value={value}
                                onChange={(e) => onChange(e.target.value)}
                        >
                                {placeholder && <option value="">{placeholder}</option>}
                                {options.map((opt) => (
                                        <option key={opt.value} value={opt.value}>
                                                {opt.label}
                                        </option>
                                ))}
                        </select>
                </div>
        );
}
