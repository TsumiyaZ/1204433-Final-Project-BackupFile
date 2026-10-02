---
title: Team Division
tags:
  - team
  - task-allocation
  - vertical-slice
status: active
---

# การแบ่งงาน 4 คนแบบ Full Stack

กลับไปที่ [[00-Project Home]]

## รูปแบบการแบ่งงาน

แบ่งงานแบบ **Vertical Slice** โดยแต่ละคนรับผิดชอบหนึ่ง Feature ตั้งแต่หน้า UI จนถึงฐานข้อมูล

```text
Frontend → Controller → Service → Repository → SQLite
```

หลักการร่วมกัน:

- ทุกคนต้องได้เขียนทั้ง Nuxt Frontend และ Go Backend
- Controller ห้ามเรียก Repository โดยตรง ต้องผ่าน Service เสมอ
- Repository ใช้ GORM Gen ผ่าน `query.Query`
- Query ของรูปภาพต้องกรองตาม Destination
- แต่ละ Feature ต้องมี Tests และอธิบาย Flow ของตัวเองใน Technical Demo ได้
- ตกลง DTO และ Wails IPC Contract ก่อนเชื่อม Frontend

## สมาชิก 1 — Setup, Settings, Scan Photos และ Integration

**Feature:** กำหนด Source/Destination, สแกนรูปภาพจาก Source และประกอบ Backend Modules

**พื้นที่รับผิดชอบ:**

- `frontend/app/pages/setup.vue`
- `controller/setting_controller.go`
- `controller/photo_controller.go`
- `service/setting_service.go`
- `service/photo_service.go`
- `repository/setting_repository.go`
- `repository/photo_repository.go`
- `model/`, `model/query/`
- `repository/dbconnect.go`, `repository/migrate.go`
- `main.go`, `app.go`, `go.mod`, `go.sum`

**งาน:**

- [x] สร้าง GORM Models: Setting, Photo, Tag และ PhotoTag
- [x] เชื่อม SQLite และกำหนดตำแหน่ง `database/database.db`
- [x] ทำ AutoMigrate
- [x] สร้าง GORM Gen Query
- [x] สร้าง SettingRepository
- [x] สร้าง SettingService และ Validation
- [x] สร้าง SettingController
- [x] Bind SettingController เข้า Wails
- [x] เขียน SettingService Test
- [x] สร้างหน้า Setup สำหรับเลือก Source/Destination
- [x] เชื่อมหน้า Setup กับ Wails Methods
- [x] เปิด Native Directory Dialog สำหรับเลือก Source และ Destination
- [x] บันทึกและโหลดค่า Source/Destination จาก SQLite ได้
- [ ] สร้าง PhotoRepository สำหรับจัดการข้อมูลรูปภาพ
- [ ] สร้าง PhotoService สำหรับสแกนเฉพาะไฟล์รูปจาก Source
- [ ] ตรวจนามสกุลไฟล์รูปที่รองรับ เช่น JPG, JPEG, PNG, WEBP และ GIF
- [ ] สร้าง PhotoController และเมธอด `ScanPhotos(source)`
- [ ] Bind PhotoController เข้า Wails
- [ ] สร้างหน้า Scan Photos และแสดงรายการรูปที่พบ
- [ ] แสดง Preview และเลือกภาพที่จะส่งต่อไปขั้นตอน Move
- [ ] เขียน Tests สำหรับการ Scan ด้วย Temporary Directory
- [ ] กำหนด DTO และ IPC Contract ร่วมกับทีม
- [ ] รวม Backend Modules ของสมาชิกทุกคนใน `main.go`
- [ ] Review Pull Requests ก่อนเข้า `develop`

**Wails Methods:**

```text
GetSetting()
SaveSetting(source, dest)
SelectSourceDirectory()
SelectDestinationDirectory()
ScanPhotos(source)
```

