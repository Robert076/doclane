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

The deployment is done in AWS, and it is mostly serverless except the database layer. It can handle spikes in load without issues.

## 🛥️ CD to AWS

Code in AWS is never manually updated. Instead, it is updated automatically on pushes to the main branch using Github Actions.

### How the backend is updated

The backend is updated in the least-privilege principle, and no permanent credentials ever exist.

The Github action generates an id token using Github's OIDC, and that signed token (Github signs it with their private key) is then sent to AWS to prove the identity of the runner.
AWS checks that the signature cryptographically matches by using Github's public key, and confirms that the runner is actually who it pretends to be.
In AWS, a provider has been created in IAM that has the provider as Github, and type of OIDC. The audience of this provider is sts.amazonaws.com and that means AWS' STS service can now federate an identity and return temporary IAM credentials, which in turn the runner uses to update the Lambda.

By doing this, we safely update code on every push to main branch.

-

## 📄 License

This project is licensed under a proprietary license.
See the `LICENSE` file for details.
