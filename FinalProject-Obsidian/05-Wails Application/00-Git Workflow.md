---
title: Git Workflow
tags:
  - git
  - teamwork
status: active
---

# Git Workflow

กลับไปที่ [[01-Cross Platform App/00-Project Home|Project Home]]

## Branches

- `main` — Version ที่ทดสอบและพร้อม Demo
- `develop` — รวมงานระหว่างพัฒนา
- `feature/database`
- `feature/backup-engine`
- `feature/frontend`
- `feature/ai-search`

## Working Rules

1. Pull `develop` ล่าสุดก่อนเริ่มงาน
2. หนึ่ง Feature ต่อหนึ่ง Branch
3. Commit เล็กและข้อความชัดเจน
4. Push Branch ของตัวเอง
5. เปิด Pull Request เข้า `develop`
6. ให้เพื่อน Review อย่างน้อยหนึ่งคน
7. Merge เมื่อ Test ผ่าน
8. Merge `develop` เข้า `main` เมื่อจบ Milestone
9. ห้าม Force Push เข้า `main`

## Files That Must Not Be Committed

- `.env`
- API Keys/Tokens
- `node_modules/`
- `.nuxt/`
- `.output/`
- `frontend/dist/`
- `build/bin/`
- Database ที่มีข้อมูลส่วนตัว

## Before Commit

- [ ] ตรวจ `git status`
- [ ] ไม่มี Secret
- [ ] Format Code
- [ ] Test Module
- [ ] Build ส่วนที่เปลี่ยน
- [ ] Commit เฉพาะงานของ Feature

## Before Merge

- [ ] Pull/Rebase `develop`
- [ ] แก้ Conflict และทดสอบใหม่
- [ ] Contract ไม่เปลี่ยนโดยไม่อัปเดตเอกสาร
- [ ] Reviewer อนุมัติ
