---
title: Screens and UX
tags:
  - frontend
  - ux
status: planning
---

# Screens and UX

กลับไปที่ [[00-Frontend Overview]]

## Setup Page

- Source Path + Browse
- Destination Path + Browse
- Path Status
- จำนวนรูปที่ Scan ล่าสุด
- ปุ่ม Continue to Move Photos

## Move Photos Page

- Scan/Refresh
- Count: files, supported images, ignored files
- Preview Grid
- Checkbox แต่ละรูป
- Select All/Clear
- Move Selected
- Warning ว่า Source จะถูกลบหลังย้ายสำเร็จ
- Progress และ Result Summary

## Gallery & AI Page

- Photo Grid จาก Database
- Filter ตาม Analysis Status
- Multi-select
- Analyze Selected
- Description และ Tag Chips
- Retry
- Missing Status; เมื่อลบสำเร็จให้เอา Card ออกจาก Gallery

## Search Page

- Description Input
- Tag Filters
- Tag Cloud
- Search/Clear
- Result Grid
- Empty Result State

## History / Integrity Page

- Current Destination
- Run Integrity Check
- Active/Missing Counts
- รายการไฟล์ผิดปกติ
- ปุ่ม Refresh

## UX Copy ที่ต้องชัดเจน

- `Move` หมายถึงไฟล์ต้นฉบับจะถูกลบหลังตรวจสอบสำเร็จ
- `Skipped` ต้องบอกเหตุผล เช่น Duplicate หรือ Unsupported Type
- `Failed` ต้องบอกว่ายังเก็บต้นฉบับไว้หรือไม่
- AI Failure ต้องระบุว่ารูปยังปลอดภัย

## Accessibility

- ปุ่มมี Label ชัดเจน
- Status ไม่สื่อด้วยสีเพียงอย่างเดียว
- Keyboard focus มองเห็นได้
- Dialog ต้องมี Cancel
- Thumbnail ต้องมี Filename เป็น alt text
