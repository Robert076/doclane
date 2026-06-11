# Doclane
<img width="900" height="508" alt="Screenshot 2026-05-30 at 12 21 37" src="https://github.com/user-attachments/assets/39ce1ae2-8d17-4b74-9a4c-c8e3f0cb1890" />

<img width="900" height="508" alt="Screenshot 2026-05-30 at 12 30 00" src="https://github.com/user-attachments/assets/53a88bb2-f4a2-447b-a20c-fe3042838c21" />

<img width="900" height="508" alt="Screenshot 2026-05-30 at 12 41 06" src="https://github.com/user-attachments/assets/df176c18-5c5b-4e4e-a76a-36e4c53728e9" />

Doclane is a web-based application that supports document workflows between citizens and public institutions. It moves the document exchange entirely online, structures it clearly for both sides, and keeps the citizen informed throughout, replacing the repeated in-person visits that paper-based processes typically require.

🔗 Live demo: [thesis.robert-beres.com](https://thesis.robert-beres.com)

## Overview

In many public institutions, document requests still rely on physical visits, printed files, and manual processing. Citizens travel to an office, wait in queues, and hand over paperwork just to initiate routine requests, often only to discover later that a document is missing or incorrectly formatted, after which the whole cycle repeats.

Most existing digital tools cover only part of the workflow, handling submission but not validation, or validation but not feedback, leaving gaps that still require in-person follow-up. Doclane closes that loop: users submit requests online, upload required documents, and track status, while institution employees review submissions, approve or reject individual documents, and send feedback directly through the system.

Doclane is a workflow support tool, not a replacement for institutional authority. Final decisions about whether a request or document is valid remain with authorised employees. It does not currently replace the physical pickup of final documents where that is legally required, but it removes most of the unnecessary visits that happen before that step.

## Main Features

### Citizen side
- Registration and login (handled by AWS Cognito)
- Profile management (contact details like address and phone number, attached automatically to requests)
- Browse request templates and submit requests
- Upload required documents into designated slots
- Track request status from a dashboard
- Receive per-document rejection reasons and resubmit corrected files through the same interface

### Department employee side
- See incoming requests routed to their department
- Claim a request so one employee is clearly responsible (avoids duplicate effort)
- Review each uploaded document independently
- Approve or reject individual documents, with a required reason on rejection
- Request is closed once every document is approved

### Administrator side
- Create and manage request templates
- Define required document slots and information fields per template
- Attach example files per slot as a reference for citizens
- Tag templates to group related request types

### AI-assisted document processing
The system integrates three fully managed, serverless AWS services. No model training or dataset curation is involved; the backend calls these APIs and surfaces the results to the reviewing employee. None of them make or record decisions on their own.
- **Textract** — OCR / text extraction from uploaded scans and photos
- **Bedrock (Amazon Nova Lite)** — interprets whether an uploaded document matches what the template requires, as advisory context only (used as-is, no fine-tuning)
- **Polly** — text-to-speech synthesis of document content for accessibility and quick review

## Workflow

1. User submits a request from a template and fills in required fields
2. User uploads the required documents
3. The request is routed to the responsible department
4. A department employee claims the request
5. Each document is reviewed and approved or rejected with feedback
6. The user uploads corrected files if needed; the request re-enters the review queue
7. Once all documents are approved, the request is marked complete

## Technology Stack

- **Backend:** Go (Chi router), layered into handlers / services / repositories / models
- **Frontend:** Next.js (React, server-side rendering)
- **Database:** PostgreSQL on AWS RDS
- **File storage:** AWS S3 (versioned, server-side encrypted, presigned URL access)
- **Authentication:** AWS Cognito, JWT stored in HTTP-only cookies
- **Document processing:** AWS Textract, Bedrock (Nova Lite), Polly
- **Cloud provider:** AWS

## Architecture & Deployment

Doclane is fully cloud-native and runs on AWS, orchestrated by Kubernetes and provisioned entirely as code.

- **Compute:** EKS (Elastic Kubernetes Service) cluster inside a VPC spanning two availability zones in `eu-west-1`
- **Containers:** Backend and frontend packaged via multi-stage Docker builds (backend ~36 MB on distroless, running as non-root; frontend ~370 MB on Node.js Alpine)
- **Scaling & health:** Backend runs 2–6 replicas behind a HorizontalPodAutoscaler (CPU-based), with liveness/readiness probes; frontend runs 2–4 replicas
- **Networking:** Worker nodes and RDS sit in private subnets; an Application Load Balancer (ALB) and NAT Gateway sit in public subnets. The ALB routes `/api` to the backend and everything else to the frontend, terminates TLS via an ACM certificate, and redirects HTTP to HTTPS. Route53 points the domain to the ALB.
- **Security:** RDS reachable only from the EKS cluster security group; all DB traffic over SSL/TLS. Backend uses **IRSA** (IAM Roles for Service Accounts) for short-lived, least-privilege credentials instead of static keys. S3 access via time-limited presigned URLs.
- **Infrastructure as Code:** Terraform, split into four stacks (`bootstrap`, `platform-data`, `platform-compute`, `workload`) with remote state in S3 and DynamoDB-based locking
- **Ephemeral environments:** `make up` provisions the full environment (~15 min); `make down` tears it down (~10 min) to bring running cost to zero
- **CI/CD:** GitHub Actions builds and pushes backend/frontend images to ECR (tagged by commit SHA), authenticating to AWS via GitHub OIDC rather than stored credentials

## Key Design Decisions

- Clear separation of backend layers (handlers, services, repositories, models)
- Role-based access control (citizen, department employee, administrator)
- Uniform API response structure (`success`, `message`, `error`, `data`) with typed errors
- Storage and authentication kept behind abstractions to reduce provider lock-in
- AI features are strictly advisory and never sit on the critical path of a decision
- Production-grade infrastructure designed to be reproducible and torn down on demand

## Limitations

- No integration with national identity systems; submissions rely on self-registered Cognito accounts with no verified identity link
- Several features depend on AWS managed services; outages of S3 or Cognito would impact file access and login respectively
- AI outputs are constrained by OCR quality and the general-purpose model's accuracy, which is why they remain advisory
- Adoption depends on organisational change and digital literacy, not just the software

## Future Improvements

- **eIDAS-compliant digital signing** — legally recognised remote identity verification, removing the last reason a citizen might need to appear in person and increasing the legal weight of submissions
- **In-person scheduling** — let employees propose time slots once a request is complete so citizens can book physical document collection through the platform
- Fallback strategy for authentication (e.g. short-window token validation caching) for production resilience
- Audit logging and document retention policies for institutional-scale deployment

## Goal

The main goal of Doclane is to simplify how citizens interact with public institutions by reducing physical paperwork and making document workflows clearer, faster, and more structured, while keeping legal authority and final decisions firmly with institution staff.
