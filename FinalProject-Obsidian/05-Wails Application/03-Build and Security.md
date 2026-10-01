---
title: Build and Security
tags:
  - build
  - security
  - cross-platform
status: planning
---

# Build and Security

กลับไปที่ [[01-Cross Platform App/00-Project Home|Project Home]]

## Build Checklist

- [ ] Nuxt Development Mode
- [ ] Nuxt Production Output ตรงกับ Path ที่ Wails Embed
- [ ] Wails Development Mode
- [ ] Wails Windows Build
- [ ] Wails macOS Build บนเครื่อง macOS จริง
- [ ] Generated Wails Bindings ตรงกับ Backend Methods

## Security Rules

- API Key อยู่ใน `.env` หรือ Secure App Configuration
- `.env` ต้องอยู่ใน `.gitignore`
- ไม่ Log API Key หรือ Raw Token
- Validate ทุก Path ใน Backend
- ป้องกัน `..` Path Traversal
- ตรวจว่า Selected File อยู่ใต้ Source
- ห้ามเขียนทับ Destination โดยไม่ยืนยัน
- จำกัดขนาดไฟล์ที่ส่ง AI หาก Provider มีข้อจำกัด

## Cross-platform Notes

- Windows และ macOS ใช้ Path/Volume ต่างกัน
- Cross-device Rename อาจล้มเหลวและต้อง Fallback
- Case sensitivity ของ Filename ต่างกันตาม File System
- ต้องทดสอบ Unicode และชื่อภาษาไทย
- การเข้าถึงโทรศัพท์โดยตรงอาจไม่ใช่ File System ปกติ

## Release Checklist

- [ ] Working tree สะอาด
- [ ] Tests ผ่าน
- [ ] ไม่มี Secret ใน Git History ล่าสุด
- [ ] Database Migration ผ่าน
- [ ] Demo Dataset พร้อม
- [ ] Version/Commit ถูกบันทึก
- [ ] Windows/macOS Artifacts เปิดได้
