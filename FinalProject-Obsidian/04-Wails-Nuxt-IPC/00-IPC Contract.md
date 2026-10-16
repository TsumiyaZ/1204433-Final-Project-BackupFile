---
title: Wails IPC Contract
tags:
  - wails
  - ipc
  - contract
status: draft
---

# Wails IPC Contract

กลับไปที่ [[01-Cross Platform App/00-Project Home|Project Home]]

> [!important]
> ทุกคนต้องตกลงชื่อ Method, Request, Response และ Error Format ก่อนเริ่ม Integration

## Proposed Methods

| Method | Purpose |
|---|---|
| `SelectSource()` | เปิด Directory Dialog และบันทึก Source |
| `SelectDestination()` | เปิด Directory Dialog และบันทึก Destination |
| `GetSettings()` | อ่าน Source/Destination ล่าสุด |
| `ScanSourceImages()` | Scan และคืนรายการรูป |
| `MovePhotos(selectedPaths)` | ย้ายรูปที่เลือก |
| `ListDestinationPhotos()` | อ่าน Gallery ตาม Destination |
| `AnalyzePhotos(photoIDs)` | Retry งานวิเคราะห์ของ Photo records ที่มีอยู่ |
| `SearchPhotos(text, tags)` | ค้น Description/Tags |
| `GetTagCloud()` | คืน Tag และจำนวนรูป |
| `RunIntegrityCheck()` | ตรวจ Database เทียบ Disk |
| `DeletePhotos(photoIDs)` | ลบไฟล์จริง พร้อมลบ `photos` และ `photo_tags` |

## Shared DTOs

### Photo Summary

- ID
- Filename
- Full/Display Path
- Thumbnail URL/Data
- Status
- Analysis Status
- Description
- Tags

### Operation Result

- Success
- Message
- Requested Count
- Success Count
- Skipped Count
- Failed Count
- Duration Milliseconds
- Items/Errors

### Error Format

- Code เช่น `SOURCE_NOT_SET`
- Message สำหรับผู้ใช้
- Detail สำหรับ Log
- Recoverable หรือไม่

## Contract Rules

- Flow หลัก: แต่ละรูป Move สำเร็จ → บันทึก Photo → Backend เข้าคิว AI โดยอัตโนมัติ ไม่ต้องให้ Frontend เรียก Analyze
- Gallery อ่าน records ตาม Destination; ไม่เพิ่ม records หรือวิเคราะห์ไฟล์เดิมด้วยการ Scan Destination
- Description และการค้นข้อความใน SQLite เป็น MVP; Tags/Tag Cloud เป็นส่วนเสริม ไม่ใช้ Vector Search
- ชื่อ Method ในตารางเป็นข้อเสนอ ส่วน Move ปัจจุบันคือ `MovePhotos(request)` ตาม DTO `MovePhotosRequest`

- ห้ามส่ง raw database error ให้ผู้ใช้
- Path validation ต้องอยู่ใน Backend
- Long-running call ต้องมี Progress Events
- Frontend ต้องไม่สร้างสถานะสำคัญเอง
- เปลี่ยน Contract แล้วต้องแจ้งทุกคนและอัปเดต Note นี้

อ่านต่อ: [[01-Events and Progress]]
