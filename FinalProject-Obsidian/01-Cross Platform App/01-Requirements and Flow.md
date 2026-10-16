---
title: Requirements and Flow
tags:
  - requirements
  - workflow
status: planning
---

# Requirements and Flow

กลับไปที่ [[00-Project Home]]

## ขอบเขต

ระบบย้ายเฉพาะรูปที่รองรับ ไม่ย้าย PDF, ZIP, เอกสาร หรือไฟล์ประเภทอื่น

### รูปแบบไฟล์ MVP

- JPG/JPEG
- PNG
- WebP

HEIC, GIF, BMP และ TIFF เป็น Feature เพิ่มภายหลัง

## User Flow

```mermaid
flowchart TD
    A[เลือก Source และ Destination] --> B[เข้าเมนู Move Photos]
    B --> C[Scan รูปใน Source]
    C --> D[แสดง Preview]
    D --> E[เลือกบางรูปหรือ Select All]
    E --> F{ตรวจ Duplicate}
    F -->|ซ้ำ| G[Skip และแจ้งเหตุผล]
    F -->|ไม่ซ้ำ| H[Move หรือ Copy + Verify + Delete]
    H --> I[บันทึก Photo ใน Database]
    I --> J[แสดงใน Gallery พร้อมสถานะ pending]
    I --> K[เข้าคิว AI อัตโนมัติทันทีต่อรูป]
    K --> L[อ่านไฟล์จาก Destination และสร้าง Description]
    L --> M[บันทึก Description และ completed]
    M --> J
    J --> N[ค้นข้อความใน Description ด้วย SQLite]
```

## Feature Requirements

### Setup

- เลือก Source และ Destination ด้วย Native Directory Dialog
- ตรวจว่า Path มีอยู่จริงและอ่าน/เขียนได้
- จำค่าล่าสุดใน `settings`
- ป้องกันการกด Scan หรือ Move เมื่อยังตั้งค่าไม่ครบ

### Scan และ Selection

- Scan รูปใน Source
- ตรวจทั้ง Extension และเนื้อหาไฟล์
- แสดง Thumbnail, Filename และขนาด
- เลือกทีละรูป, Select All และ Clear Selection
- มีปุ่ม Refresh

### Move

- ย้ายเหมือน Cut/Paste
- ตรวจชื่อซ้ำใน Destination และ Database
- แสดง Progress และเวลารวม
- Error หนึ่งรูปต้องไม่หยุดทั้งชุด
- หลังย้ายสำเร็จให้ Refresh Gallery อัตโนมัติ

### Gallery และ Integrity

- Gallery ใช้ Database เป็นรายการหลักและตรวจไฟล์จริงบน Disk
- แสดงเฉพาะ Photo records ของ Destination ปัจจุบัน ไม่สแกนนำเข้ารูปเดิมที่ไม่มี record
- แสดง Description และ Analysis Status พร้อมปุ่ม Retry สำหรับงาน AI ที่ล้มเหลว
- Record มีแต่ไฟล์หาย ให้แสดง `missing`
- ลบผ่านแอป ให้ลบไฟล์จริง พร้อมลบ Record ใน `photos` และ `photo_tags`
- ไฟล์กลับมา ให้เปลี่ยนเป็น `active`

### AI

- หลังแต่ละรูปย้ายและบันทึก Photo สำเร็จ ให้ส่ง Photo ID เข้าคิว AI อัตโนมัติ ไม่รอทั้งชุดย้ายเสร็จ
- ใช้คิว AI แยกจาก Worker ย้ายรูป เพื่อไม่ให้การเรียก AI ขวางการย้ายรูปถัดไป
- วิเคราะห์เฉพาะรูปที่ย้ายผ่านโปรแกรมสำเร็จ รูปที่ skipped/failed และรูปเดิมใน Destination ไม่เข้าคิวใหม่
- สถานะ `pending`, `processing`, `completed`, `failed`
- บันทึก Description เป็นผลลัพธ์หลัก; Tags เป็นส่วนเสริม
- Retry ได้
- AI ล้มเหลวต้องไม่กระทบไฟล์ที่ย้ายสำเร็จ

### Search

- ค้น Description แบบข้อความบางส่วน
- ใช้ SQLite สำหรับค้นข้อความ ไม่ใช้ Embedding หรือ Vector Database ใน MVP
- ส่วนเสริม: ค้นจาก Tags, Description + Tags และ Tag Cloud

## Out of Scope สำหรับ MVP

- รองรับ MTP ของ Android/iPhone โดยตรง
- Sync Cloud Storage
- Face Recognition
- Video Analysis
- HEIC Conversion
- การแก้ไขรูปภาพ
- สแกน Destination เพื่อนำไฟล์เดิมที่ไม่มี record เข้า Database หรือเข้าคิว AI
- Vector Search

โทรศัพท์ใน MVP ต้องมองเห็นเป็นโฟลเดอร์หรือเลือกโฟลเดอร์ DCIM ที่ระบบเข้าถึงได้
