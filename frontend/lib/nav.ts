import { IconType } from "react-icons";
import { MdDashboard, MdPeople, MdSettings, MdHome } from "react-icons/md";

export interface NavItem {
        label: string;
        href: string;
        icon: IconType;
}

export const SIDEBAR_CONFIG: Record<string, NavItem[]> = {
        PROFESSIONAL: [
                { label: "Dosare", href: "/dashboard", icon: MdDashboard },
                { label: "Solicitanți", href: "/dashboard/clients", icon: MdPeople },
                { label: "Setări", href: "/dashboard/settings", icon: MdSettings },
        ],
        CLIENT: [
                { label: "My Portal", href: "/dashboard", icon: MdHome },
                { label: "Settings", href: "/dashboard/settings", icon: MdSettings },
        ],
};
