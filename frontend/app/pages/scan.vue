<template>
  <main class="min-h-[calc(100dvh-4rem)] bg-base-200/70 px-4 py-10 sm:px-6">
    <section class="mx-auto max-w-6xl">
      <header
        class="mb-8 flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between"
      >
        <div>
          <span
            class="inline-flex rounded-full border border-primary/20 bg-primary/8 px-4 py-2 text-sm font-semibold text-primary"
          >
            Photo scanner
          </span>

          <h1 class="mt-5 text-3xl font-bold tracking-tight sm:text-4xl">
            สแกนรูปภาพจาก Source
          </h1>

          <p class="mt-3 max-w-3xl break-all leading-7 text-base-content/60">
            {{ source || "ยังไม่ได้กำหนด Source folder" }}
          </p>
        </div>

        <div class="flex flex-wrap gap-3">
          <NuxtLink to="/setup" class="btn btn-outline min-h-12 border-primary/35 text-primary">
            ตั้งค่าโฟลเดอร์
          </NuxtLink>

          <button
            type="button"
            class="btn btn-primary min-h-12 min-w-36"
            :disabled="!canScan"
            @click="scan"
          >
            <span v-if="scanning" class="loading loading-spinner loading-sm" />

            {{ scanning ? "กำลังสแกน..." : "สแกนรูปภาพ" }}
          </button>
        </div>
      </header>

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
        v-else-if="!source"
        role="status"
        class="rounded-[1.5rem] border border-base-300 bg-base-100 p-10 text-center shadow-[0_16px_45px_rgba(91,64,48,0.08)]"
      >
        <h2 class="text-xl font-semibold">ยังไม่มี Source folder</h2>
        <p class="mt-2 text-base-content/70">
          กรุณาไปที่หน้าตั้งค่าและเลือกโฟลเดอร์ต้นทางก่อน
        </p>
      </div>

      <div
        v-else-if="hasScanned && photos.length === 0"
        role="status"
        class="rounded-[1.5rem] border border-base-300 bg-base-100 p-10 text-center shadow-[0_16px_45px_rgba(91,64,48,0.08)]"
      >
        <h2 class="text-xl font-semibold">ไม่พบรูปภาพ</h2>
        <p class="mt-2 text-base-content/70">
          รองรับไฟล์ JPG, JPEG, PNG, WEBP และ GIF
        </p>
      </div>

      <div v-else-if="photos.length > 0" class="space-y-5">
        <div
          class="flex flex-col gap-4 rounded-[1.5rem] border border-base-300 bg-base-100 p-5 shadow-[0_12px_36px_rgba(91,64,48,0.07)] sm:flex-row sm:items-center sm:justify-between"
        >
          <div>
            <p class="text-lg font-bold">พบ {{ photos.length }} รูป</p>

            <p class="text-sm text-base-content/70">
              เลือกแล้ว {{ selectedPaths.length }} รูป
            </p>
          </div>

          <button
            type="button"
            class="btn btn-outline min-h-12 border-primary/35 text-primary"
            @click="toggleSelectAll"
          >
            {{ allSelected ? "ยกเลิกทั้งหมด" : "เลือกทั้งหมด" }}
          </button>
        </div>

        <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
          <label
            v-for="photo in photos"
            :key="photo.path"
            class="card cursor-pointer overflow-hidden border bg-base-100 shadow-[0_12px_30px_rgba(91,64,48,0.06)] transition duration-200"
            :class="
              selectedPaths.includes(photo.path)
                ? 'border-primary ring-2 ring-primary/20'
                : 'border-base-300 hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-[0_16px_36px_rgba(91,64,48,0.10)]'
            "
          >
            <div
              class="flex aspect-video items-center justify-center overflow-hidden bg-base-200"
            >
              <img
                v-if="previewUrls[photo.path]"
                :src="previewUrls[photo.path]"
                :alt="photo.filename"
                class="h-full w-full object-cover"
                loading="lazy"
                decoding="async"
              />

              <span
                v-else-if="!previewFinished[photo.path]"
                class="loading loading-spinner loading-md text-primary"
              />

              <span
                v-else
                class="text-lg font-bold uppercase text-base-content/60"
              >
                {{ photo.extension.replace(".", "") }}
              </span>
            </div>

            <div class="card-body gap-3 border-t border-base-300/70 p-4">
              <div class="flex items-start gap-3">
                <input
                  v-model="selectedPaths"
                  type="checkbox"
                  class="checkbox checkbox-primary mt-1"
                  :value="photo.path"
                />

                <div class="min-w-0">
                  <h2 class="truncate font-semibold" :title="photo.filename">
                    {{ photo.filename }}
                  </h2>

                  <p
                    class="mt-1 break-all text-sm text-base-content/60"
                    :title="photo.relativePath"
                  >
                    {{ photo.relativePath }}
                  </p>

                  <p class="mt-2 text-sm text-base-content/70">
                    {{ formatFileSize(photo.size) }}
                  </p>
                </div>
              </div>
            </div>
          </label>
        </div>
      </div>

      <div
        v-else
        role="status"
        class="rounded-[1.5rem] border-2 border-dashed border-primary/20 bg-base-100/60 p-12 text-center"
      >
        <h2 class="text-xl font-semibold">พร้อมสแกนรูปภาพ</h2>
        <p class="mt-2 text-base-content/70">
          กดปุ่มสแกนเพื่อค้นหารูปทั้งหมดใน Source และโฟลเดอร์ย่อย
        </p>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import {
  GetPhotoPreview,
  ScanPhotos,
} from "../../wailsjs/go/controller/PhotoController";
import { GetSetting } from "../../wailsjs/go/controller/SettingController";
import type { dto } from "../../wailsjs/go/models";

