import { IconType } from "react-icons";
import {
  MdDashboard,
  MdPeople,
  MdDescription,
  MdSettings,
  MdHome,
  MdFolder,
  MdPerson,
} from "react-icons/md";

export interface NavItem {
  label: string;
  href: string;
  icon: IconType;
}

export const SIDEBAR_CONFIG: Record<string, NavItem[]> = {
  PROFESSIONAL: [
    { label: "Overview", href: "/dashboard", icon: MdDashboard },
    { label: "My Clients", href: "/dashboard/clients", icon: MdPeople },
    { label: "Document Requests", href: "/dashboard/requests", icon: MdDescription },
    { label: "Settings", href: "/dashboard/settings", icon: MdSettings },
  ],
  CLIENT: [
    { label: "My Portal", href: "/dashboard", icon: MdHome },
    { label: "Documents", href: "/dashboard/documents", icon: MdFolder },
    { label: "Professional", href: "/dashboard/my-pro", icon: MdPerson },
  ],
};
