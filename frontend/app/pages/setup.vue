<template>
  <main class="min-h-[calc(100dvh-4rem)] bg-base-200/70 px-4 py-10 sm:px-6">
    <section class="mx-auto max-w-4xl">
      <header class="mb-8">
        <span
          class="inline-flex rounded-full border border-primary/20 bg-primary/8 px-4 py-2 text-sm font-semibold text-primary"
        >
          ตั้งค่าการสำรองรูป
        </span>

        <h1 class="mt-5 text-3xl font-bold tracking-tight sm:text-4xl">
          กำหนดโฟลเดอร์ต้นทางและปลายทาง
        </h1>

        <p class="mt-3 max-w-2xl text-base leading-7 text-base-content/65">
          เลือกตำแหน่งที่ต้องการค้นหารูป และโฟลเดอร์สำหรับจัดเก็บไฟล์
          โปรแกรมจะจดจำค่านี้ไว้สำหรับการใช้งานครั้งถัดไป
        </p>
      </header>

      <form
        class="overflow-hidden rounded-[1.75rem] border border-base-300 bg-base-100 shadow-[0_20px_60px_rgba(91,64,48,0.10)]"
        @submit.prevent="saveSetting"
      >
        <div class="border-b border-base-300 bg-primary/5 px-6 py-5 sm:px-8">
          <h2 class="text-lg font-bold">ตำแหน่งจัดเก็บ</h2>
          <p class="mt-1 text-sm text-base-content/60">
            Source และ Destination ต้องเป็นคนละโฟลเดอร์
          </p>
        </div>

        <div class="space-y-8 p-6 sm:p-8">
          <div>
            <div class="mb-3 flex items-center justify-between gap-4">
              <label for="source" class="font-bold">Source folder</label>
              <span class="badge badge-outline border-primary/30 text-primary">ต้นทาง</span>
            </div>

            <div class="flex flex-col gap-3 sm:flex-row">
              <input
                id="source"
                v-model="source"
                class="input input-bordered min-h-12 w-full bg-base-100 focus:border-primary focus:outline-primary"
                placeholder="ยังไม่ได้เลือกโฟลเดอร์ต้นทาง"
                readonly
              />

              <button
                type="button"
                class="btn btn-outline min-h-12 shrink-0 border-primary/40 text-primary hover:border-primary hover:bg-primary hover:text-primary-content"
                :disabled="selecting !== null || loading"
                @click="chooseSource"
              >
                <span
                  v-if="selecting === 'source'"
                  class="loading loading-spinner loading-sm"
                />
                เลือก Source
              </button>
            </div>

            <p class="mt-2 text-sm leading-6 text-base-content/55">
              โปรแกรมจะค้นหา JPG, JPEG, PNG, WEBP และ GIF จากโฟลเดอร์นี้
            </p>
          </div>

          <div class="flex justify-center" aria-hidden="true">
            <div class="h-px w-full bg-base-300" />
          </div>

          <div>
            <div class="mb-3 flex items-center justify-between gap-4">
              <label for="destination" class="font-bold">Destination folder</label>
              <span class="badge badge-outline border-accent/40 text-accent">ปลายทาง</span>
            </div>

            <div class="flex flex-col gap-3 sm:flex-row">
              <input
                id="destination"
                v-model="destination"
                class="input input-bordered min-h-12 w-full bg-base-100 focus:border-primary focus:outline-primary"
                placeholder="ยังไม่ได้เลือกโฟลเดอร์ปลายทาง"
                readonly
              />

              <button
                type="button"
                class="btn btn-outline min-h-12 shrink-0 border-primary/40 text-primary hover:border-primary hover:bg-primary hover:text-primary-content"
                :disabled="selecting !== null || loading"
                @click="chooseDestination"
              >
                <span
                  v-if="selecting === 'destination'"
                  class="loading loading-spinner loading-sm"
                />
                เลือก Destination
              </button>
            </div>

            <p class="mt-2 text-sm leading-6 text-base-content/55">
              รูปที่เลือกจะถูกย้ายมาเก็บในโฟลเดอร์นี้
            </p>
          </div>

          <div v-if="sameDirectory" role="alert" class="alert alert-warning">
            Source และ Destination ต้องเป็นคนละโฟลเดอร์
          </div>

          <div v-if="errorMessage" role="alert" class="alert alert-error">
            {{ errorMessage }}
          </div>

          <div
            v-if="successMessage"
            role="status"
            aria-live="polite"
            class="alert alert-success"
          >
            {{ successMessage }}
          </div>
        </div>

        <div
          class="flex flex-col-reverse gap-3 border-t border-base-300 bg-base-200/60 px-6 py-5 sm:flex-row sm:items-center sm:justify-between sm:px-8"
        >
          <NuxtLink to="/scan" class="btn btn-ghost min-h-12">
            ไปหน้าสแกนรูป
          </NuxtLink>

          <button
            type="submit"
            class="btn btn-primary min-h-12 min-w-44"
            :disabled="!canSave"
          >
            <span v-if="saving" class="loading loading-spinner loading-sm" />
            {{ saving ? "กำลังบันทึก..." : "บันทึกการตั้งค่า" }}
          </button>
        </div>
      </form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import {
  SelectDestinationDirectory,
  SelectSourceDirectory,
} from "../../wailsjs/go/main/App";

import {
  GetSetting,
  SaveSetting,
} from "../../wailsjs/go/controller/SettingController";

const source = ref("");
const destination = ref("");

const loading = ref(true);
const saving = ref(false);
const selecting = ref<"source" | "destination" | null>(null);

const errorMessage = ref("");
const successMessage = ref("");

const sameDirectory = computed(() => {
  if (!source.value || !destination.value) {
    return false;
  }

  return source.value.toLowerCase() === destination.value.toLowerCase();
});

const canSave = computed(() => {
  return (
    source.value.trim() !== "" &&
    destination.value.trim() !== "" &&
    !sameDirectory.value &&
    !loading.value &&
    !saving.value
  );
});

function getErrorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }

  return String(error);
}

function clearMessage() {
  errorMessage.value = "";
  successMessage.value = "";
}

async function loadSetting() {
  loading.value = true;
  clearMessage();

  try {
    const setting = await GetSetting();

    source.value = setting.Source ?? "";
    destination.value = setting.Dest ?? "";
  } catch (error) {
    errorMessage.value = getErrorMessage(error);
  } finally {
    loading.value = false;
  }
}

async function chooseSource() {
  clearMessage();
  selecting.value = "source";

  try {
    const selected = await SelectSourceDirectory();

    if (selected) {
      source.value = selected;
    }
  } catch (error) {
    errorMessage.value = getErrorMessage(error);
  } finally {
    selecting.value = null;
  }
}

async function chooseDestination() {
  clearMessage();
  selecting.value = "destination";

  try {
    const selected = await SelectDestinationDirectory();

    if (selected) {
      destination.value = selected;
    }
  } catch (error) {
    errorMessage.value = getErrorMessage(error);
  } finally {
    selecting.value = null;
  }
}

async function saveSetting() {
  if (!canSave.value) {
    return;
  }

  clearMessage();
  saving.value = true;

  try {
    await SaveSetting(source.value, destination.value);
    successMessage.value = "บันทึก Source และ Destination สำเร็จ";
  } catch (error) {
    errorMessage.value = getErrorMessage(error);
  } finally {
    saving.value = false;
  }
}

onMounted(loadSetting);
</script>
