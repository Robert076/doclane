# Doclane

## 🚀 Table of contents

- [AWS Deployment](#️-aws-deployment)

**Doclane** is a document request and workflow platform designed to help
professionals collect documents from their clients in a clear,
structured, and reliable way.

It enables teams to request documents (one-time or recurring), while
clients upload exactly what is needed, on time, through a simple and
secure portal.

Doclane focuses on eliminating document chaos, reducing follow-ups, and
creating a smooth experience for both professionals and their clients.

---

## 🎯 Problem

Professionals and their clients often exchange documents via email,
messaging apps, or shared folders.

This results in:

- lost or duplicated files
- constant follow-ups and reminders
- lack of visibility into what is missing or overdue
- wasted time and unnecessary friction

---

## ✅ Solution

Doclane provides a single, structured workspace where:

- Professionals create clients and define required documents
- Document requests can be one-time or recurring (e.g. monthly)
- Clients upload files directly to each request
- Both sides see clear request status: missing, uploaded, overdue

Doclane is **not** a document storage service.
It is a document request and workflow platform.

---

## ☁️ AWS Deployment

The deployment is done in AWS, and is made up of three distinct environments:

- dev
- staging
- prod

For better maintainability, scalability and the ability to audit changes to the infrastructure in a version controlled way, I chose to use Terraform for the deployment.

### 🧑🏼‍💻 Dev environment

The dev environment is used for local testing. It contains only a subset of the services used in prod or staging.

Since when developing locally we run the frontend and backend on our machine, the only provisioned resources in the AWS cloud are an S3 bucket used for storing documents, a dev IAM role containing policies that adhere to the least-privilege principle, only being able to access the S3 bucket in the dev environment (the role is assumed by the local backend).

> Check out the Terraform configuration files for the dev environment in /terraform/dev

## 📄 License

This project is licensed under a proprietary license.
See the `LICENSE` file for details.
