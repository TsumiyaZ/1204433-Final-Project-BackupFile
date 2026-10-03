<template>
  <div
    role="group"
    aria-label="เลือกโฟลเดอร์ที่ต้องการสแกน"
    class="grid gap-4 md:grid-cols-2"
  >
    <button
      v-for="option in options"
      :key="option.mode"
      type="button"
      class="min-h-24 cursor-pointer rounded-[1.5rem] border p-5 text-left transition duration-200 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:cursor-not-allowed disabled:opacity-60"
      :class="
        modelValue === option.mode
          ? 'border-primary bg-primary text-primary-content'
          : 'border-base-300 bg-base-100 text-base-content hover:border-primary/45 hover:bg-primary/5'
      "
      :aria-pressed="modelValue === option.mode"
      :disabled="disabled"
      @click="emit('update:modelValue', option.mode)"
    >
      <span class="flex items-start gap-4">
        <span
          class="grid size-11 shrink-0 place-items-center rounded-xl"
          :class="
            modelValue === option.mode
              ? 'bg-primary-content/15'
              : 'bg-primary/10 text-primary'
          "
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
            <path
              v-if="option.mode === 'source'"
              d="m9 12 3-3 3 3M12 9v7"
            />
            <path v-else d="m9 13 3 3 3-3M12 16V9" />
          </svg>
        </span>

        <span>
          <span class="block text-lg font-bold">{{ option.label }}</span>
          <span
            class="mt-1 block text-sm leading-6"
            :class="
              modelValue === option.mode
                ? 'text-primary-content/75'
                : 'text-base-content/60'
            "
          >
            {{ option.description }}
          </span>
        </span>
      </span>
    </button>
  </div>
</template>

<script setup lang="ts">
import type { ScanMode, ScanOption } from "../../types/photo-scan";

defineProps<{
  modelValue: ScanMode;
  options: ScanOption[];
  disabled: boolean;
}>();

const emit = defineEmits<{
  "update:modelValue": [mode: ScanMode];
}>();
</script>
