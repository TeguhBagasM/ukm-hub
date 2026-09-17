# Product Requirements Document (PRD)

## UKM Management Platform

**Version:** 1.0.0
**Status:** MVP
**Last Updated:** 17 September 2026

---

## 1. Product Overview

### 1.1 Product Name

**UKM Management Platform**

### 1.2 Product Description

UKM Management Platform adalah aplikasi web untuk membantu organisasi mahasiswa seperti Unit Kegiatan Mahasiswa (UKM) mengelola organisasi, divisi, kegiatan, pendaftaran anggota, proses seleksi, dan data anggota secara terpusat.

Sistem dirancang agar **generic dan configurable**, sehingga setiap UKM dapat memiliki struktur divisi, kegiatan, serta formulir pendaftaran yang berbeda tanpa perlu mengubah source code aplikasi.

Sebagai initial seed/demo organization, aplikasi menggunakan:

> **Creative Student Association (CSA)**

dengan dua divisi:

* Programming
* Multimedia

### 1.3 Core Value

Fitur utama yang menjadi pembeda aplikasi adalah **Dynamic Registration Form Builder**.

Admin organisasi dapat membuat formulir pendaftaran sendiri dengan memilih dan mengatur field secara visual.

Contoh:

```text
Event:
Open Recruitment CSA 2026

Form:
- Nama Lengkap       [Text]
- NIM                [Text]
- Email              [Email]
- Nomor WhatsApp     [Phone]
- Divisi Pilihan     [Select]
- Alasan Bergabung   [Textarea]
- Portfolio          [URL]
```

Admin tidak perlu mengubah source code untuk membuat formulir dengan struktur berbeda.

---

# 2. Product Goals

## 2.1 Primary Goals

1. Menyediakan sistem manajemen UKM yang generic.
2. Memungkinkan organisasi memiliki banyak divisi.
3. Memungkinkan admin membuat event/kegiatan.
4. Memungkinkan admin membuat registration form secara dinamis.
5. Memungkinkan mahasiswa mendaftar melalui public event page.
6. Memungkinkan admin melakukan review terhadap pendaftar.
7. Memungkinkan pendaftar yang diterima dikonversi menjadi anggota.
8. Menyediakan dashboard dengan data aktual.
9. Menyediakan export data dalam format CSV.
10. Memisahkan frontend dan backend menggunakan REST API.

## 2.2 Portfolio Goals

Aplikasi juga dirancang sebagai portfolio project untuk menunjukkan kemampuan:

* Golang backend development
* REST API design
* Clean Architecture
* Authentication & authorization
* Database design
* Dynamic data modeling
* React frontend development
* TypeScript
* State management
* Form handling
* API integration
* Responsive UI
* Data validation
* CSV export

---

# 3. Target Users

## 3.1 Super Admin

Super Admin mengelola organisasi yang tersedia dalam platform.

Responsibilities:

* Manage organizations
* Manage organization administrators
* View organization information
* Manage platform-level data

## 3.2 Organization Admin

Admin yang bertanggung jawab terhadap satu organisasi/UKM.

Responsibilities:

* Manage organization profile
* Manage divisions
* Manage events
* Build registration forms
* Review registrations
* Accept/reject applicants
* Convert accepted applicants into members
* Manage members
* Export data

## 3.3 Public Applicant

Mahasiswa atau pengguna umum yang ingin mendaftar ke suatu kegiatan/rekrutmen.

Responsibilities:

* View published events
* Open registration form
* Fill registration form
* Submit registration
* View submission result

---

# 4. MVP Scope

MVP harus memprioritaskan **functional end-to-end flow** daripada jumlah fitur.

Urutan implementasi:

```text
1. Data Layer
2. Authentication
3. Organization Management
4. Division Management
5. Event Management
6. Dynamic Form Builder
7. Public Registration
8. Registration Management
9. Member Management
10. Dashboard
11. CSV Export
12. UI Polish
```

---

# 5. Core User Flow

Core flow yang wajib bekerja:

```text
Login
  ↓
Organization Dashboard
  ↓
Organization
  ↓
Create/Edit Division
  ↓
Create Event
  ↓
Create Registration Form
  ↓
Form Builder
  ↓
Add Dynamic Fields
  ↓
Save & Publish
  ↓
Public Event Page
  ↓
Applicant submits registration
  ↓
Registration stored
  ↓
Admin opens Applications
  ↓
View Applicant
  ↓
Accept / Reject
  ↓
Accept Applicant
  ↓
Convert to Member
  ↓
Member Management
```

