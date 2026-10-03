<template>
  <main class="min-h-[calc(100dvh-4rem)] bg-base-200/70 px-4 py-10 sm:px-6">
    <section class="mx-auto max-w-6xl">
      <header class="mb-8">
        <span
          class="inline-flex rounded-full border border-primary/20 bg-primary/8 px-4 py-2 text-sm font-semibold text-primary"
        >
          Photo scanner
        </span>

        <h1 class="mt-5 text-3xl font-bold tracking-tight sm:text-4xl">
          สแกน Source และ Destination
        </h1>

        <p class="mt-3 max-w-3xl text-base leading-7 text-base-content/65">
          ตรวจดูรูปจากโฟลเดอร์ต้นทางก่อนย้าย
          หรือสแกนโฟลเดอร์ปลายทางเพื่อดูรูปที่จัดเก็บไว้
        </p>
      </header>

      <ScanModeSelector
        :model-value="activeMode"
        :options="scanOptions"
        :disabled="scanningMode !== null"
        class="mb-6"
        @update:model-value="setActiveMode"
      />

      <ScanToolbar
        :scan-mode="activeMode"
        :current-path="currentPath"
        :can-scan="canScan"
        :is-scanning="isCurrentScanning"
        class="mb-6"
        @scan="scan"
      />

      <div v-if="errorMessage" role="alert" class="alert alert-error mb-6">
        {{ errorMessage }}
      </div>

      <div
        v-if="loadingSetting"
        class="flex min-h-64 items-center justify-center"
      >
        <span class="loading loading-spinner loading-lg text-primary" />
      </div>

      <div
        v-else-if="!currentPath"
        role="status"
        class="rounded-[1.5rem] border border-base-300 bg-base-100 p-10 text-center"
      >
        <h2 class="text-xl font-semibold">ยังไม่มี {{ currentFolderLabel }}</h2>
        <p class="mt-2 text-base-content/70">
          ไปที่หน้าตั้งค่าแล้วเลือกโฟลเดอร์ก่อนเริ่มสแกน
        </p>
      </div>

      <div
        v-else-if="currentHasScanned && currentPhotos.length === 0"
        role="status"
        class="rounded-[1.5rem] border border-base-300 bg-base-100 p-10 text-center"
      >
        <h2 class="text-xl font-semibold">
          ไม่พบรูปภาพใน {{ currentFolderLabel }}
        </h2>
        <p class="mt-2 text-base-content/70">
          รองรับไฟล์ JPG, JPEG, PNG, WEBP และ GIF
        </p>
      </div>

      <div v-else-if="currentPhotos.length > 0" class="space-y-5">
        <div
          class="flex flex-col gap-4 rounded-[1.5rem] border border-base-300 bg-base-100 p-5 sm:flex-row sm:items-center sm:justify-between"
        >
          <div>
            <p class="text-lg font-bold">
              พบ {{ currentPhotos.length }} รูปใน {{ currentFolderLabel }}
            </p>

            <p class="text-sm text-base-content/70">
              <template v-if="activeMode === 'source'">
                เลือกแล้ว {{ selectedPaths.length }} รูป
                สำหรับส่งต่อไปขั้นตอนย้าย
              </template>
              <template v-else>
                แสดงรูปที่อยู่ในโฟลเดอร์ปลายทางปัจจุบัน
              </template>
            </p>
          </div>

          <button
            v-if="activeMode === 'source'"
            type="button"
            class="btn btn-outline min-h-12 border-primary/35 text-primary"
            @click="toggleSelectAll"
          >
            {{ allSelected ? "ยกเลิกทั้งหมด" : "เลือกทั้งหมด" }}
          </button>
        </div>

        <ScanPhotoGrid
          v-model:selected-paths="selectedPaths"
          :photos="currentPhotos"
          :scan-mode="activeMode"
          :preview-urls="previewUrls"
          :preview-finished="previewFinished"
        />
      </div>

      <div
        v-else
        role="status"
        class="rounded-[1.5rem] border-2 border-dashed border-primary/20 bg-base-100/60 p-12 text-center"
      >
        <h2 class="text-xl font-semibold">
          พร้อมสแกน {{ currentFolderLabel }}
        </h2>
        <p class="mt-2 text-base-content/70">
          กดปุ่มสแกนเพื่อค้นหารูปทั้งหมดในโฟลเดอร์และโฟลเดอร์ย่อย
        </p>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import ScanModeSelector from "../components/scan/ModeSelector.vue";
import ScanPhotoGrid from "../components/scan/PhotoGrid.vue";
import ScanToolbar from "../components/scan/ScanToolbar.vue";
import { usePhotoScanner } from "../composables/usePhotoScanner";

const {
  scanOptions,
  activeMode,
  selectedPaths,
  loadingSetting,
  scanningMode,
  errorMessage,
  previewUrls,
  previewFinished,
  currentPath,
  currentFolderLabel,
  currentPhotos,
  currentHasScanned,
  isCurrentScanning,
  canScan,
  allSelected,
  setActiveMode,
  scan,
  toggleSelectAll,
} = usePhotoScanner();
</script>