**งาน Demo:** อธิบาย Layered Architecture, SQLite, GORM, GORM Gen, Models, Relations, Dependency Injection, Wails Bind และ Flow การ Scan รูปภาพ

## สมาชิก 2 — Move Photos

**Feature:** รับรายการรูปที่ผู้ใช้เลือก แล้วทำการย้ายจาก Source ไป Destination

**พื้นที่รับผิดชอบ:**

- `frontend/app/pages/move.vue`
- Components สำหรับ Preview Grid, Selection และ Progress
- `controller/backup_controller.go`
- `service/backup/`
- `repository/backup_repository.go`

**งาน:**

- [ ] สร้างหน้า Move Photos
- [ ] ตรวจ Extension และ File Signature
- [ ] รับรายการรูปที่เลือกมาจาก Feature Scan Photos
- [ ] Move ด้วย Rename เมื่ออยู่ไดรฟ์เดียวกัน
- [ ] Copy + Verify + Delete เมื่ออยู่คนละไดรฟ์
- [ ] คำนวณ SHA-256 Checksum
- [ ] ป้องกันไฟล์ซ้ำ
- [ ] ใช้ Worker Pool
- [ ] ส่ง Progress Events ไป Frontend
- [ ] บันทึก Photo หลังย้ายและ Verify สำเร็จ
- [ ] เขียน Unit Tests ด้วย Temporary Directory

**Wails Methods และ Events:**

```text
MovePhotos(request)

backup:progress
backup:completed
backup:error
```

**งาน Demo:** อธิบาย Scan, Rename, Cross-device fallback, Checksum, Duplicate Protection และ Progress Events

## สมาชิก 3 — Gallery, Delete และ Integrity

**Feature:** ดูรูปที่ Backup แล้ว ลบรูป และตรวจสอบไฟล์

**พื้นที่รับผิดชอบ:**

- `frontend/app/pages/gallery.vue`
- Components สำหรับ Gallery, Preview และ Confirmation Dialog
- `controller/gallery_controller.go`
- `service/gallery/`
- `repository/gallery_repository.go`

**งาน:**

- [ ] สร้างหน้า Gallery
- [ ] แสดงรูปทั้งหมดของ Destination ปัจจุบัน
- [ ] Preview และเลือกรูปหลายรูป
- [ ] ทำ Destination-scoped queries
- [ ] ลบไฟล์จริงและลบ Record ใน Database
- [ ] ใช้ Transaction หรือกำหนด Rollback/Recovery ที่ชัดเจน
- [ ] ใช้ Foreign Key Cascade ลบ PhotoTag
- [ ] ตรวจ Integrity ว่าไฟล์จริงยังอยู่หรือไม่
- [ ] แสดงสถานะ `active` และ `missing`
- [ ] ทำ Loading, Empty, Error และ Confirmation States
- [ ] เขียน Unit Tests สำหรับ Gallery/Delete/Integrity

**Wails Methods:**

```text
GetPhotos(destination)
DeletePhotos(request)
CheckIntegrity(destination)
```

**งาน Demo:** อธิบาย Destination Scope, Delete Flow, Foreign Key Cascade และ Integrity Check

## สมาชิก 4 — AI Analysis และ Search

**Feature:** วิเคราะห์รูป สร้าง Description/Tags และค้นหารูป

**พื้นที่รับผิดชอบ:**

- `frontend/app/pages/ai-search.vue`
- Components สำหรับ AI Result, Search และ Tag Cloud
- `controller/analyzer_controller.go`
- `controller/search_controller.go`
- `service/analyzer/`
- `service/search/`
- `repository/analysis_repository.go`

**งาน:**

