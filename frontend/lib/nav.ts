import { IconType } from "react-icons";
import {
        MdDashboard,
        MdPeople,
        MdSettings,
        MdHome,
        MdArchive,
        MdHelp,
        MdHelpCenter,
        MdAddBox,
        MdEditDocument,
        MdEditRoad,
        MdEditSquare,
        MdUnarchive,
        MdOutlineArchive,
        MdOutlineSendAndArchive,
} from "react-icons/md";

export interface NavItem {
        label: string;
        href: string;
        icon: IconType;
}

export const SIDEBAR_CONFIG: Record<string, NavItem[]> = {
        PROFESSIONAL: [
                { label: "Dosare", href: "/dashboard/requests", icon: MdDashboard },
                {
                        label: "Dosare arhivate",
                        href: "/dashboard/archived-requests",
                        icon: MdArchive,
                },
                {
                        label: "Şabloane",
                        href: "/dashboard/templates",
                        icon: MdEditDocument,
                },
                {
                        label: "Şabloane arhivate",
                        href: "/dashboard/archived-templates",
                        icon: MdOutlineSendAndArchive,
                },
                { label: "Solicitanți", href: "/dashboard/clients", icon: MdPeople },
                { label: "Setări", href: "/dashboard/settings", icon: MdSettings },
                { label: "Ajutor", href: "/dashboard/help", icon: MdHelpCenter },
        ],
        CLIENT: [
                { label: "My Portal", href: "/dashboard", icon: MdHome },
                { label: "Settings", href: "/dashboard/settings", icon: MdSettings },
        ],
};
