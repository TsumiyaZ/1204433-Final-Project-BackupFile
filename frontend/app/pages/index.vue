<template>
  <main class="relative isolate min-h-[calc(100dvh-4rem)] overflow-hidden">
    <div
      class="pointer-events-none absolute -left-24 top-8 size-80 rounded-full bg-secondary/15 blur-3xl"
    />
    <div
      class="pointer-events-none absolute -right-20 bottom-0 size-96 rounded-full bg-accent/12 blur-3xl"
    />

    <section
      class="relative mx-auto grid max-w-7xl gap-10 px-6 py-16 lg:grid-cols-[1.1fr_0.9fr] lg:items-center lg:py-24"
    >
      <div>
        <span
          class="inline-flex rounded-full border border-primary/20 bg-primary/8 px-4 py-2 text-sm font-semibold text-primary"
        >
          Desktop Photo Backup
        </span>

        <h1
          class="mt-6 max-w-3xl text-4xl font-bold leading-tight tracking-tight text-base-content sm:text-5xl lg:text-6xl"
        >
          ย้ายรูปจากอุปกรณ์ของคุณ
          <span class="text-primary">ง่ายและชัดเจน</span>
        </h1>

        <p class="mt-6 max-w-2xl text-lg leading-8 text-base-content/65">
          กำหนดโฟลเดอร์ต้นทางและปลายทาง สแกนดูรูปทั้งหมดก่อนเลือก
          แล้วจึงย้ายไปยังพื้นที่จัดเก็บที่ต้องการ
        </p>

        <div class="mt-8 flex flex-col gap-3 sm:flex-row">
          <NuxtLink to="/setup" class="btn btn-primary min-h-12 px-6">
            เริ่มตั้งค่าโฟลเดอร์
          </NuxtLink>
          <NuxtLink to="/scan" class="btn btn-outline min-h-12 px-6">
            ไปหน้าสแกนรูป
          </NuxtLink>
        </div>
      </div>

      <div
        class="rounded-[2rem] border border-base-300 bg-base-100 p-5 shadow-[0_24px_70px_rgba(91,64,48,0.12)]"
      >
        <div class="rounded-[1.5rem] bg-base-200 p-6 sm:p-8">
          <div class="mb-8 flex items-center justify-between">
            <div>
              <p class="text-sm font-semibold text-primary">ขั้นตอนการใช้งาน</p>
              <h2 class="mt-1 text-2xl font-bold">สำรองรูปใน 3 ขั้นตอน</h2>
            </div>
            <span class="badge badge-primary badge-outline h-8 px-3">พร้อมใช้งาน</span>
          </div>

          <ol class="space-y-4">
            <li
              v-for="(step, index) in steps"
              :key="step.title"
              class="flex gap-4 rounded-2xl border border-base-300 bg-base-100 p-4"
            >
              <span
                class="grid size-10 shrink-0 place-items-center rounded-xl bg-primary/10 font-bold text-primary"
              >
                {{ index + 1 }}
              </span>
              <div>
                <h3 class="font-bold">{{ step.title }}</h3>
                <p class="mt-1 text-sm leading-6 text-base-content/60">
                  {{ step.description }}
                </p>
              </div>
            </li>
          </ol>
        </div>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
const steps = [
  {
    title: "เลือกตำแหน่ง",
    description: "กำหนด Source และ Destination เพียงครั้งเดียว",
  },
  {
    title: "ตรวจดูรูป",
    description: "สแกน ดูตัวอย่าง และเลือกเฉพาะรูปที่ต้องการ",
  },
  {
    title: "ย้ายอย่างมั่นใจ",
    description: "ตรวจสอบผลลัพธ์ก่อนบันทึกข้อมูลลงฐานข้อมูล",
  },
];
</script>
