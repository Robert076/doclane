-- Extensions
CREATE EXTENSION IF NOT EXISTS citext;

-- Departments
CREATE TABLE departments (
    id          SERIAL PRIMARY KEY,
    name        CITEXT        NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Users
CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    email         CITEXT        NOT NULL UNIQUE,
    first_name    TEXT          NOT NULL,
    last_name     TEXT          NOT NULL,
    password_hash TEXT          NOT NULL,
    role          TEXT          NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'member', 'citizen')),
    department_id INTEGER       REFERENCES departments(id) ON DELETE SET NULL,
    is_active     BOOLEAN       NOT NULL DEFAULT TRUE,
    phone         TEXT,
    street        TEXT,
    locality      TEXT,
    last_notified TIMESTAMPTZ,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Invitation codes
CREATE TABLE invitation_codes (
    id            SERIAL PRIMARY KEY,
    code          TEXT          NOT NULL UNIQUE,
    created_by    INTEGER       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    department_id INTEGER       NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    used_at       TIMESTAMPTZ,
    expires_at    TIMESTAMPTZ,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Request templates
CREATE TABLE request_templates (
    id              SERIAL PRIMARY KEY,
    title           TEXT        NOT NULL,
    description     TEXT,
    department_id   INTEGER     NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    is_recurring    BOOLEAN     NOT NULL DEFAULT FALSE,
    recurrence_cron TEXT,
    is_closed       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_by      INTEGER     NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Expected document templates
CREATE TABLE expected_document_templates (
    id                          SERIAL PRIMARY KEY,
    document_request_template_id INTEGER NOT NULL REFERENCES request_templates(id) ON DELETE CASCADE,
    title                       TEXT    NOT NULL,
    description                 TEXT    NOT NULL DEFAULT '',
    example_file_path           TEXT,
    example_mime_type           TEXT
);

-- Document requests
CREATE TABLE document_requests (
    id               SERIAL PRIMARY KEY,
    title            TEXT        NOT NULL,
    description      TEXT,
    assignee         INTEGER     NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    department_id    INTEGER     NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    template_id      INTEGER     REFERENCES request_templates(id) ON DELETE SET NULL,
    is_recurring     BOOLEAN     NOT NULL DEFAULT FALSE,
    recurrence_cron  TEXT,
    is_scheduled     BOOLEAN     NOT NULL DEFAULT FALSE,
    scheduled_for    TIMESTAMPTZ,
    is_cancelled     BOOLEAN     NOT NULL DEFAULT FALSE,
    is_closed        BOOLEAN     NOT NULL DEFAULT FALSE,
    last_uploaded_at TIMESTAMPTZ,
    next_due_at      TIMESTAMPTZ,
    due_date         TIMESTAMPTZ,
    claimed_by       INTEGER     REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Expected documents (per request)
CREATE TABLE expected_documents (
    id                      SERIAL PRIMARY KEY,
    document_request_id     INTEGER NOT NULL REFERENCES document_requests(id) ON DELETE CASCADE,
    title                   TEXT    NOT NULL,
    description             TEXT    NOT NULL DEFAULT '',
    status                  TEXT    NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'uploaded', 'approved', 'rejected')),
    rejection_reason        TEXT,
    example_file_path       TEXT,
    example_mime_type       TEXT
);

-- Documents (uploaded files)
CREATE TABLE documents (
    id                      SERIAL PRIMARY KEY,
    document_request_id     INTEGER     NOT NULL REFERENCES document_requests(id) ON DELETE CASCADE,
    expected_document_id    INTEGER     NOT NULL REFERENCES expected_documents(id) ON DELETE CASCADE,
    file_name               TEXT        NOT NULL,
    file_path               TEXT        NOT NULL,
    mime_type               TEXT,
    file_size               BIGINT,
    s3_version_id           TEXT,
    uploaded_by             INTEGER     REFERENCES users(id) ON DELETE SET NULL,
    uploaded_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Request comments
CREATE TABLE request_comments (
    id          SERIAL PRIMARY KEY,
    request_id  INTEGER     NOT NULL REFERENCES document_requests(id) ON DELETE CASCADE,
    user_id     INTEGER     NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    comment     TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    color TEXT NOT NULL DEFAULT '#6366f1',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE template_tags (
    template_id INT NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (template_id, tag_id)
);

-- Indexes
CREATE INDEX idx_users_department_id          ON users(department_id);
CREATE INDEX idx_users_role                   ON users(role);
CREATE INDEX idx_invitation_codes_department  ON invitation_codes(department_id);
CREATE INDEX idx_invitation_codes_created_by  ON invitation_codes(created_by);
CREATE INDEX idx_request_templates_dept       ON request_templates(department_id);
CREATE INDEX idx_document_requests_assignee   ON document_requests(assignee);
CREATE INDEX idx_document_requests_dept       ON document_requests(department_id);
CREATE INDEX idx_document_requests_claimed_by ON document_requests(claimed_by);
CREATE INDEX idx_document_requests_template   ON document_requests(template_id);
CREATE INDEX idx_expected_documents_request   ON expected_documents(document_request_id);
CREATE INDEX idx_documents_request            ON documents(document_request_id);
CREATE INDEX idx_documents_expected           ON documents(expected_document_id);
CREATE INDEX idx_documents_uploaded_by        ON documents(uploaded_by);
CREATE INDEX idx_request_comments_request     ON request_comments(request_id);
CREATE INDEX idx_request_comments_user        ON request_comments(user_id);