Jika flow ini belum berjalan end-to-end, MVP dianggap belum selesai.

---

# 6. Functional Requirements

## 6.1 Authentication

### Requirements

* User dapat login.
* Sistem melakukan authentication menggunakan backend.
* Password tidak disimpan dalam bentuk plaintext.
* Backend menghasilkan authentication token.
* Frontend menyimpan session/token dengan aman sesuai implementasi.
* Endpoint yang membutuhkan authentication harus dilindungi middleware.

### Roles

```text
SUPER_ADMIN
ORG_ADMIN
```

### Authorization

Organization Admin hanya boleh mengakses data organisasi yang menjadi tanggung jawabnya.

Contoh:

```text
Admin CSA
    ↓
CSA data       ✓
CSA divisions  ✓
CSA events     ✓
CSA members    ✓

Organization B data
    ↓
                 ✗
```

---

# 7. Organization Management

## 7.1 Organization

Organization merepresentasikan UKM atau organisasi yang menggunakan platform.

### Organization Fields

```text
id
name
slug
description
logo
email
phone
status
created_at
updated_at
```

### Features

Super Admin:

* View organizations
* Create organization
* Edit organization
* Delete organization
* Activate/deactivate organization

Organization Admin:

* View organization
* Edit organization profile

---

# 8. Division Management

Setiap organization dapat memiliki jumlah division yang berbeda.

Contoh CSA:

```text
Creative Student Association
├── Programming
└── Multimedia
```

Organization lain dapat memiliki:

```text
Organization A
├── Human Resources
├── Public Relations
├── Event
└── Research
```

### Division Fields

```text
id
organization_id
name
description
status
created_at
updated_at
```

### Requirements

Admin dapat:

* Create division
* Edit division
* Delete division
* Activate/deactivate division
* View division details

Division tidak boleh hardcoded di frontend.

---

# 9. Event Management

Event merupakan kegiatan atau program yang memiliki registration.

Contoh:

```text
Open Recruitment CSA 2026
```

### Event Fields

```text
id
organization_id
name
slug
description
location
start_date
end_date
registration_start
registration_end
quota
status
created_at
updated_at
```

### Event Status

```text
DRAFT
PUBLISHED
CLOSED
ARCHIVED
```

### Requirements

Admin dapat:

* Create event
* Edit event
* Delete event
* Publish event
* Close registration
* Archive event
* Set registration deadline
* Set quota
* View event detail
* View registration count

### Business Rules

Event hanya dapat diakses melalui public registration apabila:

```text
status = PUBLISHED
```

dan periode pendaftaran masih aktif.

---

# 10. Dynamic Registration Form Builder

## 10.1 Overview

Dynamic Form Builder adalah fitur utama aplikasi.

Admin dapat membuat formulir pendaftaran untuk setiap event tanpa mengubah source code.

### Supported Field Types

MVP wajib mendukung:

```text
TEXT
TEXTAREA
EMAIL
NUMBER
PHONE
DATE
SELECT
RADIO
CHECKBOX
URL
```

### Form Field Structure

```text
id
form_id
label
name
type
placeholder
description
required
options
sort_order
created_at
updated_at
```

### Form Builder Features

Admin dapat:

* Add field
* Edit field
* Delete field
* Duplicate field
* Reorder field
* Set field label
* Set field name
* Set placeholder
* Set description
* Set required/optional
* Configure options
* Preview form
* Save form
* Publish form

### Example

Admin membuat:

```text
Field 1
Type: TEXT
Label: Nama Lengkap
Required: true

Field 2
Type: EMAIL
Label: Email
Required: true

Field 3
Type: SELECT
Label: Pilihan Divisi
Options:
- Programming
- Multimedia
Required: true
```

Frontend kemudian merender form berdasarkan metadata tersebut.

Tidak boleh membuat form CSA secara hardcoded.

---

# 11. Public Registration

Public registration tidak membutuhkan login untuk applicant pada MVP.

### Flow

```text
GET /events/{slug}
       ↓
GET event form
       ↓
Render dynamic form
       ↓
Applicant fills form
       ↓
Client validation
       ↓
Server validation
       ↓
Submit
       ↓
Registration created
       ↓
Success page
```

### Registration Data

