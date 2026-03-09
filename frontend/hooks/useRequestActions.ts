import { closeRequest, reopenRequest } from "@/lib/api/requests";
import { useState } from "react";

import toast from "react-hot-toast";

export const useRequestActions = (requestId: number) => {
        const [loading, setLoading] = useState(false);

        const closeReq = async () => {
                setLoading(true);
                await toast.promise(closeRequest(requestId), {
                        loading: "Se închide dosarul...",
                        success: (res) => {
                                setLoading(false);
                                if (!res.success) throw new Error(res.message);
                                return "Dosar închis cu success.";
                        },
                        error: (err) => {
                                setLoading(false);
                                return "Ceva nu a mers bine.";
                        },
                });
        };

        const reopenReq = async () => {
                setLoading(true);
                await toast.promise(reopenRequest(requestId), {
                        loading: "Se deschide dosarul...",
                        success: (res) => {
                                setLoading(false);
                                if (!res.success) throw new Error(res.message);
                                return "Dosar redeschis.";
                        },
                        error: (err) => {
                                setLoading(false);
                                return "Ceva nu a mers bine.";
                        },
                });
        };

        return { closeReq, reopenReq, loading };
};
