<script setup lang="ts">
import { ref } from 'vue'
import { Field, useForm } from 'vee-validate'
import { ArrowRight, Eye, EyeOff, LockKeyhole, Mail, ShieldCheck, Smartphone, UserRound } from 'lucide-vue-next'
import { cn } from '@/utils/cn'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'

interface LoginPayload {
  account: string
  password?: string
  verificationCode?: string
  remember: boolean
  method: 'password' | 'code'
}

interface LoginFormValues {
  account: string
  password: string
  verificationCode: string
  remember: boolean
}

defineProps<{ isDark: boolean }>()
const emit = defineEmits<{
  login: [payload: LoginPayload]
}>()
const activeTab = ref<'password' | 'code'>('password')
const showPassword = ref(false)
const message = ref('')
const messageType = ref<'error' | 'info'>('error')
const { handleSubmit, isSubmitting } = useForm<LoginFormValues>({
  initialValues: { account: '', password: '', verificationCode: '', remember: true },
})

function selectTab(tab: 'password' | 'code') {
  activeTab.value = tab
  message.value = ''
}

function validateAccount(value: unknown) {
  return String(value ?? '').trim() ? true : '请输入用户名或邮箱'
}

function validatePassword(value: unknown) {
  return activeTab.value !== 'password' || String(value ?? '').length > 0 ? true : '请输入密码'
}

function validateVerificationCode(value: unknown) {
  return activeTab.value !== 'code' || String(value ?? '').trim().length > 0 ? true : '请输入验证码'
}

const submitLogin = handleSubmit(async (values) => {
  message.value = ''
  emit('login', {
    account: values.account.trim(),
    ...(activeTab.value === 'password'
      ? { password: values.password }
      : { verificationCode: values.verificationCode.trim() }),
    remember: values.remember,
    method: activeTab.value,
  })
  await new Promise((resolve) => window.setTimeout(resolve, 450))
  messageType.value = 'info'
  message.value = '登录服务尚未接入，请稍后再试。'
})

function showUnavailable(feature: string) {
  messageType.value = 'info'
  message.value = `${feature}暂未接入，请使用账号登录`
}
</script>

