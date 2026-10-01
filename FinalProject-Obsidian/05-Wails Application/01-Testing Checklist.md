---
title: Testing Checklist
tags:
  - testing
  - qa
status: active
---

# Testing Checklist

กลับไปที่ [[01-Cross Platform App/00-Project Home|Project Home]]

## Setup และ State Guard

- [ ] Scan โดยยังไม่เลือก Source
- [ ] Move โดยยังไม่เลือก Destination
- [ ] เลือก Path ที่ไม่มีอยู่
- [ ] เปิดแอปใหม่แล้ว Settings ยังอยู่

## Scan

- [ ] JPG/JPEG
- [ ] PNG
- [ ] WebP
- [ ] ข้าม PDF/ZIP/เอกสาร
- [ ] ข้ามไฟล์ปลอมที่เปลี่ยนนามสกุลเป็นรูป
- [ ] ชื่อไฟล์ภาษาไทยและช่องว่าง

## Move

- [ ] รูปเดียว
- [ ] หลายรูป
- [ ] Select All
- [ ] Duplicate Filename
- [ ] ไม่เขียนทับไฟล์เดิม
- [ ] Same-drive Rename
- [ ] Cross-drive Copy + Verify + Delete
- [ ] Copy ล้มเหลวแล้ว Source ยังอยู่
- [ ] ถอด Flash Drive ระหว่างทำงาน
- [ ] Progress และเวลารวม

## Database และ Integrity

- [ ] สร้าง Photo หลัง Move สำเร็จเท่านั้น
- [ ] History แยกตาม Destination
- [ ] ลบจาก Disk แล้วพบ `missing`
- [ ] คืนไฟล์แล้วกลับเป็น `active`
- [ ] Delete ผ่านแอปแล้วไฟล์จริง, `photos` และ `photo_tags` ถูกลบ

## AI

- [ ] วิเคราะห์หนึ่งรูป
- [ ] วิเคราะห์หลายรูป
- [ ] บันทึก Description
- [ ] Tags ไม่ซ้ำ
- [ ] Timeout/API Error → `failed`
- [ ] Retry สำเร็จ
- [ ] AI ล้มเหลวแล้วรูปไม่หาย

## Search

- [ ] Description แบบบางส่วน
- [ ] Tag เดียว
- [ ] หลาย Tags
- [ ] Description + Tags
- [ ] Tag Cloud Count ถูกต้อง
- [ ] รูปที่ลบผ่านแอปไม่ปรากฏในผลค้นหา
- [ ] แจ้งเตือน `missing`

## Cross-platform

- [ ] Windows Development Build
- [ ] Windows Production Build
- [ ] macOS Development Build
- [ ] macOS Production Build
- [ ] Directory Dialog ทั้งสองระบบ
- [ ] Path Separator และชื่อ Unicode