```text
id
event_id
status
submitted_at
created_at
updated_at
```

Status:

```text
PENDING
ACCEPTED
REJECTED
```

### Registration Answers

Jawaban dari dynamic form disimpan secara terpisah:

```text
id
registration_id
form_field_id
value
created_at
updated_at
```

Dengan pendekatan ini, setiap event dapat memiliki form dengan struktur berbeda.

---

# 12. Registration Management

Organization Admin dapat melihat seluruh pendaftar pada event.

### Features

* View registrations
* Search registrations
* Filter by status
* Filter by event
* View registration detail
* Accept applicant
* Reject applicant
* Add rejection reason

### Registration Detail

Detail applicant harus dirender berdasarkan field yang dibuat melalui Form Builder.

Contoh:

```text
Nama Lengkap
Teguh Bagas Mardiansyah

NIM
22110397

Email
example@email.com

Divisi
Programming

Alasan Bergabung
...
```

Frontend tidak boleh mengasumsikan field tertentu selalu tersedia.

---

# 13. Applicant Acceptance

Admin dapat menerima applicant.

Flow:

```text
PENDING
   ↓
ACCEPT
   ↓
ACCEPTED
```

Setelah accepted, admin dapat memilih:

```text
Convert to Member
```

---

# 14. Member Management

Applicant yang diterima dapat dikonversi menjadi member.

### Flow

```text
Registration
     ↓
Accepted
     ↓
Convert to Member
     ↓
Member
```

### Member Fields

```text
id
organization_id
division_id
registration_id
name
email
phone
student_id
status
joined_at
created_at
updated_at
```

### Features

Admin dapat:

* View members
* Search members
* Filter by division
* View member detail
* Edit member
* Change division
* Activate/deactivate member
* Delete member

---

# 15. Dashboard

Dashboard menampilkan data aktual dari backend.

### Metrics

Minimal:

```text
Total Members
Active Events
Total Registrations
Pending Applications
```

### Additional Information

Dashboard dapat menampilkan:

* Registration statistics
* Members by division
* Registration status distribution
* Recent registrations
* Active events

Chart harus menggunakan data aktual, bukan mock data setelah backend tersedia.

---

# 16. CSV Export

Admin dapat melakukan export data.

### Required Exports

1. Registrations
2. Members

### Registration Export

Minimal:

```text
Name
Email
Phone
Event
Status
Submitted At
```

Dynamic form answers juga harus dapat disertakan.

Contoh:

```text
Name
Email
Division
Reason
Portfolio
Status
Submitted At
```

Header CSV dapat berubah sesuai dynamic fields.

---

# 17. Data Model

Relasi utama:

```text
Organization
    │
    ├── Divisions
    │
    ├── Events
    │      │
    │      └── Form
    │            │
    │            └── Form Fields
    │
    ├── Members
    │
    └── Organization Admins


Event
  │
  └── Registrations
          │
          └── Registration Answers
                    │
                    └── Form Fields
```

### Main Entities

```text
users
organizations
organization_admins
divisions
events
forms
form_fields
registrations
registration_answers
members
```

---

# 18. Backend Requirements

## 18.1 Technology

Recommended stack:

```text
Golang
REST API
PostgreSQL
JWT Authentication
```

Recommended architecture:

```text
Clean Architecture
```

Suggested structure:

```text
backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── domain/
│   ├── repository/
│   ├── usecase/
│   ├── delivery/
│   │   └── http/
│   ├── middleware/
│   └── infrastructure/
│
├── migrations/
├── config/
├── pkg/
└── go.mod
```

The exact folder structure can be adjusted as long as separation of responsibility remains clear.

---

# 19. Backend API Requirements

API should follow REST conventions.

Example:

### Authentication

```http
POST /api/v1/auth/login
GET  /api/v1/auth/me
```

### Organizations

```http
GET    /api/v1/organizations
POST   /api/v1/organizations
GET    /api/v1/organizations/:id
PUT    /api/v1/organizations/:id
DELETE /api/v1/organizations/:id
```

### Divisions

```http
GET    /api/v1/organizations/:id/divisions
POST   /api/v1/organizations/:id/divisions
PUT    /api/v1/divisions/:id
DELETE /api/v1/divisions/:id
```

### Events

```http
GET    /api/v1/organizations/:id/events
POST   /api/v1/organizations/:id/events
GET    /api/v1/events/:id
PUT    /api/v1/events/:id
DELETE /api/v1/events/:id
POST   /api/v1/events/:id/publish
POST   /api/v1/events/:id/close
```

