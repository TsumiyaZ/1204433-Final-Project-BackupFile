<template>
  <main
    data-theme="night"
    class="min-h-[calc(100vh-4rem)] bg-base-200 px-6 py-10"
  >
    <button class="btn btn-primary">
      <NuxtLink to="/">index</NuxtLink>
    </button>
    <section class="mx-auto max-w-3xl">
      <div class="mb-8">
        <p
          class="mb-2 text-sm font-semibold uppercase tracking-wider text-primary"
        >
          Application setup
        </p>

        <h1 class="text-3xl font-bold text-base-content">
          กำหนดตำแหน่งจัดเก็บรูปภาพ
        </h1>

        <p class="mt-3 max-w-2xl text-base leading-7 text-base-content/70">
          เลือกโฟลเดอร์ต้นทางที่ต้องการสแกนรูป
          และโฟลเดอร์ปลายทางสำหรับเก็บไฟล์ที่ Backup
        </p>
      </div>

      <form
        class="card border border-base-300 bg-base-100 shadow-xl"
        @submit.prevent="saveSetting"
      >
        <div class="card-body gap-7">
          <div>
            <label
              for="source"
              class="mb-2 block font-semibold text-base-content"
            >
              Source folder
            </label>

            <div class="flex flex-col gap-3 sm:flex-row">
              <input
                id="source"
                v-model="source"
                class="input input-bordered min-h-12 w-full"
                placeholder="ยังไม่ได้เลือกโฟลเดอร์ต้นทาง"
                readonly
              />

              <button
                type="button"
                class="btn btn-outline min-h-12 shrink-0"
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

            <p class="mt-2 text-sm text-base-content/60">
              โปรแกรมจะค้นหาเฉพาะไฟล์รูปภาพจากโฟลเดอร์นี้
            </p>
          </div>

          <div>
            <label
              for="destination"
              class="mb-2 block font-semibold text-base-content"
            >
              Destination folder
            </label>

            <div class="flex flex-col gap-3 sm:flex-row">
              <input
                id="destination"
                v-model="destination"
                class="input input-bordered min-h-12 w-full"
                placeholder="ยังไม่ได้เลือกโฟลเดอร์ปลายทาง"
                readonly
              />

              <button
                type="button"
                class="btn btn-outline min-h-12 shrink-0"
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

            <p class="mt-2 text-sm text-base-content/60">
              รูปที่เลือก Backup จะถูกย้ายมาเก็บในโฟลเดอร์นี้
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

          <div class="card-actions justify-end border-t border-base-300 pt-6">
            <button
              type="submit"
              class="btn btn-primary min-h-12 min-w-36"
              :disabled="!canSave"
            >
              <span v-if="saving" class="loading loading-spinner loading-sm" />
              {{ saving ? "กำลังบันทึก..." : "บันทึกการตั้งค่า" }}
            </button>
          </div>
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
