"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { MdCalendarMonth } from "react-icons/md";

import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Input from "@/components/InputComponents/Input";
import TextArea from "@/components/InputComponents/TextArea";
import RadioInput from "@/components/InputComponents/RadioInput";
import ExpectedDocumentsList, { ExpectedDocumentInput } from "./ExpectedDocumentsList";
import DeadlineInput from "./DeadlineInput";
import CronInput from "./CronInput";
import ScheduleModal from "./ScheduleModal";

// API & Types
import { createDocumentRequest } from "@/lib/api/requests";
import { RecurrenceUnit } from "@/types";
import "./AddRequestForClientForm.css";

interface AddRequestForClientFormProps {
        id: string;
}

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

export default function AddRequestForClientForm({ id }: AddRequestForClientFormProps) {
        const router = useRouter();

        // State-uri principale
        const [requestName, setRequestName] = useState("");
        const [requestDescription, setRequestDescription] = useState("");
        const [expectedDocuments, setExpectedDocuments] = useState<ExpectedDocumentInput[]>([
                { title: "", description: "" },
        ]);

        // State-uri pentru setările de timp
        const [isNoneSelected, setIsNoneSelected] = useState(true);
        const [isRecurring, setIsRecurring] = useState(false);
        const [isDeadline, setIsDeadline] = useState(false);

        // State-uri ajutătoare pentru timp
        const [showScheduleModal, setShowScheduleModal] = useState(false);
        const [unit, setUnit] = useState<RecurrenceUnit>("month");
        const [hour, setHour] = useState("09");
        const [minute, setMinute] = useState("00");
        const [dueDate, setDueDate] = useState("");

        const validateForm = () => {
                if (!requestName.trim()) {
                        toast.error("Te rog să introduci un titlu pentru dosar");
                        return false;
                }
                if (expectedDocuments.some((doc) => !doc.title.trim())) {
                        toast.error("Toate documentele așteptate trebuie să aibă un titlu");
                        return false;
                }
                return true;
        };

        const buildPayload = (extra?: object) => {
                const cron = isRecurring ? buildCron(hour, minute, unit) : undefined;
                const due =
                        isDeadline && dueDate ? new Date(dueDate).toISOString() : undefined;

                return {
                        title: requestName,
                        description: requestDescription,
                        client_id: parseInt(id, 10), // Asigurăm conversia id-ului în număr
                        is_recurring: isRecurring,
                        recurrence_cron: cron,
                        due_date: due,
                        expected_documents: expectedDocuments,
                        ...extra,
                };
        };

        const handleSubmit = async (e?: React.FormEvent) => {
                if (e) e.preventDefault();
                if (!validateForm()) return;

                const loadingToast = toast.loading("Se creează dosarul...");

                try {
                        const res = await createDocumentRequest(buildPayload());

                        if (!res.success) {
                                throw new Error(
                                        res.error ||
                                                res.message ||
                                                "Eroare la crearea dosarului",
                                );
                        }

                        toast.success("Dosar creat cu succes!", { id: loadingToast });
                        router.push("/dashboard/clients");
                } catch (error: any) {
                        toast.error(error.message, { id: loadingToast });
                }
        };

        const handleScheduleConfirm = async (scheduledDateValue: string) => {
                if (!validateForm()) return;

                const loadingToast = toast.loading("Se programează dosarul...");

                try {
                        const res = await createDocumentRequest(
                                buildPayload({
                                        is_scheduled: true,
                                        scheduled_for: new Date(
                                                scheduledDateValue,
                                        ).toISOString(),
                                }),
                        );

                        if (!res.success) {
                                throw new Error(
                                        res.error || res.message || "Eroare la programare",
                                );
                        }

                        toast.success("Dosar programat cu succes!", { id: loadingToast });
                        setShowScheduleModal(false);
                        router.push("/dashboard/clients");
                } catch (error: any) {
                        toast.error(error.message, { id: loadingToast });
                }
        };

        return (
                <>
                        <form className="add-request-for-client-form" onSubmit={handleSubmit}>
                                <Input
                                        label="Titlu dosar"
                                        placeholder="Introdu titlul dosarului"
                                        value={requestName}
                                        onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                                                setRequestName(e.target.value)
                                        }
                                />

                                <TextArea
                                        label="Descriere (Opțional)"
                                        value={requestDescription}
                                        placeholder="Adaugă o scurtă descriere a necesității acestui dosar..."
                                        onChange={(
                                                e: React.ChangeEvent<HTMLTextAreaElement>,
                                        ) => setRequestDescription(e.target.value)}
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
                                                label="Fără limită de timp"
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
                                                label="Recurent (Periodic)"
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
                                                label="Termen limită"
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
                                        text="Creează dosar"
                                        onClick={handleSubmit}
                                        type="button"
                                />
                                <ButtonPrimary
                                        text="Programează pentru mai târziu"
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
}