### Forms

```http
GET    /api/v1/events/:id/form
POST   /api/v1/events/:id/form
PUT    /api/v1/forms/:id
POST   /api/v1/forms/:id/publish
```

### Public Registration

```http
GET  /api/v1/public/events/:slug
GET  /api/v1/public/events/:slug/form
POST /api/v1/public/events/:slug/registrations
```

### Registrations

```http
GET  /api/v1/events/:id/registrations
GET  /api/v1/registrations/:id
POST /api/v1/registrations/:id/accept
POST /api/v1/registrations/:id/reject
```

### Members

```http
GET    /api/v1/organizations/:id/members
POST   /api/v1/registrations/:id/convert-member
GET    /api/v1/members/:id
PUT    /api/v1/members/:id
DELETE /api/v1/members/:id
```

---

# 20. API Design Rules

API harus:

* Menggunakan consistent response format.
* Menggunakan HTTP status code yang sesuai.
* Melakukan server-side validation.
* Mengembalikan error message yang jelas.
* Tidak mengembalikan password.
* Menggunakan pagination untuk list data yang berpotensi besar.
* Menggunakan authentication middleware untuk protected endpoints.
* Menggunakan authorization berdasarkan role dan organization ownership.
* Menghindari business logic di HTTP handler.
* Menggunakan transaction ketika beberapa database operation harus berhasil secara atomic.

Example response:

```json
{
  "success": true,
  "data": {},
  "message": "Event created successfully"
}
```

Error:

```json
{
  "success": false,
  "message": "Registration form is not available"
}
```

---

# 21. Frontend Requirements

## 21.1 Technology

Recommended:

```text
React
TypeScript
Vite
Tailwind CSS
React Router
TanStack Query
React Hook Form
Zod
Lucide React
```

Frontend bertanggung jawab terhadap:

* UI
* Client-side validation
* API communication
* State management
* Dynamic form rendering
* Loading/error states
* Responsive layout

Frontend tidak boleh menyimpan business logic utama yang seharusnya berada di backend.

---

# 22. Frontend Pages

## Public

```text
/
 /events
 /events/:slug
 /events/:slug/register
 /events/:slug/success
```

## Authentication

```text
/login
```

## Admin

```text
/dashboard
/organizations
/organizations/:id
/organizations/:id/divisions
/organizations/:id/events
/events/:id
/events/:id/form-builder
/events/:id/registrations
/registrations/:id
/members
/members/:id
```

---

# 23. Dynamic Form Rendering

Frontend harus memiliki reusable form renderer.

Concept:

```text
Form Schema
    ↓
Form Renderer
    ↓
Field Type
    ↓
React Component
```

Example:

```text
TEXT      → Input
TEXTAREA  → Textarea
EMAIL     → Email Input
NUMBER    → Number Input
SELECT    → Select
RADIO     → Radio Group
CHECKBOX  → Checkbox
DATE      → Date Picker
URL       → URL Input
```

Form renderer harus bekerja terhadap schema dari backend tanpa mengetahui form tersebut milik CSA atau organisasi lain.

---

# 24. UI/UX Requirements

Design direction:

**Modern SaaS dashboard**

Characteristics:

* Clean
* Professional
* Minimal
* Responsive
* Good whitespace
* Clear typography
* Strong information hierarchy
* Subtle interactions
* Not overly decorative

Avoid:

* Excessive gradients
* Excessive glassmorphism
* Huge decorative elements
* Unnecessary animations
* Generic template appearance
* Excessive rounded cards
* Too many colors

Preferred:

* White/light neutral background
* Dark text
* One primary accent color
* Clear status colors
* Subtle borders
* Consistent spacing
* Compact but readable dashboard

---

# 25. Animations & Micro-interactions

Animations should improve usability rather than become decoration.

Use subtle:

* Page transitions
* Modal transitions
* Dropdown transitions
* Hover states
* Button feedback
* Toast notifications
* Skeleton loading
* Form validation feedback
* Drag-and-drop feedback in Form Builder

Avoid:

* Long animations
* Excessive bouncing
* Decorative animations
* Animations that slow down workflows

---

# 26. Loading, Empty & Error States

Every data-driven page must handle:

### Loading

Show skeleton or loading indicator.

### Empty

