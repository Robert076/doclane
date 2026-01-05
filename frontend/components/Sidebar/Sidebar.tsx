"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { User } from "@/types";
import { SIDEBAR_CONFIG } from "@/config/nav";
import "./Sidebar.css";
import Logo from "../Logo/Logo";
import Button from "../Button/Button";
import Separator from "../Separators/Separator/Separator";
import { logout } from "@/lib/logout";

export default function Sidebar({ user }: { user: User }) {
  const pathname = usePathname();
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
              <Icon size={20} className="sidebar-link-icon" />
              {link.label}
            </Link>
          );
        })}
      </nav>

      <div className="sidebar-footer">
        <div className="sidebar-user-info">
          <p className="sidebar-user-email">{user.email}</p>
          <span className="sidebar-user-role">{user.role}</span>
        </div>
        <Button text="Log out" onClick={logout} />
      </div>
    </aside>
  );
}
