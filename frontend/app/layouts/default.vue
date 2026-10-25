<template>
  <div data-theme="photo-dark" class="min-h-dvh bg-base-200 text-base-content">
    <header
      class="sticky top-0 z-40 border-b border-base-300 bg-base-100/95 backdrop-blur"
    >
      <nav class="navbar mx-auto min-h-16 w-full max-w-[1500px] px-5 sm:px-8">
        <div class="navbar-start">
          <NuxtLink
            to="/gallery"
            class="flex min-h-12 items-center gap-3 rounded-xl px-2 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
            aria-label="ไปหน้า Gallery"
          >
            <span
              class="grid size-10 place-items-center rounded-xl bg-primary text-primary-content"
              aria-hidden="true"
            >
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.8"
                class="size-6"
              >
                <path d="M4 7.5A2.5 2.5 0 0 1 6.5 5h2l1.2 1.5h7.8A2.5 2.5 0 0 1 20 9v8.5a2.5 2.5 0 0 1-2.5 2.5h-11A2.5 2.5 0 0 1 4 17.5z" />
                <path d="m8 16 2.4-2.6a1 1 0 0 1 1.5 0l1.1 1.2 1.4-1.5a1 1 0 0 1 1.5 0L18 15.5" />
                <circle cx="15.5" cy="10" r="1.2" />
              </svg>
            </span>

            <span>
              <span class="block text-base font-bold leading-tight">Photo Backup</span>
              <span class="block text-xs text-base-content/55">จัดเก็บรูปอย่างเป็นระเบียบ</span>
            </span>
          </NuxtLink>
        </div>

        <div class="navbar-end gap-1 sm:gap-2">
          <NuxtLink
            v-for="item in navigation"
            :key="item.to"
            :to="item.to"
            class="btn min-h-11 border-0 px-3 sm:px-4"
            :class="isActive(item.to) ? 'btn-primary' : 'btn-ghost'"
          >
            {{ item.label }}
          </NuxtLink>
          <button
            type="button"
            class="btn btn-square btn-ghost min-h-11 min-w-11"
            aria-label="เปิดการตั้งค่าแอป"
            title="การตั้งค่า"
            @click="openSettings"
          >
            <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" class="size-5">
              <path d="M12 8.5a3.5 3.5 0 1 0 0 7 3.5 3.5 0 0 0 0-7Z" stroke="currentColor" stroke-width="1.8" />
              <path d="m19.4 13.5 1.1.9-1.2 2.1-1.4-.5a7.9 7.9 0 0 1-1.4.8l-.3 1.5h-2.4l-.4-1.5a7.9 7.9 0 0 1-1.5-.1l-1 1.2-2.1-1.2.4-1.5a7.8 7.8 0 0 1-1-1.2l-1.5.1-.7-2.3 1.3-.8a7.8 7.8 0 0 1 .1-1.6L6.2 8.7l1.2-2.1 1.5.4a7.8 7.8 0 0 1 1.2-.9L10 4.6l2.3-.7.8 1.3a7.8 7.8 0 0 1 1.6.1l.9-1.2 2.1 1.2-.4 1.5a7.8 7.8 0 0 1 .9 1.2l1.5-.1.7 2.3-1.3.8a7.8 7.8 0 0 1-.1 1.5Z" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round" />
            </svg>
          </button>
        </div>
      </nav>
    </header>

    <div>
      <slot />
    </div>

    <dialog ref="settingsDialog" class="settings-drawer" aria-labelledby="settings-title" @click="onSettingsBackdropClick">
      <aside class="settings-panel">
        <header class="flex items-start justify-between gap-4 border-b border-base-300 px-6 py-5">
          <div>
            <p class="text-xs font-semibold uppercase tracking-[0.18em] text-primary">Photo Backup</p>
            <h2 id="settings-title" class="mt-1 text-xl font-bold">การตั้งค่า</h2>
          </div>
          <button type="button" class="btn btn-square btn-ghost min-h-11 min-w-11" aria-label="ปิดการตั้งค่า" @click="settingsDialog?.close()">
            <svg viewBox="0 0 20 20" fill="none" aria-hidden="true" class="size-5"><path d="m5 5 10 10M15 5 5 15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" /></svg>
          </button>
        </header>

        <div class="space-y-6 p-6">
          <section aria-labelledby="settings-storage-title">
            <h3 id="settings-storage-title" class="text-sm font-semibold text-base-content/70">พื้นที่จัดเก็บ</h3>
            <div class="mt-3 space-y-3">
              <div class="rounded-xl border border-base-300 bg-base-200 p-3">
                <p class="text-xs text-base-content/60">Source · โฟลเดอร์ต้นทาง</p>
                <p class="mt-1 break-all text-sm">{{ settingPaths.source || "ยังไม่ได้ตั้งค่า" }}</p>
              </div>
              <div class="rounded-xl border border-base-300 bg-base-200 p-3">
                <p class="text-xs text-base-content/60">Destination · ที่เก็บรูป</p>
                <p class="mt-1 break-all text-sm">{{ settingPaths.destination || "ยังไม่ได้ตั้งค่า" }}</p>
              </div>
            </div>
            <NuxtLink to="/setup" class="btn btn-primary mt-4 min-h-11 w-full" @click="settingsDialog?.close()">
              แก้ไขโฟลเดอร์
            </NuxtLink>
          </section>

          <section class="border-t border-base-300 pt-5" aria-labelledby="settings-database-title">
            <h3 id="settings-database-title" class="text-sm font-semibold text-base-content/70">ฐานข้อมูล</h3>
            <button
              type="button"
              class="btn btn-outline mt-3 min-h-11 w-full"
              :disabled="databaseLocationPending"
              @click="showDatabaseInFolder"
            >
              <span v-if="databaseLocationPending" class="loading loading-spinner loading-xs" aria-hidden="true" />
              เปิดตำแหน่งไฟล์ Database
            </button>
            <p v-if="databaseLocationError" class="mt-2 break-words text-xs text-error" role="alert">{{ databaseLocationError }}</p>
          </section>

          <section class="border-t border-base-300 pt-5" aria-labelledby="settings-search-title">
            <h3 id="settings-search-title" class="text-sm font-semibold text-base-content/70">ระบบค้นหารูป</h3>
            <div class="mt-3 flex items-center gap-3 rounded-xl border border-base-300 bg-base-200 p-3">
              <span class="size-2.5 shrink-0 rounded-full" :class="modelStatus.state === 'ready' ? 'bg-success' : modelStatus.state === 'error' ? 'bg-error' : 'bg-warning'" aria-hidden="true" />
              <div class="min-w-0">
                <p class="text-sm font-medium">{{ modelStatus.state === "ready" ? "โมเดลพร้อมใช้งาน" : modelStatus.state === "error" ? "โมเดลมีปัญหา" : "กำลังเตรียมโมเดล" }}</p>
                <p class="mt-0.5 break-words text-xs leading-5 text-base-content/65">{{ modelStatus.message }}</p>
              </div>
            </div>
            <button
              type="button"
              class="btn btn-outline mt-3 min-h-11 w-full"
              :disabled="modelFolderPending"
              @click="showModelFolder"
            >
              <span v-if="modelFolderPending" class="loading loading-spinner loading-xs" aria-hidden="true" />
              เปิดตำแหน่งไฟล์โมเดล
            </button>
            <p v-if="modelFolderError" class="mt-2 break-words text-xs text-error" role="alert">{{ modelFolderError }}</p>
          </section>

          <p class="border-t border-base-300 pt-5 text-xs leading-5 text-base-content/55">ค้นหารูปด้วย EmbeddingGemma และสร้างคำบรรยายภาษาไทยด้วยโมเดลภาพในเครื่อง ผ่าน llama.cpp</p>
        </div>
      </aside>
    </dialog>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from "vue";