Example:

```text
No events yet.

Create your first event to start accepting registrations.
```

### Error

Example:

```text
Something went wrong.

Please try again.
```

### Destructive Action

Delete/reject actions must require confirmation.

---

# 27. Seed Data

Application must contain realistic demo data.

## Organization

```text
Creative Student Association
Slug:
creative-student-association
```

## Divisions

```text
Programming
Multimedia
```

## Event

```text
Open Recruitment CSA 2026
```

Example registration fields:

```text
Nama Lengkap
NIM
Email
Nomor WhatsApp
Divisi Pilihan
Alasan Bergabung
Portfolio
```

Seed data should be sufficient to demonstrate the complete application flow.

---

# 28. Business Rules

## Organization

* Organization can have multiple divisions.
* Organization can have multiple events.
* Organization can have multiple members.

## Division

* Division belongs to exactly one organization.
* Division cannot belong to another organization.

## Event

* Event belongs to exactly one organization.
* Event can have one active registration form.
* Only published events can receive public registrations.
* Closed events cannot receive new registrations.

## Registration

* Registration belongs to one event.
* Registration contains answers based on the event's form.
* Registration status starts as `PENDING`.
* Admin can change status to `ACCEPTED` or `REJECTED`.

## Member

* Member belongs to one organization.
* Member may optionally belong to one division.
* Accepted registration can be converted into a member.
* The same registration must not create duplicate members.

---

# 29. Validation

Validation must exist on both frontend and backend.

Examples:

### Required Field

```text
required = true
```

must reject empty value.

### Email

Must use valid email format.

### Number

Must contain valid numeric value.

### URL

Must contain valid URL format.

### Select

Submitted option must exist in configured options.

### Event Registration

Backend must verify:

```text
event exists
event published
registration period active
quota available
form exists
required fields completed
field values valid
```

Never trust frontend validation alone.

---

# 30. Security Requirements

Backend must implement:

* Password hashing
* JWT/session authentication
* Role-based authorization
* Organization-level authorization
* Input validation
* SQL injection prevention
* CORS configuration
* Secure error handling
* No sensitive information in API responses

Public registration endpoints must be protected against obvious invalid input and abuse where appropriate.

---

# 31. Database Requirements

Use relational database.

Recommended:

```text
PostgreSQL
```

Foreign key relationships must be enforced where appropriate.

Important relationships:

```text
organizations.id
    ↓
divisions.organization_id

organizations.id
    ↓
events.organization_id

events.id
    ↓
forms.event_id

forms.id
    ↓
form_fields.form_id

events.id
    ↓
registrations.event_id

registrations.id
    ↓
registration_answers.registration_id

form_fields.id
    ↓
registration_answers.form_field_id

organizations.id
    ↓
members.organization_id

divisions.id
    ↓
members.division_id
```

---

# 32. Repository & Business Logic

Backend should separate:

```text
HTTP Handler
     ↓
Use Case / Service
     ↓
Repository
     ↓
Database
```

Example:

```text
RegistrationHandler
        ↓
CreateRegistrationUseCase
        ↓
RegistrationRepository
        ↓
PostgreSQL
```

Business rules should not be placed directly inside HTTP handlers.

---

# 33. Error Handling

Backend should use meaningful errors.

Examples:

```text
OrganizationNotFound
EventNotFound
FormNotFound
RegistrationNotFound
Unauthorized
Forbidden
InvalidFormField
RegistrationClosed
QuotaExceeded
DuplicateRegistration
```

Frontend should translate these into understandable messages.

---

# 34. Testing Requirements

MVP should include backend tests for critical business logic.

Minimum:

### Authentication

* Valid login
* Invalid login
* Unauthorized request

### Event

* Create event
* Publish event
* Cannot register to unpublished event
* Cannot register to closed event

### Form

* Create form
* Add fields
* Validate required fields
* Validate field types

### Registration

* Successful registration
* Invalid registration
* Registration after deadline
* Registration when quota is full

### Member

* Accept applicant
* Convert applicant to member
* Prevent duplicate conversion

---

# 35. MVP Completion Checklist

## Data

* [ ] Organization CRUD
* [ ] Dynamic divisions
* [ ] Event CRUD
* [ ] Database persistence
* [ ] Seed data
* [ ] Authentication

## Form Builder

