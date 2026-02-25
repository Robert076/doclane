export const UI_TEXT = {
        roles: {
                clientSingular: "Solicitant",
                clientPlural: "Solicitanți",
                professionalSingular: "Funcționar",
                professionalPlural: "Funcționari",
        },
        common: {
                loading: "Se încarcă...",
                search: "Caută...",
                searchNotFound:
                        "Nu am găsit nimic să se potrivească cu căutarea dumneavoastră.",
                save: "Salvează",
                cancel: "Anulează",
                close: "Închide",
                continue: "Continuă",
                done: "Gata",
                status: "Status",
                dueDate: "Data scadentă",
                nextDueAt: "Următoarea scadență",
                description: "Descriere",
                createdAt: "Deschis la data",
                updatedAt: "Închis la data",
                notFoundTitleRequests: "Nu am găsit niciun dosar in lucru.",
                notFoundSubtitleRequestsProfessional: "Începe prin a crea primul dosar.",
                notFoundSubtitleRequestsClient:
                        "Nu este necesară nicio acțiune din partea dumneavoastră.",
        },
        dashboard: {
                professional: {
                        header: (name: string) => `Bine ai revenit, ${name}`,
                        subheader: (activeRequestsCount: number) =>
                                `Ai ${activeRequestsCount} dosare în lucru.`,
                },
        },
        request: {
                card: {
                        clientEmail: "Adresa de email",
                        clientName: "Numele solicitantului",
                },
                status: {
                        overdue: "În întârziere",
                },
                actions: {
                        viewDetails: "Vezi detalii",
                        closeRequest: "Arhivează dosarul",
                },
                details: {
                        title: "Detaliile dosarului",
                        actions: "Acțiuni",
                        files: "Documente",
                },
                createForm: {
                        title: "Titlul dosarului",
                        titlePlaceholder: "Scrie aici titlul dosarului...",
                        description: "Descrierea dosarului",
                        descriptionPlaceholder: "Scrie aici descrierea dosarului...",
                        time: {
                                noConstraint: "Fără constrângeri de timp",
                                recurring: "Recurent",
                                deadline: "Termen limită",
                        },
                        expectedDocuments: "Documente solicitate",
                        expectedDocumentsNotAdded: "Niciun document solicitat încă",
                        expectedDocumentTitle: "Titlul documentului solicitat",
                        expectedDocumentTitlePlaceholder: "Copia după buletin...",
                        expectedDocumentDescription: "Descrierea documentului solicitant",
                        expectedDocumentDescriptionPlaceholder:
                                "Poza să fie nemişcată si clară...",
                        addExpectedDocument: "Adaugă document",
                        createRequest: "Crează dosar",
                        scheduleRequest: "Programează dosar",
                },
        },
        client: {
                card: {
                        clientEmail: "Adresa de email",
                        joinedAt: "Cont creat",
                },
                actions: {
                        newRequest: "Dosar nou",
                        deactivateAccount: "Dezactivează contul",
                },
        },
        modals: {
                closeRequest: {
                        title: "Arhivează dosarul",
                        subtitle1: "Eşti sigur că vrei sa arhivezi dosarul ",
                        subtitle2: "Această acțiune va marca dosarul ca arhivat, ceea ce înseamnă că solicitantul nu va mai putea adăuga documentele. Acțiunea este reversibilă.",
                },
                scheduleRequest: {
                        title: "Programează dosar",
                        subtitle: "Dosarul va fi disponibil clientului începand cu data selectată.",
                        schedule: "Programează",
                },
                generateCode: {
                        title: "Crează un nou cod de invitație",
                        subtitle1: `
                                Această acțiune va genera un cod unic de invitație. Trimite acest cod noului solicitant ca el să se poată ìnregistra in Doclane.
                        `,
                        subtitle2: `
                                Asigură-te ca îi trimiți codul cât mai repede posibil, acesta va expira în 7 zile.
                        `,
                        subtitle3: `
                                Distribuie acest cod noului solicitant. Acest cod este de unică folosință.
                        `,
                        expiryNotice: "Acest cod expiră in 7 zile.",
                        errorMaxCodes:
                                "Ai atins numărul maxim de coduri disponibile. Şterge din codurile existente sau aşteaptă să fie folosite.",
                },
                codesModal: {
                        title: "Coduri de invitație",
                        subtitle: "Acestea sunt codurile tale de invitație active.",
                        createdAt: "Creat pe data ",
                        expiresAt: "Expiră pe data ",
                },
                deactivateClient: {
                        title: "Dezactivează contul solicitantului",
                        subtitle1: (text: string) =>
                                `Eşti sigur că vrei să dezactivezi contul solicitantului ${text}?`,
                        subtitle2: "Această acțiune îi va interzice accesul la cont. Acțiunea este reversibilă.",
                        confirm: "Dezactivează",
                },
        },
        buttons: {
                sendNotification: {
                        normal: "Notifică solicitantul",
                },
                uploadDocument: {
                        normal: "Încarcă",
                        inProgress: "Se încarcă",
                },
                viewFile: {
                        normal: "Deschide",
                        inProgress: "Se încarcă",
                },
                addClient: {
                        normal: "Solicitant nou",
                },
                viewInvitationCodes: {
                        normal: "Vezi coduri de invitație",
                },
        },
        sidebar: {
                overview: "Dosare",
                clients: "Solicitanți",
                settings: "Setări",
                logout: "Deconectare",
        },
};