import { GetSetting, ShowDatabaseInFolder, ShowModelFolder } from "../../wailsjs/go/controller/SettingController";
import { GetEmbeddingStatus } from "../../wailsjs/go/main/App";

const route = useRoute();
const settingsDialog = ref<HTMLDialogElement | null>(null);
const settingPaths = ref({ source: "", destination: "" });
const modelStatus = ref({ state: "loading", message: "กำลังเชื่อมต่อกับโมเดล..." });
const databaseLocationPending = ref(false);
const databaseLocationError = ref("");
const modelFolderPending = ref(false);
const modelFolderError = ref("");
let settingsStatusTimer: ReturnType<typeof setInterval> | undefined;

const navigation = [
  { label: "ตั้งค่า", to: "/setup" },
  { label: "สแกนรูป", to: "/scan" },
  { label: "ย้ายรูป", to: "/move" },
  { label: "รูปภาพ", to: "/gallery"},
];

function isActive(path: string) {
  return route.path === path;
}

async function refreshSettingsData() {
  try {
    const setting = await GetSetting();
    settingPaths.value = { source: setting.Source ?? "", destination: setting.Dest ?? "" };
  } catch {
    settingPaths.value = { source: "", destination: "" };
  }
  try {
    modelStatus.value = await GetEmbeddingStatus();
  } catch {
    modelStatus.value = { state: "error", message: "ไม่สามารถอ่านสถานะโมเดลได้" };
  }
}

function openSettings() {
  settingsDialog.value?.showModal();
  void refreshSettingsData();
}

async function showDatabaseInFolder() {
  if (databaseLocationPending.value) return;
  databaseLocationPending.value = true;
  databaseLocationError.value = "";
  try {
    await ShowDatabaseInFolder();
  } catch (error) {
    databaseLocationError.value = error instanceof Error ? error.message : String(error);
  } finally {
    databaseLocationPending.value = false;
  }
}

async function showModelFolder() {
  if (modelFolderPending.value) return;
  modelFolderPending.value = true;
  modelFolderError.value = "";
  try {
    await ShowModelFolder();
  } catch (error) {
    modelFolderError.value = error instanceof Error ? error.message : String(error);
  } finally {
    modelFolderPending.value = false;
  }
}

function onSettingsBackdropClick(event: MouseEvent) {
  if (event.target === settingsDialog.value) settingsDialog.value?.close();
}

onMounted(() => {
  settingsStatusTimer = setInterval(() => {
    if (settingsDialog.value?.open) void refreshSettingsData();
  }, 1500);
});

onBeforeUnmount(() => {
  if (settingsStatusTimer) clearInterval(settingsStatusTimer);
});
</script>

<style scoped>
.settings-drawer {
  position: fixed;
  inset: 0;
  width: 100%;
  height: 100%;
  max-width: none;
  max-height: none;
  margin: 0;
  padding: 0;
  border: 0;
  background: transparent;
  color: inherit;
}

.settings-drawer[open] {
  display: flex;
  justify-content: flex-end;
}

.settings-drawer::backdrop {
  background: rgb(0 0 0 / 58%);
}

.settings-panel {
  width: min(390px, 100%);
  height: 100%;
  overflow-y: auto;
  border-left: 1px solid var(--color-base-300);
  background: var(--color-base-100);
  box-shadow: -24px 0 80px rgb(0 0 0 / 35%);
  animation: settings-enter 180ms ease-out both;
}

@keyframes settings-enter {
  from { transform: translateX(18px); opacity: 0.7; }
  to { transform: translateX(0); opacity: 1; }
}

@media (prefers-reduced-motion: reduce) {
  .settings-panel { animation: none; }
}
</style>
