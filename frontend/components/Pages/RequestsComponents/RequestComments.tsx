"use client";
import { useTransition, useState } from "react";
import { useRouter } from "next/navigation";
import SectionTitle from "./SectionTitle";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import TextArea from "@/components/InputComponents/TextArea";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import { RequestComment } from "@/types";
import { addComment } from "@/lib/api/requests";
import { formatDate } from "@/lib/client/formatDate";
import toast from "react-hot-toast";
import "./RequestComments.css";

const COMMENTS_PER_PAGE = 5;

interface RequestCommentsProps {
        comments: RequestComment[];
        requestId: number;
}

export default function RequestComments({
        comments: initial,
        requestId,
}: RequestCommentsProps) {
        const [text, setText] = useState("");
        const [isPending, startTransition] = useTransition();
        const [currentPage, setCurrentPage] = useState(1);
        const router = useRouter();

        const sorted = [...initial].sort(
                (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
        );

        const totalPages = Math.ceil(sorted.length / COMMENTS_PER_PAGE);
        const startIndex = (currentPage - 1) * COMMENTS_PER_PAGE;
        const currentComments = sorted.slice(startIndex, startIndex + COMMENTS_PER_PAGE);

        const handleSubmit = () => {
                const trimmed = text.trim();
                if (trimmed.length < 3) {
                        toast.error("Comentariul trebuie să aibă cel puțin 3 caractere.");
                        return;
                }
                startTransition(async () => {
                        const res = await addComment(requestId, trimmed);
                        if (res.success) {
                                toast.success("Comentariu adăugat.");
                                setText("");
                                router.refresh();
                        } else {
                                toast.error(res.message);
                        }
                });
        };

        return (
                <div className="comments-section">
                        <SectionTitle text="Comentarii" />
                        <div className="comment-compose">
                                <TextArea
                                        value={text}
                                        onChange={(e) => setText(e.target.value)}
                                        placeholder="Scrie un comentariu..."
                                        fullWidth={true}
                                        minHeight={80}
                                        maxHeight={200}
                                />
                                <div className="comment-compose-footer">
                                        <span className="comment-char-count">
                                                {text.length}/200
                                        </span>
                                        <ButtonPrimary
                                                text="Adaugă"
                                                disabled={isPending}
                                                onClick={handleSubmit}
                                                variant="ghost"
                                        />
                                </div>
                        </div>

                        {initial.length === 0 ? (
                                <NotFound
                                        text="Nu există niciun comentariu încă"
                                        subtext="Începe prin a adăuga primul comentariu."
                                />
                        ) : (
                                <>
                                        <ul className="comments-list">
                                                {currentComments.map((c) => (
                                                        <li
                                                                key={c.id}
                                                                className="comment-item"
                                                        >
                                                                <div className="comment-header">
                                                                        <span className="comment-author">
                                                                                {
                                                                                        c.user_first_name
                                                                                }{" "}
                                                                                {
                                                                                        c.user_last_name
                                                                                }
                                                                        </span>
                                                                        <span className="comment-date">
                                                                                {formatDate(
                                                                                        c.created_at,
                                                                                )}
                                                                        </span>
                                                                </div>
                                                                <p className="comment-body">
                                                                        {c.comment}
                                                                </p>
                                                        </li>
                                                ))}
                                        </ul>
                                        {totalPages > 1 && (
                                                <PaginationFooter
                                                        currentPage={currentPage}
                                                        totalPages={totalPages}
                                                        setCurrentPage={setCurrentPage}
                                                />
                                        )}
                                </>
                        )}
                </div>
        );
}
