<template>
  <main class="mx-auto max-w-6xl px-6 py-10">
    <div class="flex justify-between">
      <h1 class="mb-6 text-3xl font-bold">Gallery</h1>
      <button
        type="button"
        class="btn btn-primary"
        :disabled="loading"
        @click="refreshPhoto"
      >
        {{ loading ? "กำลังโหลด..." : "รีเฟรชรูปภาพ" }}
      </button>
    </div>

    <p v-if="loading">กำลังโหลดรูปภาพ...</p>

    <p v-else-if="errMessage" role="alert">
      {{ errMessage }}
    </p>

    <p v-else-if="photos.length === 0">
      ยังไม่มีรูปที่ย้ายมาใน Destination ปัจจุบัน
    </p>

    <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <article
        v-for="photo in photos"
        :key="photo.id"
        class="rounded-2xl border border-base-300 bg-base-100 p-5"
      >
        <img
          v-if="previewUrls[photo.id]"
          :src="previewUrls[photo.id]"
          :alt="photo.filename"
          class="mb-4 h-48 w-full rounded-xl object-cover"
        />
        <p v-else class="mb-4 text-base-content/60">
          {{
            photo.status === "missing"
              ? "ไม่พบไฟล์รูปภาพ"
              : "ไม่สามารถแสดงตัวอย่างได้"
          }}
        </p>

        <h2 class="break-all font-bold">{{ photo.filename }}</h2>

        <p class="mt-2">สถานะไฟล์: {{ photo.status }}</p>
        <p>สถานะ AI: {{ photo.analysisStatus }}</p>

        <p class="mt-3 text-base-content/70">
          {{ photo.description || "ยังไม่มีคำอธิบาย" }}
        </p>
      </article>
    </div>
  </main>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  GetGalleryPhotos,
  GetGalleryPhotoPreview,
} from "../../wailsjs/go/controller/GalleryController";
import type { dto } from "../../wailsjs/go/models";

const photos = ref<dto.GalleryPhoto[]>([]);
const loading = ref(true);
const errMessage = ref("");
const previewUrls = ref<Record<number, string>>({});

async function loadPhotos() {
  loading.value = true;
  errMessage.value = "";

  try {
    photos.value = await GetGalleryPhotos();

    previewUrls.value = {};
    await loadPreviews(photos.value);
  } catch (err) {
    errMessage.value = err instanceof Error ? err.message : String(err);
  } finally {
    loading.value = false;
  }
}

async function loadPreviews(items: dto.GalleryPhoto[]) {
  for (const photo of items) {
    if (photo.status !== "active") continue;

    try {
      previewUrls.value[photo.id] = await GetGalleryPhotoPreview(photo.id);
    } catch (err) {
      previewUrls.value[photo.id] = "";

      const message = err instanceof Error ? err.message : String(err);
      errMessage.value = `${photo.filename}: ${message}`;
      console.error("โหลด Preview ไม่สำเร็จ", photo.id, err);
    }
  }
}

function refreshPhoto() {
  loadPhotos();
}

onMounted(loadPhotos);
</script>

<style scoped></style>
