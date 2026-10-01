---
title: Frontend Overview
tags:
  - frontend
  - nuxt
status: planning
---

# Frontend Overview

กลับไปที่ [[01-Cross Platform App/00-Project Home|Project Home]]

## Main Navigation

- Setup
- Move Photos
- Gallery & AI
- Search
- History / Integrity

## Shared Components

- App Sidebar/Header
- Path Selector
- Photo Card
- Photo Grid
- Selection Toolbar
- Progress Bar
- Result Summary
- Status Badge
- Error Alert
- Confirm Dialog
- Empty State
- Tag Chip/Tag Cloud

## State Rules

- ไม่มี Source → Disable Scan
- ไม่มี Destination → Disable Move/Gallery
- ไม่มี Selection → Disable Move/Analyze
- กำลัง Move → Disable Path Change และ Move ซ้ำ
- กำลัง Analyze → แสดง Progress ต่อรูป
- Missing File → แสดง Badge และไม่ส่งให้ AI
- Failed Analysis → แสดง Retry

## Development Strategy

1. ตกลง DTO กับทีม Backend
2. สร้าง UI ด้วย Mock Data
3. เชื่อม Wails Methods
4. เชื่อม Events
5. ทดสอบ Error/Empty/Loading States

อ่านต่อ: [[01-Screens and UX]] และ [[04-Wails-Nuxt-IPC/00-IPC Contract|IPC Contract]]