* [ ] Add field
* [ ] Edit field
* [ ] Delete field
* [ ] Duplicate field
* [ ] Reorder field
* [ ] Required/optional
* [ ] Configure options
* [ ] Preview
* [ ] Save
* [ ] Publish

## Registration

* [ ] Public event page
* [ ] Dynamic form
* [ ] Client validation
* [ ] Server validation
* [ ] Submit registration
* [ ] Persist registration
* [ ] Success page

## Admin Review

* [ ] Registration list
* [ ] Search
* [ ] Filter
* [ ] Registration detail
* [ ] Accept
* [ ] Reject
* [ ] Rejection reason

## Members

* [ ] Convert applicant to member
* [ ] Member list
* [ ] Search
* [ ] Filter by division
* [ ] Member detail
* [ ] Edit member

## Dashboard

* [ ] Total members
* [ ] Active events
* [ ] Total registrations
* [ ] Pending applications
* [ ] Registration statistics

## Export

* [ ] Export registrations CSV
* [ ] Export members CSV
* [ ] Dynamic registration answers included

## UX

* [ ] Responsive
* [ ] Loading states
* [ ] Empty states
* [ ] Error states
* [ ] Toast notifications
* [ ] Confirmation dialogs
* [ ] Form validation feedback

---

# 36. Post-MVP Features

Features below are intentionally excluded from MVP.

Potential future features:

* Advanced analytics
* Email notifications
* File upload
* QR attendance
* Attendance management
* Certificate generation
* Event calendar
* Multiple organization administrators
* Advanced permission system
* Audit logs
* Password reset
* Email verification
* Applicant scoring
* Recruitment stages
* Interview scheduling
* Member activity tracking
* Organization announcements
* Member dues/payment management
* Certificate templates
* Report generation
* Dashboard customization

These features should only be implemented after the MVP core flow is stable.

---

# 37. Development Principles

## 37.1 Functional First

Prioritize working functionality over visual complexity.

Do not create a visually polished page with fake/mock interactions if the underlying functionality does not exist.

## 37.2 Generic Architecture

Do not hardcode:

```text
CSA
Programming
Multimedia
Open Recruitment
```

into business logic.

These are only seed/demo data.

The system must also work with:

```text
UKM A
UKM B
Organization C
```

without source code changes.

## 37.3 API First

Frontend and backend communicate exclusively through defined APIs.

Do not directly access the backend database from React.

## 37.4 Reusable Components

Build reusable components for:

```text
DataTable
Modal
FormField
FormRenderer
StatusBadge
SearchInput
Filter
Pagination
ConfirmDialog
Toast
EmptyState
LoadingState
```

## 37.5 Consistency

Maintain consistent:

* Naming
* API response structure
* Error handling
* Validation
* UI components
* Database conventions
* Code organization

---

# 38. Definition of Done

A feature is considered complete when:

1. Backend endpoint exists.
2. Backend validation exists.
3. Business logic is implemented.
4. Database persistence works.
5. Frontend consumes the real API.
6. Loading state exists.
7. Error state exists.
8. Empty state exists where applicable.
9. Authorization is implemented where required.
10. The feature works after page refresh.
11. No critical feature relies on hardcoded mock data.
12. Critical business logic has tests where applicable.

---

# 39. Final MVP Demonstration

The application should be demonstrable through this scenario:

### Step 1

Login as Organization Admin.

### Step 2

Open:

```text
Creative Student Association
```

### Step 3

View divisions:

```text
Programming
Multimedia
```

### Step 4

Create:

```text
Open Recruitment CSA 2026
```

### Step 5

Open Form Builder.

### Step 6

Create fields dynamically:

```text
Nama Lengkap
NIM
Email
Nomor WhatsApp
Divisi Pilihan
Alasan Bergabung
Portfolio
```

### Step 7

Publish event and registration form.

### Step 8

Open public event page.

### Step 9

Submit registration as applicant.

### Step 10

Return to admin dashboard.

### Step 11

Open Applications.

### Step 12

View applicant detail.

### Step 13

Accept applicant.

### Step 14

Convert applicant to member.

### Step 15

Open Member Management.

The newly created member must appear in the selected division.

---

# 40. Portfolio Product Statement

The core capability of this application can be summarized as:

> **An organization administrator can create an event, visually build a custom registration form, publish it, receive registrations, review applicants, accept them, and convert accepted applicants into organization members — without changing the application's source code.**

This capability represents the primary product value and should remain functional throughout development.
