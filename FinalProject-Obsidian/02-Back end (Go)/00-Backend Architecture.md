---
title: Backend Architecture
tags:
  - backend
  - golang
status: planning
---

# Backend Architecture

กลับไปที่ [[01-Cross Platform App/00-Project Home|Project Home]]

## Layers

```text
Nuxt Frontend
      ↕
Wails Controller
      ↓
Services
 ├─ Backup
 ├─ Analyzer
 └─ Search
      ↓
Repository / Queries
      ↓
SQLite / GORM
```

## Folder Responsibilities

| Folder | Responsibility |
|---|---|
| `controller/` | Wails-facing methods, validation และ DTO mapping |
| `model/` | GORM Models และ DTOs |
| `model/query/` | Queries ที่ใช้ซ้ำได้ |
| `repository/` | DB connection, AutoMigrate และ Transaction |
| `service/backup/` | Scan, Move, Duplicate และ Integrity |
| `service/analyzer/` | AI analysis และ Retry |
| `service/search/` | Description, Tag และ Tag Cloud queries |

## Design Rules

- Controller ต้องบางและไม่ใส่ File Operation ขนาดใหญ่
- Service ไม่ควรอ้างอิง Nuxt
- Database access ผ่าน Repository/Query layer
- AI Provider ต้องเปลี่ยนได้ผ่าน Interface
- Error ส่งกลับเป็นข้อมูลที่ Frontend แสดงได้
- Long-running tasks ต้องส่ง Progress Event
- ใช้ `context.Context` สำหรับ Cancel และ Timeout

## Related Notes

- [[01-Backup Engine]]
- [[02-Database and GORM]]
- [[03-AI and Search]]
- [[04-Wails-Nuxt-IPC/00-IPC Contract|IPC Contract]]
