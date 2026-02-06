import "./NotFound.css";

interface NotFoundProps {
  text: string;
  subtext?: string;
}

const NotFound = ({ text, subtext }: NotFoundProps) => {
  return (
    <div className="not-found-container">
      <p className="not-found-text">{text}</p>
      {subtext && <p className="not-found-subtext">{subtext}</p>}
    </div>
  );
};

export default NotFound;
