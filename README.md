# Doclane

Doclane is a web-based system for handling document workflows between citizens and public institutions.  
It is designed to replace manual paper-based processes with a structured digital workflow.

## Overview

In many public institutions, document requests still require physical visits, printed files, and manual processing.  
Doclane moves this process into a digital system where users can submit requests online, upload documents, and track progress.

On the other side, institution employees can review submissions, approve or reject documents, and communicate feedback directly in the system.

## Main Features

### User side
- User registration and login
- Profile management (personal data like address, phone number)
- Submit document requests based on templates
- Upload required files
- Track request status
- Receive feedback from institutions

### Institution side
- Review submitted requests
- Approve or reject individual documents
- Add rejection reasons and feedback
- Claim requests (one responsible employee per request)
- Final approval and completion of requests

### Admin side
- Create request templates
- Add required document fields for each template
- Provide example files for guidance
- Tag templates for organization
- Manage structure of available request types

### AI assistance
Doclane integrates AI-based tools to help with document processing:
- Image-to-text extraction (OCR) for uploaded files
- Document quality checks (e.g. blurry or unreadable scans)
- Content interpretation to verify if document matches requirements
- Text-to-speech (TTS) for accessibility in some cases

## System idea

The system is built around a simple workflow:

1. User submits a request
2. User uploads required documents
3. Department employee claims the request
4. Employee reviews documents
5. Documents are approved or rejected with feedback
6. User fixes issues if needed and resubmits
7. Final approval is given when everything is correct

## Technologies

- Cloud provider: AWS
- Backend: Go (Golang)
- Frontend: Next.js
- Database: PostgreSQL
- File storage: AWS S3
- Authentication: JWT (HTTP-only cookies)

## Key Design Decisions

- Separation of backend layers (handlers, services, repositories)
- Role-based access control (user, employee, admin)
- Stateless authentication using JWT
- Cloud-based file storage using S3
- Clear workflow-based request system instead of ad-hoc document handling

## Future Improvements

- Two-factor authentication (2FA)
- eIDAS-compliant digital signatures
- Fully digital identity verification (remove physical visits completely)
- Request timeline view (full history of actions per request)
- Better audit logging for administrative actions

## Goal

The main goal of Doclane is to simplify how citizens interact with public institutions by reducing physical paperwork and making document workflows clearer, faster, and more structured.