<template>
  <div>
    <div class="mb-7 grid grid-cols-2 border-b" :class="cn(isDark ? 'border-blue-200/15' : 'border-slate-200')">
      <Button
        variant="ghost"
        class="relative w-full! rounded-none! px-0! text-sm font-semibold"
        :class="activeTab === 'password' ? (isDark ? 'text-white' : 'text-blue-700') : (isDark ? 'text-blue-100/55' : 'text-slate-400')"
        @click="selectTab('password')"
      >
        <span class="inline-flex items-center gap-2"><UserRound :size="16" />账号密码登录</span>
        <span v-if="activeTab === 'password'" class="absolute inset-x-5 -bottom-px h-0.75 rounded-full bg-linear-to-r from-cyan-400 to-violet-500 shadow-lg shadow-violet-500/50"></span>
      </Button>
      <Button
        variant="ghost"
        class="relative w-full! rounded-none! px-0! text-sm font-semibold"
        :class="activeTab === 'code' ? (isDark ? 'text-white' : 'text-blue-700') : (isDark ? 'text-blue-100/55' : 'text-slate-400')"
        @click="selectTab('code')"
      >
        <span class="inline-flex items-center gap-2"><Smartphone :size="16" />验证码登录</span>
        <span v-if="activeTab === 'code'" class="absolute inset-x-5 -bottom-px h-0.75 rounded-full bg-linear-to-r from-cyan-400 to-violet-500 shadow-lg shadow-violet-500/50"></span>
      </Button>
    </div>

    <form class="space-y-4" @submit="submitLogin">
      <Field v-slot="{ componentField, errorMessage }" name="account" :rules="validateAccount">
        <Input
          v-bind="componentField"
          type="text"
          appearance="outlined"
          autocomplete="username"
          placeholder="请输入用户名或邮箱"
          class="[&_input]:text-sm! [&_input]:text-white! [&_input]:placeholder:text-blue-100! light:[&_input]:text-slate-900! light:[&_input]:placeholder:text-slate-500! [&>div]:gap-3! [&>div]:rounded-xl! [&>div]:border-blue-200/20! [&>div]:bg-blue-950/55! [&>div]:px-4! [&>div]:focus-within:border-violet-400! [&>div]:focus-within:ring-violet-500/15! light:[&>div]:border-slate-300! light:[&>div]:bg-white! light:[&>div]:focus-within:border-blue-500! light:[&>div]:focus-within:ring-blue-500/15!"
          :disabled="isSubmitting"
          :error="errorMessage"
        >
          <template #leading><Mail :size="19" /></template>
        </Input>
      </Field>

      <Field v-if="activeTab === 'password'" v-slot="{ componentField, errorMessage }" name="password" :rules="validatePassword">
        <Input
          v-bind="componentField"
          :type="showPassword ? 'text' : 'password'"
          appearance="outlined"
          autocomplete="current-password"
          placeholder="请输入密码"
          class="[&_input]:text-sm! [&_input]:text-white! [&_input]:placeholder:text-blue-100! light:[&_input]:text-slate-900! light:[&_input]:placeholder:text-slate-500! [&>div]:gap-3! [&>div]:rounded-xl! [&>div]:border-blue-200/20! [&>div]:bg-blue-950/55! [&>div]:px-4! [&>div]:focus-within:border-violet-400! [&>div]:focus-within:ring-violet-500/15! light:[&>div]:border-slate-300! light:[&>div]:bg-white! light:[&>div]:focus-within:border-blue-500! light:[&>div]:focus-within:ring-blue-500/15!"
          :disabled="isSubmitting"
          :error="errorMessage"
        >
          <template #leading><LockKeyhole :size="19" /></template>
          <template #trailing>
            <Button
              variant="ghost"
              class="w-8! gap-0! rounded-lg! px-0!"
              :aria-label="showPassword ? '隐藏密码' : '显示密码'"
              :disabled="isSubmitting"
              @click="showPassword = !showPassword"
            >
              <Eye v-if="!showPassword" :size="18" /><EyeOff v-else :size="18" />
            </Button>
          </template>
        </Input>
      </Field>

      <div v-else class="flex items-center gap-2">
        <Field v-slot="{ componentField, errorMessage }" name="verificationCode" :rules="validateVerificationCode">
          <Input
            v-bind="componentField"
            type="text"
            appearance="outlined"
            inputmode="numeric"
            autocomplete="one-time-code"
            placeholder="请输入验证码"
            class="min-w-0 flex-1 [&_input]:text-sm! [&_input]:text-white! [&_input]:placeholder:text-blue-100! light:[&_input]:text-slate-900! light:[&_input]:placeholder:text-slate-500! [&>div]:gap-3! [&>div]:rounded-xl! [&>div]:border-blue-200/20! [&>div]:bg-blue-950/55! [&>div]:px-4! light:[&>div]:border-slate-300! light:[&>div]:bg-white!"
            :disabled="isSubmitting"
            :error="errorMessage"
          >
            <template #leading><ShieldCheck :size="19" /></template>
          </Input>
        </Field>
        <Button
          variant="secondary"
          class="shrink-0 rounded-xl! border-violet-300/30! bg-violet-500/15! px-3! text-violet-200! light:border-violet-200! light:bg-violet-50! light:text-violet-700!"
          :disabled="isSubmitting"
          @click="showUnavailable('验证码发送')"
        >
          获取验证码
        </Button>
      </div>

      <div class="flex items-center justify-between pt-1 text-xs">
        <label class="inline-flex cursor-pointer items-center gap-2" :class="cn(isDark ? 'text-blue-100/75' : 'text-slate-600')">
          <Field v-slot="{ field }" name="remember" type="checkbox" :value="true">
            <input v-bind="field" type="checkbox" class="size-4 accent-violet-600" />
          </Field>
          记住登录状态
        </label>
        <Button
          variant="ghost"
          type="button"
          class="rounded-none! px-0! font-medium text-violet-400 hover:bg-transparent! hover:text-violet-300 light:text-violet-700 light:hover:text-violet-800"
          @click="showUnavailable('找回密码')"
        >
          忘记密码？
        </Button>
      </div>

      <p
        v-if="message"
        class="rounded-lg border px-3 py-2 text-xs"
        :class="messageType === 'error' ? (isDark ? 'border-rose-400/20 bg-rose-500/10 text-rose-200' : 'border-rose-200 bg-rose-50 text-rose-600') : (isDark ? 'border-sky-400/20 bg-sky-500/10 text-sky-200' : 'border-sky-200 bg-sky-50 text-sky-700')"
        role="status"
      >
        {{ message }}
      </p>

      <Button
        variant="primary"
        type="submit"
        block
        :loading="isSubmitting"
        class="group mt-2! rounded-xl! border-0! bg-linear-to-r! from-sky-500! via-blue-600! to-violet-600! text-base! font-bold! shadow-lg shadow-blue-700/25"
      >
        登录<ArrowRight :size="18" class="transition group-hover:translate-x-1" />
      </Button>
    </form>
  </div>
</template>