- [ ] เลือก AI Provider และกำหนด Analyzer Interface
- [ ] สร้างหน้า AI Analysis และ Search
- [ ] เลือกรูปจาก Destination เพื่อวิเคราะห์
- [ ] ส่งรูปให้ AI และ Parse Description/Tags
- [ ] จัดการ Analysis Status
- [ ] ทำ Retry และ Timeout
- [ ] ป้องกันการวิเคราะห์ซ้ำ
- [ ] บันทึก Description และ Tags ใน Transaction เดียว
- [ ] ค้นหาจาก Description
- [ ] ค้นหาจาก Tag
- [ ] ค้นหาแบบ Description + Tag
- [ ] แสดง Tag Cloud
- [ ] เตรียม Demo Dataset และเขียน Tests

**Wails Methods:**

```text
AnalyzePhoto(photoID, destination)
SearchPhotos(request)
GetTagCloud(destination)
```

**งาน Demo:** อธิบาย AI Request/Response, Transaction, Analysis Status, Tags, Search และ Error Handling

## Repository แยกตาม Feature

เพื่อให้ทุกคนได้ทำ Data Access และไม่แก้ไฟล์เดียวกันพร้อมกัน ให้แบ่ง Repository เป็น:

```text
repository/
├── dbconnect.go
├── migrate.go
├── setting_repository.go
├── photo_repository.go
├── backup_repository.go
├── gallery_repository.go
└── analysis_repository.go
```

แต่ละ Repository ใช้ Model และ Generated Query ชุดเดียวกันได้:

```go
r.q.Photo
r.q.PhotoTag
r.q.Setting
r.q.Tag
```

## Ownership เพื่อลด Merge Conflict

| พื้นที่ | Owner หลัก |
|---|---|
| Models, Migration, GORM Gen | สมาชิก 1 |
| `main.go`, `app.go`, `go.mod`, `go.sum` | สมาชิก 1 |
| Setup/Settings และ Scan Photos Feature | สมาชิก 1 |
| Move/Backup Feature | สมาชิก 2 |
| Gallery/Delete/Integrity Feature | สมาชิก 3 |
| AI/Search Feature | สมาชิก 4 |
| IPC Contract และ DTO | ทุกคนตกลงร่วมกัน สมาชิก 1 เป็นผู้รวม |

หากต้องแก้ไฟล์ของ Owner คนอื่น ให้ทำผ่าน Pull Request และแจ้งก่อนแก้

## Branches

```text
feature/settings
feature/scan-photos
feature/backup-move
feature/gallery-integrity
feature/ai-search
```

## Definition of Done ของทุก Feature

- [ ] Frontend ใช้งานได้จริง
- [ ] Controller ถูก Bind กับ Wails
- [ ] Service มี Business Logic และ Validation
- [ ] Repository ใช้ GORM Gen
- [ ] Queries ของรูปกรองตาม Destination
- [ ] มี Unit Tests หรือ Integration Tests
- [ ] มี Loading, Empty และ Error States
- [ ] ไม่มี API Key หรือ `database.db` ถูก Push ขึ้น Git
- [ ] ผ่าน `gofmt` และ `go test ./...`
- [ ] Owner สามารถอธิบาย Flow ใน Technical Demo ได้

## งานร่วมกัน

- [ ] ยืนยัน Requirements
- [ ] ยืนยัน Database Schema
- [ ] ยืนยัน DTO และ IPC Contract
- [ ] Review Code ของเพื่อน
- [ ] ทดสอบบน Windows
- [ ] ทดสอบบน macOS จริง
- [ ] ทดลองกับ Flash Drive จริง
- [ ] ตรวจ API Key และไฟล์ Database ก่อน Push
- [ ] ซ้อม Progress Demo
- [ ] ซ้อม Technical Demo

## Weekly Sync Checklist

- แต่ละคนรายงาน Feature ที่ทำเสร็จ
- ระบุ Blocker และ Contract ที่เปลี่ยน
- ตรวจว่าไม่มีการแก้ไฟล์ชนกัน
- Merge เฉพาะงานที่ผ่าน Tests
- อัปเดต [[03-Milestones]]
