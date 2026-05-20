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
import InvitationCodesModal from "./InvitationCodesModal";

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

        return (
                <div className="section">
                        <div className="section-actions">
                                <div className="department-members-actions">
                                        <div>
                                                <ButtonPrimary
                                                        text="Adaugă membru"
                                                        fullWidth
                                                        onClick={() => setIsAddModalOpen(true)}
                                                />
                                        </div>
                                        <div>
                                                <ButtonPrimary
                                                        text="Vezi șabloane"
                                                        fullWidth
                                                        variant="ghost"
                                                        onClick={() =>
                                                                router.push(
                                                                        `/dashboard/templates?department=${departmentId}`,
                                                                )
                                                        }
                                                />
                                        </div>
                                        <div>
                                                <ButtonPrimary
                                                        text="Coduri invitație"
                                                        fullWidth
                                                        variant="ghost"
                                                        onClick={() => {setIsCodesModalOpen(true)}}
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

                        <InvitationCodesModal
                                isOpen={isCodesModalOpen}
                                onClose={() => setIsCodesModalOpen(false)}
                        />
                </div>
        );
}
