/**
 * Maps backend error/status messages (English) to Romanian.
 *
 * The backend returns all messages in English and exposes no stable error code,
 * so translation is done by matching the exact message string. Any message not
 * present in the dictionary falls through unchanged, so untranslated strings
 * degrade gracefully instead of breaking.
 *
 * Keep keys in sync with the `Msg:` strings in the Go backend.
 */
const ERROR_TRANSLATIONS: Record<string, string> = {
	// Auth / access
	"Unauthorized.": "Nu ești autentificat.",
	"Forbidden.": "Acces interzis.",
	"Access denied.": "Acces interzis.",
	"Your account is deactivated.": "Contul tău este dezactivat.",
	"Invalid token.": "Token invalid.",
	"Invalid token claims.": "Token invalid.",

	// Generic request validation
	"Invalid request body.": "Date trimise invalide.",
	"Invalid JSON body.": "Date trimise invalide.",
	"The body received is invalid.": "Datele trimise sunt invalide.",
	"Invalid ID received.": "ID invalid.",
	"Failed to parse form.": "Nu s-a putut procesa formularul.",

	// IDs
	"Invalid request ID.": "ID-ul cererii este invalid.",
	"Invalid request ID format.": "Format invalid pentru ID-ul cererii.",
	"Invalid document request ID format.": "Format invalid pentru ID-ul cererii.",
	"Invalid assignee ID format.": "Format invalid pentru ID-ul responsabilului.",
	"Invalid code ID format.": "Format invalid pentru ID-ul codului.",
	"Invalid department ID format.": "Format invalid pentru ID-ul departamentului.",
	"Invalid department id.": "ID de departament invalid.",
	"Invalid department_id value.": "Valoare invalidă pentru departament.",
	"department_id is required.": "Departamentul este obligatoriu.",
	"Missing department_id query parameter.": "Lipsește parametrul departament.",
	"Invalid expected document ID format.":
		"Format invalid pentru ID-ul documentului.",
	"Invalid expected document template ID format.":
		"Format invalid pentru ID-ul șablonului de document.",
	"Invalid file ID.": "ID de fișier invalid.",
	"Invalid file ID format.": "Format invalid pentru ID-ul fișierului.",
	"Invalid tag ID.": "ID de etichetă invalid.",
	"Invalid template ID.": "ID de șablon invalid.",
	"Invalid template ID format.": "Format invalid pentru ID-ul șablonului.",
	"Invalid user ID.": "ID de utilizator invalid.",
	"Invalid limit value.": "Valoare invalidă pentru limită.",
	"Invalid offset value.": "Valoare invalidă pentru offset.",
	"Invalid or missing expected_document_id.":
		"ID-ul documentului așteptat lipsește sau este invalid.",

	// Requests
	"Request not found.": "Cererea nu a fost găsită.",
	"You are not allowed to create requests": "Nu ai voie să creezi cereri.",
	"You are not allowed to view these requests.":
		"Nu ai voie să vizualizezi aceste cereri.",
	"You don't have access to this request.": "Nu ai acces la această cerere.",
	"You don't have access to edit the request.":
		"Nu ai acces să editezi cererea.",
	"You are not allowed to cancel this request.":
		"Nu ai voie să anulezi această cerere.",
	"Only admins can view all requests.":
		"Doar administratorii pot vedea toate cererile.",
	"Only admins and department members can view archived requests.":
		"Doar administratorii și membrii departamentului pot vedea cererile arhivate.",
	"Only admins and department members can view cancelled requests.":
		"Doar administratorii și membrii departamentului pot vedea cererile anulate.",
	"A request cannot be closed if the status is not 'pending'.":
		"Cererea nu poate fi închisă dacă statusul nu este „în așteptare”.",
	"A request marked as scheduled must have a scheduled_for field.":
		"O cerere programată trebuie să aibă o dată de programare.",
	"Due date cannot be in the past.":
		"Termenul limită nu poate fi în trecut.",

	// Claiming
	"Failed to claim request.": "Nu s-a putut prelua cererea.",
	"Failed to unclaim request.": "Nu s-a putut renunța la cerere.",
	"Cannot claim a closed or cancelled request.":
		"Nu poți prelua o cerere închisă sau anulată.",
	"Cannot unclaim an archived request.":
		"Nu poți renunța la o cerere arhivată.",
	"Only department members can claim requests.":
		"Doar membrii departamentului pot prelua cereri.",
	"You can only claim requests from your department.":
		"Poți prelua doar cereri din departamentul tău.",
	"You can only unclaim requests you have claimed.":
		"Poți renunța doar la cererile pe care le-ai preluat.",
	"You have already claimed this request.": "Ai preluat deja această cerere.",
	"This request has already been claimed by another member.":
		"Această cerere a fost deja preluată de alt membru.",
	"This request is not claimed.": "Această cerere nu este preluată.",
	"You must claim the request before working on it.":
		"Trebuie să preiei cererea înainte de a lucra la ea.",
	"You must claim the request before closing it.":
		"Trebuie să preiei cererea înainte de a o închide.",
	"You can only work on requests that are claimed by you.":
		"Poți lucra doar la cereri pe care le-ai preluat.",
	"Only the member who claimed this request can close it.":
		"Doar membrul care a preluat cererea o poate închide.",
	"Only the member who claimed this request can reopen it.":
		"Doar membrul care a preluat cererea o poate redeschide.",

	// Documents
	"Expected document not found.": "Documentul nu a fost găsit.",
	"Must provide a reason for rejecting the document.":
		"Trebuie să oferi un motiv pentru refuzul documentului.",
	"Only the department handling this request can update document status.":
		"Doar departamentul care gestionează cererea poate actualiza statusul documentului.",
	"This document has no example file.": "Acest document nu are un fișier exemplu.",
	"This template document does not have an example.":
		"Acest document din șablon nu are un exemplu.",
	"Only staff can analyze documents.":
		"Doar personalul poate analiza documente.",
	"Only staff can interpret documents.":
		"Doar personalul poate interpreta documente.",
	"Only staff can use this feature.":
		"Doar personalul poate folosi această funcție.",

	// Files / uploads
	"File is empty.": "Fișierul este gol.",
	"File size must be less than 20MB.":
		"Dimensiunea fișierului trebuie să fie sub 20MB.",
	"Could not get file from request.": "Nu s-a putut citi fișierul din cerere.",
	"Failed to upload to S3.": "Încărcarea fișierului a eșuat.",
	"You are not allowed to view this file.":
		"Nu ai voie să vizualizezi acest fișier.",

	// Templates
	"Template not found.": "Șablonul nu a fost găsit.",
	"This template is archived.": "Acest șablon este arhivat.",
	"You are not allowed to create templates.":
		"Nu ai voie să creezi șabloane.",
	"You don't have access to this template.": "Nu ai acces la acest șablon.",
	"Title is required.": "Titlul este obligatoriu.",
	"Title cannot be empty.": "Titlul nu poate fi gol.",
	"Title cannot exceed 255 characters.":
		"Titlul nu poate depăși 255 de caractere.",
	"Title must be between 3 and 30 characters.":
		"Titlul trebuie să aibă între 3 și 30 de caractere.",
	"New title is too short or too long. Minimum 3 characters, maximum 30 characters.":
		"Titlul este prea scurt sau prea lung. Minim 3, maxim 30 de caractere.",
	"Description cannot exceed 1000 characters.":
		"Descrierea nu poate depăși 1000 de caractere.",
	"A template marked as recurring must have a recurrence_cron field.":
		"Un șablon recurent trebuie să aibă un program de recurență.",
	"Recurrence cron is required when is_recurring is true.":
		"Programul de recurență este obligatoriu pentru șabloane recurente.",
	"Invalid cron expression.": "Expresie cron invalidă.",
	"Invalid recurrence_cron format.": "Format invalid pentru recurență.",
	"A template must be provided.": "Trebuie selectat un șablon.",

	// Tags
	"Only admins can manage tags.": "Doar administratorii pot gestiona etichete.",
	"Maximum tag count has been reached.":
		"Numărul maxim de etichete a fost atins.",
	"Tag name is required.": "Numele etichetei este obligatoriu.",
	"Tag color must be a valid hex color (e.g. #ff5722).":
		"Culoarea etichetei trebuie să fie un cod hex valid (ex. #ff5722).",

	// Comments
	"Cannot add comment to closed or cancelled request.":
		"Nu poți adăuga comentarii la o cerere închisă sau anulată.",
	"Comment is too long (max 200 characters).":
		"Comentariul este prea lung (maxim 200 de caractere).",
	"Comment must contain at least 3 visible characters.":
		"Comentariul trebuie să conțină cel puțin 3 caractere vizibile.",
	"Please wait before posting another comment.":
		"Așteaptă puțin înainte de a posta un alt comentariu.",

	// Users / departments
	"User not found.": "Utilizatorul nu a fost găsit.",
	"User already exists.": "Utilizatorul există deja.",
	"First and last name cannot be empty.":
		"Numele și prenumele nu pot fi goale.",
	"Failed to update profile.": "Actualizarea profilului a eșuat.",
	"Only admins can list users by department.":
		"Doar administratorii pot lista utilizatorii pe departament.",
	"Only admins can move users between departments.":
		"Doar administratorii pot muta utilizatori între departamente.",
	"Cannot change department for a user who is not a department member.":
		"Nu poți schimba departamentul unui utilizator care nu este membru.",
	"User must unclaim all requests before being moved to another department.":
		"Utilizatorul trebuie să renunțe la toate cererile înainte de a fi mutat în alt departament.",
	"You do not have permission to deactivate this account.":
		"Nu ai permisiunea să dezactivezi acest cont.",
	"You must update your phone number first.":
		"Trebuie să îți actualizezi mai întâi numărul de telefon.",
	"You must update the locality where you live first.":
		"Trebuie să îți actualizezi mai întâi localitatea.",
	"You must update the street where you live first.":
		"Trebuie să îți actualizezi mai întâi strada.",

	// Departments
	"Department not found.": "Departamentul nu a fost găsit.",
	"Department name cannot be empty.":
		"Numele departamentului nu poate fi gol.",
	"Only admins can create departments.":
		"Doar administratorii pot crea departamente.",

	// Invitation codes
	"Invalid invitation code.": "Cod de invitație invalid.",
	"Invitation code not found.": "Codul de invitație nu a fost găsit.",
	"This invitation code has expired.": "Codul de invitație a expirat.",
	"This invitation code has already been used.":
		"Codul de invitație a fost deja folosit.",
	"Code is required.": "Codul este obligatoriu.",
	"Email is required.": "Adresa de email este obligatorie.",
	"Only 3 active codes are allowed at one time.":
		"Sunt permise maxim 3 coduri active simultan.",
	"Only admins can generate invitation codes.":
		"Doar administratorii pot genera coduri de invitație.",
	"Only admins can delete invitation codes.":
		"Doar administratorii pot șterge coduri de invitație.",
	"Only admins can view invitation codes.":
		"Doar administratorii pot vedea codurile de invitație.",
	"You can only access your own invitation codes.":
		"Poți accesa doar propriile coduri de invitație.",
	"Failed to generate invitation code.":
		"Generarea codului de invitație a eșuat.",
	"Failed to save invitation code.":
		"Salvarea codului de invitație a eșuat.",
	"Failed to delete invitation code.":
		"Ștergerea codului de invitație a eșuat.",

	// Stats
	"Only admins can view stats.":
		"Doar administratorii pot vedea statisticile.",
};

/**
 * Translates a backend message to Romanian.
 *
 * Tries an exact match first, then a few prefix matches for templated
 * (format-string) backend messages. Returns the original message if no
 * translation is found.
 */
export function translateError(message?: string | null): string {
	if (!message) return "";

	const exact = ERROR_TRANSLATIONS[message];
	if (exact) return exact;

	// Templated backend messages (contain runtime values via fmt.Sprintf).
	if (message.startsWith("Invalid role:")) return "Rol invalid.";
	if (message.endsWith("is not allowed.") && message.startsWith("Extension"))
		return "Tipul de fișier nu este permis.";
	if (message.startsWith("Templates can have up to"))
		return "Numărul maxim de etichete pentru un șablon a fost depășit.";
	if (message.startsWith("Could not fetch users"))
		return "Nu s-au putut încărca utilizatorii.";
	if (message.startsWith("Failed to upload"))
		return "Încărcarea fișierului a eșuat.";

	return message;
}
