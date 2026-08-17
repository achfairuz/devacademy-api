# Task: Endpoint Course Card (`GET /api/v1/courses/cards`)

Endpoint baru untuk menampilkan daftar kursus dalam bentuk **card** yang berisi info ringkas yang dibutuhkan halaman daftar kursus (misal `CourseView` di frontend) dan "Lanjutkan Belajar" di Home.

## Ringkasan

| Field | Sumber | Keterangan |
|---|---|---|
| `category` | `Course.Category` | Nama kategori |
| `title` | `Course.Title` | Judul kursus |
| `description` | `Course.Description` | Deskripsi kursus |
| `progress` | `LessonProgress` | % lesson selesai untuk user yang login (0 kalau belum pernah dikerjain / belum login) |
| `total_modules` | `Course.Sections` | Jumlah section (bab) |
| `total_bought` | `Enrollment` | Jumlah user yang membeli/enroll kursus |
| `total_rated` | — | Skip dulu (belum ada fitur rating) → selalu `0` |

## Scope

- Semua kursus `status = published` (public), pagination `page` / `page_size`.
- Opsional kirim JWT (Authorization header) → `progress` terisi sesuai user tersebut. Tanpa JWT → `progress` = `0`.

## Keputusan Desain

- **Total modul = jumlah section (bab)**, mengikuti `course.modulesData.length` di frontend `CourseView.vue`.
- **Progress = lesson selesai / total lesson * 100**, dihitung dari `Enrollment → LessonProgress (is_completed)`. Bukan per-section.
- **Total rated = 0** karena model rating belum ada di backend.

## File yang Diubah

| File | Perubahan |
|---|---|
| `internal/middleware/middleware.go` | Tambah `OptionalAuth(secret)`: parse JWT kalau header ada & valid → set `user_id`; kalau header hilang / token invalid → tetap lanjut (tidak abort). |
| `internal/functions/course/course_repository.go` | Tambah di interface `CourseRepository` + implementasi: `FindCards`, `CountEnrollments`, `CountLessons`, `CountCompletedLessons`. |
| `internal/functions/course/course_response.go` | Tambah struct `CourseCard` + mapper `toCourseCard` (reuse `CategorySummary`, `LevelSummary`, `UserSummary`). |
| `internal/functions/course/course_service.go` | Tambah `GetCards(ctx, userID, page, pageSize)` di interface `CourseService` + implementasi. |
| `internal/functions/course/course_controller.go` | Tambah handler `GetCards` + doc swagger. |
| `internal/routes/routes.go` | Daftarkan `courses.GET("/cards", middleware.OptionalAuth(...), ctr.Course.GetCards)`. |

## Detail Implementasi Repository

```
FindCards(ctx, limit, offset) []models.Course
  - WHERE status = 'published', ORDER BY created_at DESC, LIMIT/OFFSET
  - Preload: Category, Level, Mentor, Sections

CountEnrollments(ctx, courseIDs) map[courseID]int
  - SELECT course_id, COUNT(*) FROM enrollments
    WHERE course_id IN (...) GROUP BY course_id

CountLessons(ctx, courseIDs) map[courseID]int
  - SELECT cs.course_id, COUNT(*) FROM lessons l
    JOIN course_sections cs ON l.section_id = cs.id
    WHERE cs.course_id IN (...) GROUP BY cs.course_id

CountCompletedLessons(ctx, userID, courseIDs) map[courseID]int
  - SELECT cs.course_id, COUNT(*) FROM lesson_progress lp
    JOIN enrollments e ON lp.enrollment_id = e.id
    JOIN lessons l ON lp.lesson_id = l.id
    JOIN course_sections cs ON l.section_id = cs.id
    WHERE e.user_id = ? AND lp.is_completed = TRUE
      AND cs.course_id IN (...) GROUP BY cs.course_id
```

Semua query agregat pakai `IN (...)` + `GROUP BY` supaya tidak terjadi N+1.

## Detail Service

```
GetCards(ctx, userID, page, pageSize):
  1. courses = repo.FindCards(pageSize, (page-1)*pageSize)
  2. stats  = repo.CountEnrollments(courseIDs)
  3. progressMap:
       - kalau userID != Nil:
           total     = repo.CountLessons(courseIDs)
           completed = repo.CountCompletedLessons(userID, courseIDs)
           progress = round(completed / total * 100), 0 kalau total = 0
       - kalau userID == Nil: progress = 0
  4. Map ke []CourseCard dengan total_modules = len(course.Sections)
```

## Detail Response DTO

```go
type CourseCard struct {
    ID            uuid.UUID        `json:"id"`
    Title         string           `json:"title"`
    Slug          string           `json:"slug"`
    Description   string           `json:"description"`
    Thumbnail     string           `json:"thumbnail"`
    Price         float64          `json:"price"`
    Duration      int              `json:"duration"`
    Level         *LevelSummary    `json:"level"`
    Mentor        *UserSummary     `json:"mentor"`
    Category      *CategorySummary `json:"category"`
    TotalModules  int              `json:"total_modules"`
    Progress      int              `json:"progress"`
    TotalBought   int              `json:"total_bought"`
    TotalRated    int              `json:"total_rated"`
}
```

## Routing

```go
courses.GET("/cards", middleware.OptionalAuth(cfg.JWT.Secret), ctr.Course.GetCards)
```

Catatan: static route `/cards` memiliki prioritas lebih tinggi dari `/:id` di gin, jadi aman dari konflik.

## Contoh Response

```
GET /api/v1/courses/cards?page=1&page_size=10

{
  "status": 200,
  "message": "course cards retrieved",
  "data": [
    {
      "id": "79eaf3e2-...",
      "title": "Belajar Golang dari Dasar",
      "slug": "belajar-golang-dari-dasar",
      "description": "Belajar bahasa pemrograman Go mulai dari nol...",
      "thumbnail": "",
      "price": 150000,
      "duration": 12,
      "level": { "id": "...", "name": "Beginner", "slug": "beginner" },
      "mentor": { "id": "...", "full_name": "Mentor DevAcademy", ... },
      "category": { "id": "...", "name": "Backend Development", "slug": "backend-development", "icon": "server" },
      "total_modules": 2,
      "progress": 0,
      "total_bought": 0,
      "total_rated": 0
    }
  ]
}
```

## Catatan / Limitasi

- `total_bought` dan `progress` masih `0` karena belum ada data `Enrollment` / `LessonProgress` di DB (belum ada endpoint enrollment juga). Opsional: tambahkan seed demo di `internal/database/seeder.go`.
- `total_rated` sengaja di-skip; saat fitur rating dibuat, tinggal ganti nilai `0` dengan query `course_ratings`.
- Endpoint lain (`GetAll`, `GetByMentor`, dll) tidak terpengaruh karena endpoint card menggunakan query terpisah.

## Verifikasi

- `go build ./...` — pastikan tidak ada error.
- Jalankan server, cek:
  - `GET /api/v1/courses/cards` → 200, list card.
  - Tanpa JWT → `progress: 0`.
  - Dengan JWT student yang punya `LessonProgress` → `progress` terisi.
