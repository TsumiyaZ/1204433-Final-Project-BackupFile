# Photo Backup

แอปเดสก์ท็อปสำหรับจัดเก็บและค้นหารูปภาพ ใช้ EmbeddingGemma กับ llama.cpp เพื่อค้นหารูปจากข้อความ โดยประมวลผลบนเครื่องของคุณ

## ดาวน์โหลด

เลือกดาวน์โหลดตามระบบปฏิบัติการ:

[![ดาวน์โหลดสำหรับ Windows](https://img.shields.io/badge/ดาวน์โหลด-Windows-0078D4?logo=windows&logoColor=white)](https://github.com/TsumiyaZ/1204433-Final-Project-BackupFile/releases/latest/download/FinalProject.exe)

[![ดาวน์โหลดสำหรับ macOS](https://img.shields.io/badge/ดาวน์โหลด-macOS-555555?logo=apple&logoColor=white)](https://github.com/TsumiyaZ/1204433-Final-Project-BackupFile/releases/latest/download/PhotoBackup-macOS-universal.zip)

แอป macOS รองรับทั้ง Apple Silicon และ Intel

## เริ่มต้นใช้งาน

1. ดาวน์โหลดไฟล์สำหรับระบบปฏิบัติการของคุณจากด้านบน
2. เปิดแอปและตั้งค่าโฟลเดอร์ต้นทางกับโฟลเดอร์ปลายทาง
3. เมื่อเปิดครั้งแรก แอปจะดาวน์โหลด llama.cpp และโมเดลสำหรับค้นหารูป ต้องเชื่อมต่ออินเทอร์เน็ตและมีพื้นที่ว่างในเครื่อง
4. เมื่อติดตั้งโมเดลเสร็จแล้ว การสร้าง embedding และค้นหารูปจะทำงานในเครื่อง

บน macOS อาจมีข้อความเตือนความปลอดภัยเพราะแอปยังไม่ได้เซ็นรับรองโดย Apple

## ความสามารถ

- สแกนและย้ายรูปภาพไปยังโฟลเดอร์จัดเก็บ
- ตรวจหารูปซ้ำด้วย checksum
- สร้าง embedding เบื้องหลัง พร้อมแสดงความคืบหน้า
- ค้นหารูปด้วยข้อความและเลือกดูตามแท็ก
- เก็บฐานข้อมูลและโมเดลไว้ในเครื่อง

## สำหรับนักพัฒนา

ต้องติดตั้ง Go ตามเวอร์ชันใน `go.mod`, Node.js และ npm รวมถึง Wails CLI v2.16.0

```bash
wails dev
wails build
```

GitHub Actions จะ build Windows และ macOS เมื่อมีการ push tag เวอร์ชันที่ขึ้นต้นด้วย `v` เช่น `v1.0.0` แล้วสร้าง GitHub Release พร้อมไฟล์ดาวน์โหลดให้อัตโนมัติ
