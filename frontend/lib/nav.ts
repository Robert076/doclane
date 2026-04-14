import { IconType } from "react-icons";
import {
        MdDashboard,
        MdPeople,
        MdArchive,
        MdEditDocument,
        MdOutlineSendAndArchive,
        MdApartment,
        MdHome,
        MdDoNotDisturb,
} from "react-icons/md";

export interface NavItem {
        label: string;
        href: string;
        icon: IconType;
}

export const SIDEBAR_CONFIG: Record<string, NavItem[]> = {
        admin: [
                { label: "Dosare", href: "/dashboard/requests", icon: MdDashboard },
                {
                        label: "Dosare arhivate",
                        href: "/dashboard/archived-requests",
                        icon: MdArchive,
                },
                {
                        label: "Dosare retrase",
                        href: "/dashboard/cancelled-requests",
                        icon: MdDoNotDisturb,
                },
                { label: "Şabloane", href: "/dashboard/templates", icon: MdEditDocument },
                {
                        label: "Şabloane arhivate",
                        href: "/dashboard/archived-templates",
                        icon: MdOutlineSendAndArchive,
                },
                { label: "Utilizatori", href: "/dashboard/users", icon: MdPeople },
                { label: "Departamente", href: "/dashboard/departments", icon: MdApartment },
                // { label: "Setări", href: "/dashboard/settings", icon: MdSettings },
                // { label: "Ajutor", href: "/dashboard/help", icon: MdHelpCenter },
        ],
        member: [
                { label: "Portal", href: "/dashboard/requests", icon: MdHome },
                { label: "Şabloane", href: "/dashboard/templates", icon: MdEditDocument },
                // { label: "Setări", href: "/dashboard/settings", icon: MdSettings },
                // { label: "Ajutor", href: "/dashboard/help", icon: MdHelpCenter },
        ],
};
