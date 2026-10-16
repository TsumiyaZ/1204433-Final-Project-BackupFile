---
title: Photo Transfer and AI Analysis
aliases:
  - Final Project Home
tags:
  - final-project
  - project-home
status: planning
updated: 2026-10-16
---

# Final Project: Photo Transfer and AI Analysis

> [!summary] Project Summary
> แอป Desktop ด้วย Wails + Go + Nuxt สำหรับย้ายรูปที่ผู้ใช้เลือกจาก Source ไป Destination เมื่อแต่ละรูปย้ายสำเร็จและบันทึก Photo แล้ว ให้เข้าคิว AI อัตโนมัติเพื่อสร้าง Description และค้นด้วยข้อความใน SQLite ส่วน Tags/Tag Cloud เป็นส่วนเสริม

## Project Navigation

### 01 — Cross Platform App

- [[01-Requirements and Flow]]
- [[02-Team Division]]
- [[03-Milestones]]

### 02 — Back end (Go)

- [[02-Back end (Go)/00-Backend Architecture|Backend Architecture]]
- [[02-Back end (Go)/01-Backup Engine|Backup Engine]]
- [[02-Back end (Go)/02-Database and GORM|Database and GORM]]
- [[02-Back end (Go)/03-AI and Search|AI and Search]]

### 03 — Front End

- [[03-Front End/00-Frontend Overview|Frontend Overview]]
- [[03-Front End/01-Screens and UX|Screens and UX]]

### 04 — Wails-Nuxt-IPC

- [[04-Wails-Nuxt-IPC/00-IPC Contract|IPC Contract]]
- [[04-Wails-Nuxt-IPC/01-Events and Progress|Events and Progress]]

### 05 — Wails Application

- [[05-Wails Application/00-Git Workflow|Git Workflow]]
- [[05-Wails Application/01-Testing Checklist|Testing Checklist]]
- [[05-Wails Application/02-Demo Plan|Demo Plan]]
- [[05-Wails Application/03-Build and Security|Build and Security]]

## Current Scope

- [ ] เลือก Source และ Destination
- [ ] Scan เฉพาะรูป JPG/JPEG, PNG และ WebP
- [ ] เลือกรูปบางส่วนหรือ Select All
- [ ] ย้ายไฟล์แบบ Cut/Paste อย่างปลอดภัย
- [ ] บันทึกข้อมูลด้วย SQLite และ GORM
- [ ] แสดง Gallery จาก Photo records ของ Destination ปัจจุบัน
- [ ] สร้าง Description ด้วย AI อัตโนมัติหลังย้ายแต่ละรูปสำเร็จ
- [ ] ค้นจาก Description ด้วยข้อความใน SQLite
- [ ] ส่วนเสริม: Tags, Tag Search และ Tag Cloud
- [ ] ตรวจ Integrity ระหว่าง Database กับไฟล์จริง
- [ ] ทำงานบน Windows และ macOS

## System Overview

```mermaid
flowchart LR
    U[User] --> F[Nuxt Frontend]
    F <--> I[Wails IPC]
    I --> B[Backup Service]
    I --> A[AI Service]
    I --> S[Search Service]
    B --> D[(SQLite / GORM)]
    A --> D
    S --> D
    B --> FS[Source / Destination]
    A --> FS
```

## Definition of Success

โปรเจกต์พร้อมส่งเมื่อผู้ใช้สามารถทำ Flow นี้ได้โดยไม่ใช้ Terminal:

```text
เลือก Source/Destination
→ Scan และเลือกรูป
→ Move
→ แต่ละรูปย้ายสำเร็จและบันทึก DB → เข้าคิว AI อัตโนมัติ
→ ดู Gallery พร้อมสถานะและ Description
→ ค้นจาก Description
```

AI อ่านไฟล์จาก Destination และทำงานขนานกับการย้ายรูปถัดไป รูปเดิมที่อยู่ใน Destination ไม่ถูกนำเข้า DB หรือวิเคราะห์อัตโนมัติ ไม่ใช้ Vector Search ใน MVP

และต้องไม่มีกรณีลบไฟล์ต้นฉบับเมื่อ Copy หรือการตรวจสอบปลายทางล้มเหลว
