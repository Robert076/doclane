import { RecurrenceUnit } from "@/types";

export const buildCronExpression = (
        unit: RecurrenceUnit,
        hour: string,
        minute: string,
): string => {
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

export const parseCronExpression = (
        cron: string,
): { unit: RecurrenceUnit; hour: string; minute: string } => {
        const defaults = { unit: "month" as RecurrenceUnit, hour: "09", minute: "00" };
        if (!cron) return defaults;

        const parts = cron.split(" ");
        if (parts.length !== 5) return defaults;

        const [minute, hour, dom, month, dow] = parts;

        let unit: RecurrenceUnit;
        if (month !== "*") unit = "year";
        else if (dom !== "*") unit = "month";
        else if (dow !== "*") unit = "week";
        else unit = "day";

        return { unit, hour, minute };
};
