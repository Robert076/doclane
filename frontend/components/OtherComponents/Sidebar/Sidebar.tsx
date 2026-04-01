"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState, useEffect } from "react";
import { useUser } from "@/context/UserContext";
import { SIDEBAR_CONFIG } from "@/lib/nav";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Separator from "@/components/OtherComponents/Separators/Separator/Separator";
import Logo from "@/components/OtherComponents/Logo/Logo";
import { UI_TEXT } from "@/locales/ro";
import { logout } from "@/lib/api/auth";
import "./Sidebar.css";

export default function Sidebar() {
        const pathname = usePathname();
        const user = useUser();
        const links = SIDEBAR_CONFIG[user.role] ?? [];

        const [isOpen, setIsOpen] = useState(false);

        useEffect(() => {
                setIsOpen(false);
        }, [pathname]);

        useEffect(() => {
                document.body.style.overflow = isOpen ? "hidden" : "";
                return () => {
                        document.body.style.overflow = "";
                };
        }, [isOpen]);

        return (
                <>
                        <button
                                className="sidebar-hamburger"
                                onClick={() => setIsOpen((prev) => !prev)}
                                aria-label={isOpen ? "Close menu" : "Open menu"}
                                aria-expanded={isOpen}
                        >
                                <span className={`hamburger-bar ${isOpen ? "open" : ""}`} />
                                <span className={`hamburger-bar ${isOpen ? "open" : ""}`} />
                                <span className={`hamburger-bar ${isOpen ? "open" : ""}`} />
                        </button>

                        {isOpen && (
                                <div
                                        className="sidebar-overlay"
                                        onClick={() => setIsOpen(false)}
                                        aria-hidden="true"
                                />
                        )}

                        <aside className={`sidebar ${isOpen ? "sidebar--open" : ""}`}>
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
                                                <p className="sidebar-user-email">
                                                        {user.email}
                                                </p>
                                                <span className="sidebar-user-role">
                                                        {user.role === "admin"
                                                                ? "Administrator"
                                                                : "Membru"}
                                                </span>
                                        </div>
                                        <ButtonPrimary
                                                text={UI_TEXT.sidebar.logout}
                                                variant="primary"
                                                fullWidth={true}
                                                onClick={logout}
                                        />
                                </div>
                        </aside>
                </>
        );
}
