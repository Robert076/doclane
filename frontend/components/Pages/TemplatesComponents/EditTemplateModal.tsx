"use client";
import { useState } from "react";
import { Template, RecurrenceUnit } from "@/types";
import Modal from "@/components/Modals/Modal";
import Input from "@/components/InputComponents/Input";
import TextArea from "@/components/InputComponents/TextArea";
import RadioInput from "@/components/InputComponents/RadioInput";
import CronInput from "@/components/Pages/ClientsComponents/CronInput";
import { buildCronExpression, parseCronExpression } from "@/lib/cron";

interface EditTemplateModalProps {
        isOpen: boolean;
        onClose: () => void;
        onConfirm: (data: {
                title?: string;
                description?: string;
                is_recurring?: boolean;
                recurrence_cron?: string;
        }) => void;
        template: Template;
}

export default function EditTemplateModal({
        isOpen,
        onClose,
        onConfirm,
        template,
}: EditTemplateModalProps) {
        const parsed = parseCronExpression(template.recurrence_cron ?? "");

        const [title, setTitle] = useState(template.title);
        const [description, setDescription] = useState(template.description ?? "");
        const [isNoneSelected, setIsNoneSelected] = useState(!template.is_recurring);
        const [isRecurring, setIsRecurring] = useState(template.is_recurring);
        const [unit, setUnit] = useState<RecurrenceUnit>(parsed.unit);
        const [hour, setHour] = useState(parsed.hour);
        const [minute, setMinute] = useState(parsed.minute);

        const handleConfirm = () => {
                const payload: {
                        title?: string;
                        description?: string;
                        is_recurring?: boolean;
                        recurrence_cron?: string;
                } = {};

                if (title !== template.title) payload.title = title;
                if (description !== (template.description ?? ""))
                        payload.description = description;
                if (isRecurring !== template.is_recurring) payload.is_recurring = isRecurring;

                if (isRecurring) {
                        const cron = buildCronExpression(unit, hour, minute);
                        if (cron !== template.recurrence_cron) payload.recurrence_cron = cron;
                }

                onConfirm(payload);
        };

        return (
                <Modal
                        isOpen={isOpen}
                        onClose={onClose}
                        onConfirm={handleConfirm}
                        title="Editează șablon"
                >
                        <Input
                                label="Titlu"
                                value={title}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                                        setTitle(e.target.value)
                                }
                                placeholder="Titlul șablonului"
                        />
                        <TextArea
                                label="Descriere"
                                value={description}
                                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                                        setDescription(e.target.value)
                                }
                                placeholder="Descrierea șablonului"
                                fullWidth
                                minHeight={100}
                        />
                        <div className="radio-inputs-time">
                                <RadioInput
                                        isChecked={isNoneSelected}
                                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                                setIsNoneSelected(e.target.checked);
                                                setIsRecurring(false);
                                        }}
                                        label="Fără recurenţă"
                                />
                                <RadioInput
                                        isChecked={isRecurring}
                                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                                setIsRecurring(e.target.checked);
                                                setIsNoneSelected(false);
                                        }}
                                        label="Recurent"
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
                </Modal>
        );
}
