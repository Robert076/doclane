"use client";

import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Input from "@/components/InputComponents/Input";
import TextArea from "@/components/InputComponents/TextArea";
import { useState } from "react";
import "./AddRequestForClientForm.css";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import CronInput from "./CronInput";
import { RecurrenceUnit } from "@/types";
import RadioInput from "@/components/InputComponents/RadioInput";
import DeadlineInput from "./DeadlineInput";
import { MdCalendarMonth } from "react-icons/md";
import ScheduleModal from "./ScheduleModal";

interface AddRequestForClientFormProps {
        id: string;
}

const AddRequestForClientForm: React.FC<AddRequestForClientFormProps> = ({ id }) => {
        const [requestName, setRequestName] = useState("");
        const [requestDescription, setRequestDescription] = useState("");

        const [isNoneSelected, setIsNoneSelected] = useState(true);
        const [isRecurring, setIsRecurring] = useState(false);
        const [isDeadline, setIsDeadline] = useState(false);
        const [showScheduleModal, setShowScheduleModal] = useState(false);
        const [scheduledDate, setScheduledDate] = useState("");

        const [unit, setUnit] = useState<RecurrenceUnit>("month");
        const [hour, setHour] = useState("09");
        const [minute, setMinute] = useState("00");
        const [dueDate, setDueDate] = useState("");
        const router = useRouter();

        const handleSubmit = async (e: React.FormEvent) => {
                e.preventDefault();

                if (!requestName) return;

                const recurrenceCron = isRecurring ? buildCron(hour, minute, unit) : undefined;
                const dueDateRFC3339 =
                        dueDate && isDeadline ? new Date(dueDate).toISOString() : undefined;

                const createRequestPromise = fetch("/api/backend/document-requests", {
                        method: "POST",
                        credentials: "include",
                        headers: {
                                "Content-Type": "application/json",
                        },
                        body: JSON.stringify({
                                title: requestName,
                                description: requestDescription,
                                client_id: +id,
                                ...(recurrenceCron && isRecurring
                                        ? { recurrence_cron: recurrenceCron }
                                        : {}),
                                ...(dueDate && isDeadline ? { due_date: dueDateRFC3339 } : {}),
                        }),
                }).then(async (res) => {
                        if (!res.ok) {
                                const errorData = await res.json();
                                throw new Error(errorData.error || "Failed to create request");
                        }
                        return res.json();
                });

                toast.promise(createRequestPromise, {
                        loading: "Creating request...",
                        success: "Request created successfully!",
                        error: (err) => `Failed: ${err.message}`,
                });

                createRequestPromise.then(() => {
                        setRequestName("");
                        setRequestDescription("");
                        router.push("/dashboard/clients");
                });
        };

        const handleScheduleConfirm = (scheduledDateValue: string) => {
                if (!requestName) {
                        toast.error("Please fill in the request title");
                        return;
                }

                const recurrenceCron = isRecurring ? buildCron(hour, minute, unit) : undefined;
                const dueDateRFC3339 =
                        dueDate && isDeadline ? new Date(dueDate).toISOString() : undefined;
                const scheduledDateRFC3339 = new Date(scheduledDateValue).toISOString();

                const createRequestPromise = fetch("/api/backend/document-requests", {
                        method: "POST",
                        credentials: "include",
                        headers: {
                                "Content-Type": "application/json",
                        },
                        body: JSON.stringify({
                                title: requestName,
                                description: requestDescription,
                                client_id: +id,
                                is_scheduled: true,
                                scheduled_for: scheduledDateRFC3339,
                                ...(recurrenceCron && isRecurring
                                        ? { recurrence_cron: recurrenceCron }
                                        : {}),
                                ...(dueDate && isDeadline ? { due_date: dueDateRFC3339 } : {}),
                        }),
                }).then(async (res) => {
                        if (!res.ok) {
                                const errorData = await res.json();
                                throw new Error(
                                        errorData.error || "Failed to schedule request",
                                );
                        }
                        return res.json();
                });

                toast.promise(createRequestPromise, {
                        loading: "Scheduling request...",
                        success: "Request scheduled successfully!",
                        error: (err) => `Failed: ${err.message}`,
                });

                createRequestPromise.then(() => {
                        setRequestName("");
                        setRequestDescription("");
                        setScheduledDate("");
                        router.push("/dashboard/clients");
                });
        };

        return (
                <>
                        <form className="add-request-for-client-form" onSubmit={handleSubmit}>
                                <Input
                                        label="Request title"
                                        placeholder="Title goes here..."
                                        value={requestName}
                                        onChange={(e: any) => setRequestName(e.target.value)}
                                />

                                <TextArea
                                        label="Request description"
                                        value={requestDescription}
                                        placeholder="Description goes here..."
                                        onChange={(e: any) =>
                                                setRequestDescription(e.target.value)
                                        }
                                />

                                <div className="radio-inputs-time">
                                        <RadioInput
                                                isChecked={isNoneSelected}
                                                onChange={(
                                                        e: React.ChangeEvent<HTMLInputElement>,
                                                ) => {
                                                        setIsNoneSelected(e.target.checked);
                                                        setIsRecurring(false);
                                                        setIsDeadline(false);
                                                }}
                                                label="No time constraint"
                                        />
                                        <RadioInput
                                                isChecked={isRecurring}
                                                onChange={(
                                                        e: React.ChangeEvent<HTMLInputElement>,
                                                ) => {
                                                        setIsRecurring(e.target.checked);
                                                        setIsNoneSelected(false);
                                                        setIsDeadline(false);
                                                }}
                                                label="Recurring"
                                        />
                                        <RadioInput
                                                isChecked={isDeadline}
                                                onChange={(
                                                        e: React.ChangeEvent<HTMLInputElement>,
                                                ) => {
                                                        setIsDeadline(e.target.checked);
                                                        setIsRecurring(false);
                                                        setIsNoneSelected(false);
                                                }}
                                                label="Deadline"
                                        />
                                </div>

                                {isRecurring && (
                                        <CronInput
                                                unit={unit}
                                                setUnit={setUnit}
                                                hour={hour}
                                                minute={minute}
                                                setHour={setHour}
                                                setMinute={setMinute}
                                        />
                                )}

                                {isDeadline && (
                                        <DeadlineInput
                                                dueDate={dueDate}
                                                setDueDate={setDueDate}
                                        />
                                )}
                        </form>

                        <div className="button-group">
                                <ButtonPrimary
                                        text="Create request"
                                        onClick={handleSubmit}
                                        type="button"
                                />
                                <ButtonPrimary
                                        text="Schedule request"
                                        variant="ghost"
                                        icon={MdCalendarMonth}
                                        type="button"
                                        onClick={() => setShowScheduleModal(true)}
                                />
                        </div>

                        {showScheduleModal && (
                                <ScheduleModal
                                        onClose={() => setShowScheduleModal(false)}
                                        onConfirm={handleScheduleConfirm}
                                />
                        )}
                </>
        );
};

export default AddRequestForClientForm;

const buildCron = (hour: string, minute: string, unit: RecurrenceUnit) => {
        const h = hour || "0";
        const m = minute || "0";

        switch (unit) {
                case "day":
                        return `${m} ${h} * * *`;
                case "week":
                        return `${m} ${h} * * 1`;
                case "month":
                        return `${m} ${h} 1 * *`;
                case "year":
                        return `${m} ${h} 1 1 *`;
        }
};
