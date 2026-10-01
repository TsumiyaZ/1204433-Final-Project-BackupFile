---
title: Events and Progress
tags:
  - wails
  - events
status: draft
---

# Events and Progress

กลับไปที่ [[00-IPC Contract]]

## Events ที่ควรมี

### Move Progress

- Event name: `move:progress`
- Current index
- Total
- Filename
- Status
- Percent

### Move Completed

- Event name: `move:completed`
- Success/Skipped/Failed counts
- Duration

### Analysis Progress

- Event name: `analysis:progress`
- Photo ID
- Filename
- Status
- Current/Total

### Integrity Result

- Event name: `integrity:completed`
- Active count
- Missing count

## UI Behaviour

- Subscribe เมื่อ Page Mount
- Unsubscribe เมื่อ Page Unmount
- ป้องกัน Event Listener ซ้ำ
- Update เฉพาะ Card ที่เกี่ยวข้อง
- Disable ปุ่มที่ทำงานซ้ำไม่ได้
- Operation จบแล้ว Refresh จาก Backend อีกครั้ง

## Cancellation

ถ้ามีเวลาเพิ่ม:

- Cancel Scan
- Cancel Pending Move ก่อนเริ่มไฟล์ถัดไป
- Cancel AI Analysis

การ Cancel ห้ามทิ้ง Temporary File หรือสถานะ `processing` ค้าง
