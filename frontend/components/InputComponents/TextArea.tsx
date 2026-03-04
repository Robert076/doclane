import "./TextArea.css";

interface TextAreaProps {
        value: string;
        onChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
        label?: string;
        placeholder?: string;
        minHeight?: number;
        maxHeight?: number;
        fullWidth?: boolean;
}

const TextArea: React.FC<TextAreaProps> = ({
        value,
        onChange,
        label,
        placeholder,
        minHeight,
        maxHeight,
        fullWidth,
}) => {
        return (
                <div
                        className={`textarea-wrapper${fullWidth ? " textarea-wrapper--full-width" : ""}`}
                >
                        {label && <label>{label}</label>}
                        <textarea
                                value={value}
                                onChange={onChange}
                                placeholder={placeholder}
                                style={{
                                        ...(minHeight !== undefined && {
                                                minHeight: `${minHeight}px`,
                                        }),
                                        ...(maxHeight !== undefined && {
                                                maxHeight: `${maxHeight}px`,
                                        }),
                                }}
                        />
                </div>
        );
};

export default TextArea;
