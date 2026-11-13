<script setup lang="ts">
import { Field, useForm } from 'vee-validate'
import { Check, X } from 'lucide-vue-next'
import type { ProfileData } from '@/data/profile'

export type ProfileFormValues = Pick<ProfileData, 'name' | 'email' | 'bio' | 'avatar'>

const props = defineProps<{ profile: ProfileData }>()
const emit = defineEmits<{ close: []; save: [values: ProfileFormValues] }>()
const { handleSubmit } = useForm<ProfileFormValues>({
  initialValues: { name: props.profile.name, email: props.profile.email, bio: props.profile.bio, avatar: props.profile.avatar },
})

const submit = handleSubmit((values) => emit('save', values))
</script>

<template>
  <div class="fixed inset-0 z-20 grid place-items-center bg-slate-950/75 p-5 backdrop-blur-sm" @click.self="emit('close')">
    <form
      class="flex max-h-[90vh] w-full max-w-[440px] flex-col gap-[13px] overflow-auto rounded-2xl border border-slate-700 bg-[#151f31] p-[22px] shadow-2xl"
      aria-labelledby="profile-editor-title"
      @submit="submit"
    >
      <div class="mb-1 flex items-start justify-between">
        <div>
          <h2 id="profile-editor-title" class="font-semibold">编辑个人资料</h2>
          <p class="mt-1 text-[10px] text-slate-400">更新你的 DevHub 个人信息</p>
        </div>
        <button
          type="button"
          class="grid size-8 place-items-center rounded-lg hover:bg-slate-700"
          aria-label="关闭"
          @click="emit('close')"
        >
          <X :size="18" />
        </button>
      </div>
      <label class="flex flex-col gap-1.5 text-[11px] text-slate-300">
        用户名
        <Field v-slot="{ field, errorMessage }" name="name" :rules="(value: unknown) => String(value ?? '').trim() ? true : '请输入用户名'">
          <input
            v-bind="field"
            class="rounded-md border border-slate-600 bg-slate-900 px-3 py-2.5 text-slate-100 outline-none focus:border-violet-500"
            maxlength="30"
            required
          />
          <span v-if="errorMessage" class="text-rose-400" role="alert">{{ errorMessage }}</span>
        </Field>
      </label>
      <label class="flex flex-col gap-1.5 text-[11px] text-slate-300">
        邮箱
        <Field
          v-slot="{ field, errorMessage }"
          name="email"
          :rules="(value: unknown) => !value || /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(value)) || '请输入有效邮箱地址'"
        >
          <input
            v-bind="field"
            class="rounded-md border border-slate-600 bg-slate-900 px-3 py-2.5 text-slate-100 outline-none focus:border-violet-500"
            type="email"
            placeholder="name@example.com"
          />
          <span v-if="errorMessage" class="text-rose-400" role="alert">{{ errorMessage }}</span>
        </Field>
      </label>
      <label class="flex flex-col gap-1.5 text-[11px] text-slate-300">
        个人简介
        <Field v-slot="{ field }" name="bio">
          <textarea
            v-bind="field"
            class="resize-y rounded-md border border-slate-600 bg-slate-900 px-3 py-2.5 text-slate-100 outline-none focus:border-violet-500"
            rows="3"
            maxlength="100"
          />
        </Field>
      </label>
      <label class="flex flex-col gap-1.5 text-[11px] text-slate-300">
        头像文字
        <Field v-slot="{ field }" name="avatar">
          <input
            v-bind="field"
            class="rounded-md border border-slate-600 bg-slate-900 px-3 py-2.5 text-slate-100 outline-none focus:border-violet-500"
            maxlength="2"
            placeholder="L"
          />
        </Field>
      </label>
      <div class="flex justify-end gap-2 border-t border-slate-700 pt-3">
        <button
          type="button"
          class="rounded-md border border-slate-600 bg-slate-800 px-4 py-2 text-xs text-slate-300"
          @click="emit('close')"
        >
          取消
        </button>
        <button type="submit" class="flex items-center gap-2 rounded-md bg-violet-600 px-4 py-2 text-xs text-white">
          <Check :size="15" />保存修改
        </button>
      </div>
    </form>
  </div>
</template>
