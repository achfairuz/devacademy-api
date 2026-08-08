# Fitur yang Belum Dibuat

Daftar fitur yang model-nya sudah ada di database (`internal/models`) tetapi **belum** memiliki endpoint/controller/service/repository, lengkap dengan prioritas pengerjaan.

## Tabel Prioritas

| # | Feature | Priority | Model | Est. Endpoints | Ketergantungan |
|---|---|---|---|---|---|
| 1 | Enrollment & LessonProgress | 🔴 High | `Enrollment`, `LessonProgress` | 6 | Fondasi fitur belajar (dipakai Certificate) |
| 2 | Quiz Attempt | 🔴 High | `QuizAttempt` | 3 | Quiz CRUD (sudah ada) |
| 3 | Assignment Submission | 🔴 High | `AssignmentSubmission` | 5 | Assignment CRUD (sudah ada) |
| 4 | Certificate | 🟡 Medium | `Certificate` | 3 | Enrollment + LessonProgress (100%) |
| 5 | Subscription & Payment | 🟢 Low | `SubscriptionPlan`, `UserSubscription`, `Payment` | 8 | Bisa jalan paralel / menyusul |

---

## 1. Enrollment & LessonProgress — 🔴 High

| No | Model | Status |
|---|---|---|
| 1 | `Enrollment` | Belum ada endpoint |
| 2 | `LessonProgress` | Belum ada endpoint |

- Student mendaftar ke kursus (unique per user+course) dan melacak progres belajar.
- Endpoint yang dibutuhkan:

| Method | Endpoint | Akses |
|---|---|---|
| POST | `/api/v1/enrollments` | Student |
| GET | `/api/v1/enrollments` | Student |
| GET | `/api/v1/enrollments/:id` | Student/Mentor |
| DELETE | `/api/v1/enrollments/:id` | Student |
| PUT | `/api/v1/enrollments/:id/lessons/:lesson_id/progress` | Student |
| GET | `/api/v1/enrollments/:id/progress` | Student/Mentor |

## 2. Quiz Attempt — 🔴 High

| No | Model | Status |
|---|---|---|
| 1 | `QuizAttempt` | Belum ada endpoint |

- Student submit jawaban quiz, sistem hitung skor vs `passing_score`.
- Endpoint yang dibutuhkan:

| Method | Endpoint | Akses |
|---|---|---|
| POST | `/api/v1/courses/:id/sections/:section_id/lessons/:lesson_id/quizzes/:quiz_id/attempts` | Student |
| GET | `/api/v1/.../quizzes/:quiz_id/attempts` | Student/Mentor |
| GET | `/api/v1/.../attempts/:attempt_id` | Student/Mentor |

- Catatan: saat ini `is_correct` pada option ikut tampil di response publik — perlu filter per role Student.

## 3. Assignment Submission — 🔴 High

| No | Model | Status |
|---|---|---|
| 1 | `AssignmentSubmission` | Belum ada endpoint |

- Student upload file jawaban, Mentor beri nilai (`score`) + `feedback`.
- Endpoint yang dibutuhkan:

| Method | Endpoint | Akses |
|---|---|---|
| POST | `/api/v1/.../assignments/:assignment_id/submissions` | Student |
| GET | `/api/v1/.../assignments/:assignment_id/submissions` | Mentor |
| GET | `/api/v1/.../submissions/:submission_id` | Student/Mentor |
| PUT | `/api/v1/.../submissions/:submission_id` | Mentor |
| DELETE | `/api/v1/.../submissions/:submission_id` | Mentor |

## 4. Certificate — 🟡 Medium

| No | Model | Status |
|---|---|---|
| 1 | `Certificate` | Belum ada endpoint |

- Sertifikat dibuat otomatis saat user lulus kursus (100% progress + lolos quiz/assignment).
- Endpoint yang dibutuhkan:

| Method | Endpoint | Akses |
|---|---|---|
| GET | `/api/v1/certificates` | Student |
| GET | `/api/v1/certificates/:id` | Student |
| GET | `/api/v1/certificates/verify/:number` | Public |

## 5. Subscription & Payment — 🟢 Low

| No | Model | Status |
|---|---|---|
| 1 | `SubscriptionPlan` | Belum ada endpoint |
| 2 | `UserSubscription` | Belum ada endpoint |
| 3 | `Payment` | Belum ada endpoint |

- CRUD paket subscription (admin), user subscribe, dan pembayaran.
- Endpoint yang dibutuhkan:

| Method | Endpoint | Akses |
|---|---|---|
| GET | `/api/v1/subscriptions` | Public |
| GET | `/api/v1/subscriptions/:id` | Public |
| POST | `/api/v1/subscriptions` | Admin |
| PUT/DELETE | `/api/v1/subscriptions/:id` | Admin |
| POST | `/api/v1/subscriptions/:id/subscribe` | Student |
| GET | `/api/v1/users/:id/subscription` | Student |
| POST | `/api/v1/payments` | Student |
| GET | `/api/v1/payments` | Student |

---

## Ringkasan Status

| Domain | Model | Implementasi | Priority |
|---|---|---|---|
| User/Auth | `User`, `Role` | ✅ Selesai | — |
| Master | `Category`, `Level` | ✅ Selesai | — |
| Konten Kursus | `Course`, `CourseSection`, `Lesson`, `LessonFile` | ✅ Selesai | — |
| Evaluasi | `Quiz`, `QuizQuestion`, `QuizOption`, `Assignment` | ✅ CRUD selesai | — |
| Pembelajaran | `Enrollment`, `LessonProgress` | ❌ Belum | 🔴 High |
| Evaluasi | `QuizAttempt`, `AssignmentSubmission` | ❌ Belum | 🔴 High |
| Kelulusan | `Certificate` | ❌ Belum | 🟡 Medium |
| Monetisasi | `SubscriptionPlan`, `UserSubscription`, `Payment` | ❌ Belum | 🟢 Low |

## Perbaikan yang Disarankan

- Filter `is_correct` pada option quiz agar tidak bocor ke role Student (misal dengan response DTO terpisah).
- Query pagination yang konsisten di semua endpoint list (page/page_size).
- Penjadwalan kerja: kerjakan **Enrollment → Quiz Attempt → Assignment Submission** lebih dulu, baru Certificate, lalu Subscription/Payment.
