import "./PaginationCounter.css";

interface PaginationCounterProps {
  startIndex: number;
  itemsPerPage: number;
  length: number;
}

const PaginationCounter: React.FC<PaginationCounterProps> = ({
  startIndex,
  itemsPerPage,
  length,
}) => {
  return (
    <span className="pagination-counter">
      {startIndex + 1}-{Math.min(startIndex + itemsPerPage, length)} of {length}
    </span>
  );
};

export default PaginationCounter;
