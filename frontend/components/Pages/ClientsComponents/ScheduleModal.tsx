"use client";

import { useState } from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { MdClose } from "react-icons/md";
import "./ScheduleModal.css";
import { UI_TEXT } from "@/locales/ro";

interface ScheduleModalProps {
        onClose: () => void;
        onConfirm: (scheduledDate: string) => void;
}

const ScheduleModal: React.FC<ScheduleModalProps> = ({ onClose, onConfirm }) => {
        const [scheduledDate, setScheduledDate] = useState("");

        const handleConfirm = () => {
                if (!scheduledDate) {
                        return;
                }
                onConfirm(scheduledDate);
                onClose();
        };

        return (
                <div className="modal-overlay" onClick={onClose}>
                        <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                                <div className="modal-header">
                                        <h3>{UI_TEXT.modals.scheduleRequest.title}</h3>
                                        <button className="modal-close" onClick={onClose}>
                                                <MdClose />
                                        </button>
                                </div>

                                <div className="modal-body">
                                        <label htmlFor="scheduled-date">
                                                {UI_TEXT.modals.scheduleRequest.subtitle}
                                        </label>
                                        <input
                                                id="scheduled-date"
                                                type="datetime-local"
                                                value={scheduledDate}
                                                onChange={(e) =>
                                                        setScheduledDate(e.target.value)
                                                }
                                                min={new Date().toISOString().slice(0, 16)}
                                        />
                                </div>

                                <div className="modal-footer">
                                        <ButtonPrimary
                                                text={UI_TEXT.common.cancel}
                                                variant="ghost"
                                                onClick={onClose}
                                                type="button"
                                        />
                                        <ButtonPrimary
                                                text={UI_TEXT.modals.scheduleRequest.schedule}
                                                onClick={handleConfirm}
                                                type="button"
                                        />
                                </div>
                        </div>
                </div>
        );
};

export default ScheduleModal;
