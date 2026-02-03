import { RecurrenceUnit } from "@/types";
import React, { Dispatch, SetStateAction } from "react";
import "./CronInput.css";

interface CronInputProps {
  unit: RecurrenceUnit;
  setUnit: Dispatch<SetStateAction<RecurrenceUnit>>;
  hour: string;
  setHour: Dispatch<SetStateAction<string>>;
  minute: string;
  setMinute: Dispatch<SetStateAction<string>>;
}

const CronInput: React.FC<CronInputProps> = ({
  unit,
  setUnit,
  hour,
  setHour,
  minute,
  setMinute,
}) => {
  return (
    <div className="recurrence-ui">
      <label>
        Every
        <select value={unit} onChange={(e) => setUnit(e.target.value as RecurrenceUnit)}>
          <option value="day">Day</option>
          <option value="week">Week</option>
          <option value="month">Month</option>
          <option value="year">Year</option>
        </select>
      </label>

      <label>
        At
        <input
          type="number"
          min="0"
          max="23"
          value={hour}
          onChange={(e) => setHour(e.target.value)}
        />
        :
        <input
          type="number"
          min="0"
          max="59"
          value={minute}
          onChange={(e) => setMinute(e.target.value)}
        />
      </label>
    </div>
  );
};

export default CronInput;
