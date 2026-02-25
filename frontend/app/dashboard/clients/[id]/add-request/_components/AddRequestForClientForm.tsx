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
import ExpectedDocumentsList, { ExpectedDocumentInput } from "./ExpectedDocumentsList";
import { createDocumentRequest } from "@/lib/api/api";
import { UI_TEXT } from "@/locales/ro";

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
        const [expectedDocuments, setExpectedDocuments] = useState<ExpectedDocumentInput[]>([
                { title: "", description: "" },
        ]);
        const router = useRouter();

        const validateForm = () => {
                if (!requestName) {
                        toast.error("Please fill in the request title");
                        return false;
                }
                if (expectedDocuments.some((doc) => !doc.title)) {
                        toast.error("Please fill in all expected document titles");
                        return false;
                }
                return true;
        };

        const buildPayload = (extra?: object) => ({
                title: requestName,
                description: requestDescription,
                client_id: +id,
                expected_documents: expectedDocuments,
                ...(isRecurring && buildCron(hour, minute, unit)
                        ? { recurrence_cron: buildCron(hour, minute, unit) }
                        : {}),
                ...(dueDate && isDeadline
                        ? { due_date: new Date(dueDate).toISOString() }
                        : {}),
                ...extra,
        });

        const handleSubmit = async (e: React.FormEvent) => {
                e.preventDefault();
                if (!validateForm()) return;

                toast.promise(createDocumentRequest(buildPayload()), {
                        loading: "Creating request...",
                        success: (res) => {
                                if (!res.success) throw new Error(res.error);
                                router.push("/dashboard/clients");
                                return "Request created successfully!";
                        },
                        error: (err) => `Failed: ${err.message}`,
                });
        };

        const handleScheduleConfirm = (scheduledDateValue: string) => {
                if (!validateForm()) return;

                toast.promise(
                        createDocumentRequest(
                                buildPayload({
                                        is_scheduled: true,
                                        scheduled_for: new Date(
                                                scheduledDateValue,
                                        ).toISOString(),
                                }),
                        ),
                        {
                                loading: "Scheduling request...",
                                success: (res) => {
                                        if (!res.success) throw new Error(res.error);
                                        router.push("/dashboard/clients");
                                        return "Request scheduled successfully!";
                                },
                                error: (err) => `Failed: ${err.message}`,
                        },
                );
        };

        return (
                <>
                        <form className="add-request-for-client-form" onSubmit={handleSubmit}>
                                <Input
                                        label={UI_TEXT.request.createForm.title}
                                        placeholder={
                                                UI_TEXT.request.createForm.titlePlaceholder
                                        }
                                        value={requestName}
                                        onChange={(e: any) => setRequestName(e.target.value)}
                                />
                                <TextArea
                                        label={UI_TEXT.request.createForm.description}
                                        value={requestDescription}
                                        placeholder={
                                                UI_TEXT.request.createForm
                                                        .descriptionPlaceholder
                                        }
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
                                                label={
                                                        UI_TEXT.request.createForm.time
                                                                .noConstraint
                                                }
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
                                                label={
                                                        UI_TEXT.request.createForm.time
                                                                .recurring
                                                }
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
                                                label={
                                                        UI_TEXT.request.createForm.time
                                                                .deadline
                                                }
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
                                <ExpectedDocumentsList
                                        documents={expectedDocuments}
                                        onChange={setExpectedDocuments}
                                />
                        </form>
                        <div className="button-group">
                                <ButtonPrimary
                                        text={UI_TEXT.request.createForm.createRequest}
                                        onClick={handleSubmit}
                                        type="button"
                                />
                                <ButtonPrimary
                                        text={UI_TEXT.request.createForm.scheduleRequest}
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
