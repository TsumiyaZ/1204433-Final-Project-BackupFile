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
- แต่ละ Feature ต้องใช้งานได้จริง และอธิบาย Flow ของตัวเองใน Technical Demo ได้
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
- [x] สร้างหน้า Setup สำหรับเลือก Source/Destination
- [x] เชื่อมหน้า Setup กับ Wails Methods
- [x] เปิด Native Directory Dialog สำหรับ Source และ Destination
- [x] บันทึกและโหลดค่า Source/Destination จาก SQLite ได้
- [x] สร้าง PhotoRepository สำหรับจัดการข้อมูลรูปภาพ
- [x] สร้าง `ScannedPhoto` DTO สำหรับส่งผลการสแกนไป Frontend
- [x] สร้าง PhotoService สำหรับสแกนเฉพาะไฟล์รูปจาก Source และโฟลเดอร์ย่อย
- [x] ตรวจนามสกุลไฟล์รูปที่รองรับ เช่น JPG, JPEG, PNG, WEBP และ GIF
- [x] สร้าง PhotoController และเมธอด `ScanPhotos(source)`
- [x] Bind PhotoController เข้า Wails และสร้าง Frontend Bindings
- [x] สร้างหน้า Scan Photos และแสดงรายการรูปที่พบ
- [x] เพิ่มโหมด Scan Source และ Scan Destination ในหน้าเดียวกัน
- [x] เลือกภาพทีละรูปและเลือกทั้งหมดเพื่อส่งต่อไปขั้นตอน Move
- [x] สร้าง `GetPhotoPreview(source, path)` และแสดง Preview ภาพจริงในหน้า Scan Photos
- [x] จำกัดขนาดไฟล์ Preview และตรวจว่าไฟล์อยู่ภายใน Source
- [x] โหลด Preview พร้อมกันสูงสุด 4 รูปเพื่อไม่ให้หน้าจอค้าง
- [x] ปรับ UI เป็นธีมครีม-ขาว พร้อม Navbar, Responsive Layout และสถานะการใช้งาน
- [x] กำหนด DTO และ IPC Contract ร่วมกับทีม

**Wails Methods:**

