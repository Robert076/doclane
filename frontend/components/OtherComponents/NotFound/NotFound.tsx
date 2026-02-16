import "./NotFound.css";

interface NotFoundProps {
  text: string;
  subtext?: string;
  background?: string;
}

const NotFound = ({ text, subtext, background }: NotFoundProps) => {
  return (
    <div className="not-found-container" style={{ background: background }}>
      <p className="not-found-text">{text}</p>
      {subtext && <p className="not-found-subtext">{subtext}</p>}
    </div>
  );
};

export default NotFound;