const source = ref("");
const photos = ref<dto.ScannedPhoto[]>([]);
const selectedPaths = ref<string[]>([]);

const loadingSetting = ref(true);
const scanning = ref(false);
const hasScanned = ref(false);
const errorMessage = ref("");

const previewUrls = ref<Record<string, string>>({});
const previewFinished = ref<Record<string, boolean>>({});

const canScan = computed(() => {
  return source.value.trim() !== "" && !loadingSetting.value && !scanning.value;
});

const allSelected = computed(() => {
  return (
    photos.value.length > 0 &&
    selectedPaths.value.length === photos.value.length
  );
});

function getErrorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }

  return String(error);
}

async function loadSetting() {
  loadingSetting.value = true;

  try {
    const setting = await GetSetting();
    source.value = setting.Source ?? "";
  } catch (error) {
    errorMessage.value = getErrorMessage(error);
  } finally {
    loadingSetting.value = false;
  }
}

async function scan() {
  if (!canScan.value) {
    return;
  }

  scanning.value = true;
  hasScanned.value = false;
  errorMessage.value = "";
  selectedPaths.value = [];

  try {
    const scannedPhotos = await ScanPhotos(source.value);

    photos.value = scannedPhotos;
    previewUrls.value = {};
    previewFinished.value = {};
    hasScanned.value = true;

    void loadPreviews(scannedPhotos);
  } catch (error) {
    photos.value = [];
    errorMessage.value = getErrorMessage(error);
  } finally {
    scanning.value = false;
  }
}

function toggleSelectAll() {
  if (allSelected.value) {
    selectedPaths.value = [];
    return;
  }

  selectedPaths.value = photos.value.map((photo) => photo.path);
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) {
    return `${bytes} B`;
  }

  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }

  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

async function loadPreviews(items: dto.ScannedPhoto[]) {
  let nextIndex = 0;

  async function worker() {
    while (nextIndex < items.length) {
      const photo = items[nextIndex];
      nextIndex += 1;

      if (!photo) {
        return;
      }

      try {
        previewUrls.value[photo.path] = await GetPhotoPreview(
          source.value,
          photo.path,
        );
      } catch {
        previewUrls.value[photo.path] = "";
      } finally {
        previewFinished.value[photo.path] = true;
      }
    }
  }

  const workerCount = Math.min(4, items.length);

  await Promise.all(Array.from({ length: workerCount }, () => worker()));
}

onMounted(loadSetting);
</script>