```text
GetSetting()
SaveSetting(source, dest)
SelectSourceDirectory()
SelectDestinationDirectory()
ScanPhotos(source)
GetPhotoPreview(source, path)
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

- [x] สร้างหน้า Move Photos
- [x] ตรวจ Extension และ File Signature
- [x] รับรายการรูปที่เลือกมาจาก Feature Scan Photos
- [x] Move ด้วย Rename เมื่ออยู่ไดรฟ์เดียวกัน
- [x] Copy + Verify + Delete เมื่ออยู่คนละไดรฟ์
- [x] คำนวณ SHA-256 Checksum
- [x] ป้องกันไฟล์ซ้ำ
- [x] ใช้ Worker Pool
- [x] ส่ง Progress Events ไป Frontend
- [x] บันทึก Photo หลังย้ายและ Verify สำเร็จ

**สถานะล่าสุด:** Implementation ครบแล้ว และผ่าน `go test ./...` กับ Nuxt Production Build ส่วนการทดสอบย้ายรูปด้วย Flash Drive จริงยังอยู่ในงานร่วมกัน

**งานเพิ่มตาม Flow วันที่ 2026-10-16 (ร่วมสมาชิก 4):**

- [ ] ส่ง Photo ID เข้าคิว AI หลังแต่ละรูปย้ายและบันทึก DB สำเร็จ
- [ ] ส่งเฉพาะรูปที่ย้ายสำเร็จ ไม่ส่งรูป skipped/failed และไม่รอจบทั้งชุด
- [ ] แยกผล Move ออกจากผล AI; AI ล้มเหลวไม่เปลี่ยนผลย้ายที่สำเร็จ

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
- `service/gallery_service.go`
- `dto/gallery_photo.go`
- `repository/photo_repository.go` (ใช้ query รูปร่วมกับสมาชิก 1)

**งาน:**

- [x] สร้างหน้า Gallery และเพิ่มเมนูรูปภาพใน Navbar
- [x] สร้าง GalleryPhoto DTO, GalleryService และ GalleryController
- [x] Bind GalleryController ใน main.go และสร้าง Wails Bindings
- [x] แสดงรายการชื่อรูปจาก DB ของ Destination ปัจจุบัน
- [x] ใช้ Photo records ใน Database เป็นรายการหลัก ไม่สแกนนำเข้ารูปเดิมใน Destination
- [x] แสดง Description, Status และ Analysis Status ที่อ่านจาก DB
- [ ] รับ Event อัปเดตสถานะ AI ตาม Photo ID
- [x] เพิ่ม GetPhotoByID โดยกรอง ID และ Destination ใน PhotoRepository
- [x] เพิ่ม GetGalleryPhotoPreview ใน Service โดยใช้ PhotoService อ่านไฟล์และแปลง Preview
- [ ] เปิด GetGalleryPhotoPreview ผ่าน Controller และสร้าง Wails Bindings ใหม่
- [ ] เชื่อม Frontend โหลด Preview ตาม Photo ID และแสดงภาพจริง
- [ ] Preview และเลือกรูปหลายรูป
- [x] ทำ Destination-scoped queries สำหรับรายการ Gallery และค้นรูปตาม ID
- [ ] ลบไฟล์จริงและลบ Record ใน Database
- [ ] ใช้ Transaction หรือกำหนด Rollback/Recovery ที่ชัดเจน
- [ ] ใช้ Foreign Key Cascade ลบ PhotoTag
- [ ] ตรวจ Integrity ว่าไฟล์จริงยังอยู่หรือไม่
- [x] ทำ Loading, Empty และ Error States สำหรับโหลดรายการ Gallery
- [ ] ตรวจไฟล์จริงและอัปเดตสถานะ `active` และ `missing`
- [ ] ทำ Confirmation สำหรับลบรูป

**สถานะล่าสุด (2026-10-16):** โหลดรายการจาก DB ขึ้นหน้า Gallery ได้แล้ว และ `go test ./...` ผ่านใน commit `8e940c0` ส่วน Preview มี method ใน Service แล้ว แต่ยังไม่เปิดผ่าน Controller หรือแสดงภาพใน Frontend; ยังไม่มี AI, Delete และ Integrity

**Wails Methods:**

```text
GetGalleryPhotos()                  // ใช้งานแล้ว
GetGalleryPhotoPreview(photoID)      // มีใน Service; Controller/Frontend ยังไม่เชื่อม
DeletePhotos(request)
CheckIntegrity(destination)
```

**งาน Demo:** อธิบาย Destination Scope, Delete Flow, Foreign Key Cascade และ Integrity Check

## สมาชิก 4 — AI Analysis และ Search

**Feature:** วิเคราะห์รูปที่ย้ายสำเร็จอัตโนมัติ สร้าง Description และค้นข้อความด้วย SQLite; Tags/Tag Cloud เป็นส่วนเสริม

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
- [ ] สร้างคิว AI แยกจาก Worker ย้ายรูป และเชื่อมกับสมาชิก 2
- [ ] รับ Photo ID หลังย้ายและบันทึก DB สำเร็จ แล้วอ่านรูปจาก Destination
- [ ] ส่งรูปให้ AI และ Validate Description
- [ ] แสดงผลและสถานะ AI ใน Gallery และทำหน้า Search
- [ ] จัดการ Analysis Status
- [ ] ทำ Retry และ Timeout
- [ ] ป้องกันการวิเคราะห์ซ้ำ
- [ ] บันทึก Description และสถานะ completed ให้สำเร็จใน Transaction เดียว
- [ ] กู้คืนงาน pending/processing ที่ค้างหลังปิดแอป โดยตรวจ record และไฟล์ก่อนเข้าคิวใหม่
- [ ] ค้นหาจาก Description
- [ ] ส่วนเสริม: สร้าง Tags และบันทึกพร้อม Description ใน Transaction เดียว
- [ ] ส่วนเสริม: ค้นหาจาก Tag และ Description + Tag
- [ ] ส่วนเสริม: แสดง Tag Cloud
- [ ] เตรียม Demo Dataset

**Wails Methods:**

```text
AnalyzePhoto(photoID, destination)
SearchPhotos(request)
GetTagCloud(destination)
```

`AnalyzePhoto` ใช้สำหรับ Retry ของรูปที่อยู่ในระบบ ส่วนการเริ่มวิเคราะห์ครั้งแรกให้ Backend เข้าคิวอัตโนมัติหลัง Move ไม่ต้องกด Analyze ใน Frontend และ `GetTagCloud` เป็นส่วนเสริม

**งาน Demo:** อธิบาย AI Request/Response, Transaction, Analysis Status, Tags, Search และ Error Handling

## Repository แยกตาม Feature

แผนเดิมแยก Repository ตาม Feature; Gallery ปัจจุบันใช้ `photo_repository.go` ร่วมกัน และยังไม่ได้สร้าง `gallery_repository.go` หากจะแยกเมื่อเพิ่ม Delete/Integrity ให้ตกลงกับ Owner ก่อน

โครงสร้างตามแผน:

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
- [ ] มี Loading, Empty และ Error States
- [ ] ไม่มี API Key หรือ `database.db` ถูก Push ขึ้น Git
- [ ] ผ่าน `gofmt` และ Build สำเร็จ
- [ ] Owner สามารถอธิบาย Flow ใน Technical Demo ได้

## งานร่วมกัน

- [x] ยืนยัน Requirements
- [x] ยืนยัน Database Schema
- [x] ยืนยัน DTO และ IPC Contract
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
- Merge เฉพาะงานที่ผ่านการ Review และ Build
- อัปเดต [[03-Milestones]]
