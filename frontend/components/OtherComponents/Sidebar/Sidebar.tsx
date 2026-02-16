"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useUser } from "@/context/UserContext";
import { SIDEBAR_CONFIG } from "@/lib/nav";
import { logout } from "@/lib/api/api";

import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Separator from "@/components/OtherComponents/Separators/Separator/Separator";
import "./Sidebar.css";
import Logo from "@/components/OtherComponents/Logo/Logo";

export default function Sidebar() {
        const pathname = usePathname();
        const user = useUser();

        const links = SIDEBAR_CONFIG[user.role] || [];

        return (
                <aside className="sidebar">
                        <Logo />
                        <Separator />

                        <nav className="sidebar-nav">
                                {links.map((link) => {
                                        const isActive = pathname === link.href;
                                        const Icon = link.icon;

                                        return (
                                                <Link
                                                        key={link.href}
                                                        href={link.href}
                                                        className={`sidebar-link ${isActive ? "active" : ""}`}
                                                >
                                                        {Icon && (
                                                                <Icon
                                                                        size={20}
                                                                        className="sidebar-link-icon"
                                                                />
                                                        )}
                                                        <span className="sidebar-link-label">
                                                                {link.label}
                                                        </span>
                                                </Link>
                                        );
                                })}
                        </nav>

                        <div className="sidebar-footer">
                                <div className="sidebar-user-info">
                                        <p className="sidebar-user-email">{user.email}</p>
                                        <span className="sidebar-user-role">
                                                {user.role.toLowerCase()}
                                        </span>
                                </div>

                                <ButtonPrimary
                                        text="Log out"
                                        variant="primary"
                                        fullWidth={true}
                                        onClick={logout}
                                />
                        </div>
                </aside>
        );
}
