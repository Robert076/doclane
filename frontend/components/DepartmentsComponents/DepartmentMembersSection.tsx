"use client";
import { useState } from "react";
import { User, InvitationCode, Department } from "@/types";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import UserCard from "@/components/Pages/UsersComponents/UserCard";
import Modal from "@/components/Modals/Modal";
import { deactivateUser } from "@/lib/api/users";
import {
        getInvitationCodesByDepartment,
        deleteInvitationCode,
} from "@/lib/api/invitation_codes";
import { useRouter } from "next/navigation";
import { formatDate } from "@/lib/client/formatDate";
import toast from "react-hot-toast";
import "./DepartmentMembersSection.css";
import GenerateInvitationCodeModal from "./GenerateInvitationCodeModal";
import MoveDepartmentModal from "./MoveDepartmentModal";
import { MdContentCopy } from "react-icons/md";

interface Props {
        members: User[];
        departmentId: number;
        departments: Department[];
}

export default function DepartmentMembersSection({
        members = [],
        departmentId,
        departments,
}: Props) {
        const router = useRouter();
        const [isAddModalOpen, setIsAddModalOpen] = useState(false);
        const [isCodesModalOpen, setIsCodesModalOpen] = useState(false);
        const [codes, setCodes] = useState<InvitationCode[]>([]);
        const [isLoadingCodes, setIsLoadingCodes] = useState(false);
        const [movingUser, setMovingUser] = useState<User | null>(null);

        const handleDeactivate = async (userId: number, name: string) => {
                const id = toast.loading("Se dezactivează...");
                const response = await deactivateUser(userId);
                toast.dismiss(id);
                if (response.success) {
                        toast.success(`${name} a fost dezactivat.`);
                        router.refresh();
                } else {
                        toast.error(response.message);
                }
        };

        const handleOpenCodes = async () => {
                setIsLoadingCodes(true);
                const response = await getInvitationCodesByDepartment(departmentId);
                setIsLoadingCodes(false);
                if (response.success) {
                        setCodes(response.data ?? []);
                        setIsCodesModalOpen(true);
                } else {
                        toast.error(response.message ?? "Eroare la încărcarea codurilor.");
                }
        };

        const handleDeleteCode = async (codeId: number) => {
                const id = toast.loading("Se șterge...");
                const response = await deleteInvitationCode(codeId);
                toast.dismiss(id);
                if (response.success) {
                        toast.success("Cod șters.");
                        setCodes((prev) => prev.filter((c) => c.id !== codeId));
                } else {
                        toast.error(response.message);
                }
        };

        return (
                <div className="section">
                        <div className="section-actions">
                                <div className="department-members-actions">
                                        <div>
                                                <ButtonPrimary
                                                        text="Vezi șabloane"
                                                        fullWidth
                                                        onClick={() =>
                                                                router.push(
                                                                        `/dashboard/templates?department=${departmentId}`,
                                                                )
                                                        }
                                                />
                                        </div>
                                        <div>
                                                <ButtonPrimary
                                                        text="Adaugă membru"
                                                        fullWidth
                                                        onClick={() => setIsAddModalOpen(true)}
                                                />
                                        </div>
                                        <div>
                                                <ButtonPrimary
                                                        text="Coduri invitație"
                                                        fullWidth
                                                        onClick={handleOpenCodes}
                                                        disabled={isLoadingCodes}
                                                />
                                        </div>
                                </div>
                        </div>

                        {members.length === 0 ? (
                                <NotFound
                                        text="Niciun membru în acest departament."
                                        subtext="Atribuie utilizatori acestui departament pentru a-i vedea aici."
                                        background="#fff"
                                />
                        ) : (
                                <div className="objects-grid">
                                        {members.map((member) => (
                                                <UserCard
                                                        key={member.id}
                                                        user={member}
                                                        footer={
                                                                <>
                                                                        <ButtonPrimary
                                                                                text="Mută departament"
                                                                                variant="ghost"
                                                                                fullWidth
                                                                                onClick={() =>
                                                                                        setMovingUser(
                                                                                                member,
                                                                                        )
                                                                                }
                                                                        />
                                                                        <ButtonPrimary
                                                                                text="Dezactivează cont"
                                                                                variant="ghost"
                                                                                fullWidth
                                                                                onClick={() =>
                                                                                        handleDeactivate(
                                                                                                member.id,
                                                                                                `${member.first_name} ${member.last_name}`,
                                                                                        )
                                                                                }
                                                                        />
                                                                </>
                                                        }
                                                />
                                        ))}
                                </div>
                        )}

                        <GenerateInvitationCodeModal
                                isOpen={isAddModalOpen}
                                onClose={() => setIsAddModalOpen(false)}
                                departmentId={departmentId}
                        />

                        {movingUser && (
                                <MoveDepartmentModal
                                        isOpen={!!movingUser}
                                        onClose={() => setMovingUser(null)}
                                        userId={movingUser.id}
                                        userName={`${movingUser.first_name} ${movingUser.last_name}`}
                                        currentDepartmentId={departmentId}
                                        departments={departments}
                                />
                        )}

                        <Modal
                                isOpen={isCodesModalOpen}
                                onClose={() => setIsCodesModalOpen(false)}
                                title="Coduri de invitație"
                                hideFooter
                        >
                                {codes.length === 0 ? (
                                        <p className="codes-empty">
                                                Niciun cod activ pentru acest departament.
                                        </p>
                                ) : (
                                        <div className="codes-list">
                                                {codes.map((code) => {
                                                        const inviteLink = `${window.location.origin}/register/invite?code=${code.code}`;
                                                        return (
                                                                <div
                                                                        key={code.id}
                                                                        className="code-item"
                                                                >
                                                                        <div>
                                                                                <span className="code-item-text">
                                                                                        {
                                                                                                code.code
                                                                                        }
                                                                                </span>
                                                                                {code.expires_at && (
                                                                                        <p className="code-item-meta">
                                                                                                Expiră:{" "}
                                                                                                {formatDate(
                                                                                                        code.expires_at,
                                                                                                )}
                                                                                        </p>
                                                                                )}
                                                                        </div>
                                                                        <button
                                                                                className="code-copy-btn"
                                                                                onClick={() => {
                                                                                        navigator.clipboard.writeText(
                                                                                                inviteLink,
                                                                                        );
                                                                                        toast.success(
                                                                                                "Link copiat!",
                                                                                        );
                                                                                }}
                                                                        >
                                                                                <MdContentCopy
                                                                                        size={
                                                                                                18
                                                                                        }
                                                                                />
                                                                        </button>
                                                                </div>
                                                        );
                                                })}
                                        </div>
                                )}
                        </Modal>
                </div>
        );
}
