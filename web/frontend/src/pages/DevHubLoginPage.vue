<script setup lang="ts">
import { cn } from '@/utils/cn'
import { computed, ref } from 'vue'
import {
  Code2,
  Github,
  Moon,
  Rocket,
  Star,
  Sun,
  Users,
  Zap,
} from 'lucide-vue-next'
import { storeToRefs } from 'pinia'
import Button from '@/components/ui/Button.vue'
import { usePreferencesStore } from '@/stores/preferences'
import LoginForm from '@/pages/login/components/LoginForm.vue'

interface LoginPayload {
  account: string
  password?: string
  verificationCode?: string
  remember: boolean
  method: 'password' | 'code'
}

const emit = defineEmits<{
  login: [payload: LoginPayload]
}>()

const preferences = usePreferencesStore()
const { isDark } = storeToRefs(preferences)
const unavailableMessage = ref('')
const themeButtonLabel = computed(() => (isDark.value ? '切换到浅色主题' : '切换到深色主题'))

function showUnavailable(feature: string) {
  unavailableMessage.value = `${feature}暂未接入，请使用账号登录`
}

function toggleTheme() {
  preferences.setThemeMode(isDark.value ? 'light' : 'dark')
}
</script>

<template>
  <main
    class="relative isolate min-h-screen overflow-hidden"
    :class="[
      isDark
        ? 'bg-[radial-gradient(ellipse_at_top_left,#123f9b_0%,transparent_42%),radial-gradient(ellipse_at_bottom_right,#35128b_0%,transparent_38%),#07101f] text-white'
        : 'bg-[radial-gradient(ellipse_at_top_left,#e4efff_0%,transparent_45%),radial-gradient(ellipse_at_bottom_right,#e8eaff_0%,transparent_42%),#f8fafc] text-slate-900',
      !isDark && 'light',
    ]"
  >
    <div
      class="pointer-events-none absolute inset-0 -z-10 bg-[radial-gradient(circle,rgba(148,163,184,0.18)_1px,transparent_1px)] mask-[linear-gradient(135deg,#000_0%,transparent_72%)] bg-size-[28px_28px]"
      :class="!isDark && 'bg-[radial-gradient(circle,rgba(71,85,105,0.13)_1px,transparent_1px)]'"
    ></div>

    <Button
      variant="ghost"
      class="absolute top-7 right-7 z-20 gap-1! rounded-full! border! p-1! shadow-lg backdrop-blur-xl sm:top-8 sm:right-10"
      :class="cn(isDark ? 'border-white/15 bg-slate-900/55' : 'border-slate-200 bg-white/80')"
      :aria-label="themeButtonLabel"
      @click="toggleTheme"
    >
      <span
        class="grid size-8 place-items-center rounded-full"
        :class="cn(!isDark ? 'bg-blue-600 text-white' : 'text-slate-300')"
        ><Sun :size="16"
      /></span>
      <span
        class="grid size-8 place-items-center rounded-full"
        :class="cn(isDark ? 'bg-violet-600 text-white shadow-md shadow-violet-500/40' : 'text-slate-400')"
        ><Moon :size="16"
      /></span>
    </Button>

    <div
      class="mx-auto grid min-h-screen w-full max-w-[1320px] grid-cols-1 items-center gap-10 px-5 py-24 sm:px-8 lg:grid-cols-[minmax(0,1fr)_minmax(400px,460px)] lg:gap-20 lg:px-14 lg:py-16 xl:px-20"
    >
      <section class="flex flex-col justify-center lg:min-h-[560px]">
        <a href="/" class="relative z-10 inline-flex w-fit items-center gap-3 no-underline">
          <span
            class="grid size-12 place-items-center rounded-xl bg-linear-to-br from-cyan-500 to-blue-700 text-white shadow-lg shadow-blue-900/25"
          >
            <Code2 :size="26" stroke-width="2.5" />
          </span>
          <span class="flex flex-col">
            <strong class="text-[25px] font-bold tracking-tight" :class="cn(isDark ? 'text-white' : 'text-slate-900')"
              >DevHub</strong
            >
            <span class="text-xs tracking-wide" :class="cn(isDark ? 'text-blue-100/65' : 'text-slate-500')"
              >Developer Resource Hub</span
            >
          </span>
        </a>

        <div class="mt-12 max-w-[620px] sm:mt-16 lg:mt-20">
          <p
            class="mb-4 text-xs font-semibold tracking-[0.2em]"
            :class="cn(isDark ? 'text-cyan-300' : 'text-blue-700')"
          >
            DEVELOPER RESOURCE CENTER
          </p>
          <h1
            class="text-4xl leading-[1.2] font-extrabold sm:text-5xl xl:text-[56px]"
            :class="cn(isDark ? 'text-white' : 'text-slate-900')"
          >
            汇聚全球优质<span
              class="bg-linear-to-r from-cyan-400 via-blue-500 to-indigo-500 bg-clip-text text-transparent"
              >开发资源</span
            >
          </h1>
          <p class="mt-5 max-w-lg text-base/7 sm:text-lg" :class="cn(isDark ? 'text-slate-300' : 'text-slate-600')">
            发现值得信赖的开发工具与资源，把常用链接整理成自己的工作空间。
          </p>

          <div class="mt-8 grid max-w-lg grid-cols-2 gap-x-6 gap-y-4">
            <span
              v-for="item in [
                { label: '海量资源', icon: Rocket },
                { label: '开源开放', icon: Star },
                { label: '开发必备', icon: Zap },
                { label: '共同成长', icon: Users },
              ]"
              :key="item.label"
              class="inline-flex items-center gap-2 text-sm font-medium"
              :class="cn(isDark ? 'text-slate-200' : 'text-slate-700')"
            >
              <component :is="item.icon" :size="15" class="text-cyan-400" />{{ item.label }}
            </span>
          </div>
        </div>
      </section>

      <section class="relative flex items-center justify-center lg:py-0">
        <div
          class="relative w-full max-w-115 overflow-hidden rounded-2xl border p-6 shadow-2xl backdrop-blur-xl sm:p-8"
          :class="
            isDark ? 'border-white/10 bg-slate-950/70 shadow-black/30' : 'border-white bg-white/85 shadow-slate-300/50'
          "
        >
          <div class="relative">
            <div class="mb-8">
              <p class="text-2xl font-bold" :class="cn(isDark ? 'text-white' : 'text-slate-900')">欢迎登录 DevHub</p>
              <p class="mt-2 text-sm" :class="cn(isDark ? 'text-blue-100/65' : 'text-slate-500')">
                开启你的开发资源探索之旅
              </p>
            </div>

            <LoginForm :is-dark="isDark" @login="emit('login', $event)" />

            <div
              class="my-7 flex items-center gap-4 text-xs"
              :class="cn(isDark ? 'text-blue-100/45' : 'text-slate-400')"
            >
              <span class="h-px flex-1" :class="cn(isDark ? 'bg-blue-200/15' : 'bg-slate-200')"></span
              >或使用以下方式登录<span
                class="h-px flex-1"
                :class="cn(isDark ? 'bg-blue-200/15' : 'bg-slate-200')"
              ></span>
            </div>

            <div class="flex justify-center gap-5">
              <Button
                variant="ghost"
                type="button"
                class="size-12! gap-0! rounded-full! border! p-0! transition-[transform,background] duration-180 hover:-translate-y-0.5"
                :class="
                  isDark
                    ? 'border-white/10 bg-white/5 text-white hover:bg-white/10'
                    : 'border-slate-200 bg-white text-slate-700 hover:bg-slate-100'
                "
                aria-label="GitHub 登录"
                @click="showUnavailable('GitHub 登录')"
              >
                <Github :size="21" />
              </Button>
              <Button
                variant="ghost"
                type="button"
                class="size-12! gap-0! rounded-full! border! p-0! text-[#4285f4] transition-[transform,background] duration-180 hover:-translate-y-0.5"
                :class="
                  isDark
                    ? 'border-white/10 bg-white/5 hover:bg-white/10'
                    : 'border-slate-200 bg-white hover:bg-slate-100'
                "
                aria-label="Google 登录"
                @click="showUnavailable('Google 登录')"
              >
                <span class="text-xl font-extrabold">G</span>
              </Button>
              <Button
                variant="ghost"
                type="button"
                class="size-12! gap-0! rounded-full! border! p-0! text-emerald-500 transition-[transform,background] duration-180 hover:-translate-y-0.5"
                :class="
                  isDark
                    ? 'border-white/10 bg-white/5 hover:bg-white/10'
                    : 'border-slate-200 bg-white hover:bg-slate-100'
                "
                aria-label="微信登录"
                @click="showUnavailable('微信登录')"
              >
                <Users :size="21" />
              </Button>
            </div>

            <p class="mt-7 text-center text-xs" :class="cn(isDark ? 'text-blue-100/60' : 'text-slate-500')">
              还没有账号？<Button
                variant="ghost"
                type="button"
                class="rounded-none! px-0! font-semibold text-violet-400 hover:bg-transparent! hover:text-violet-300 light:text-violet-700 light:hover:text-violet-800"
                @click="showUnavailable('注册')"
              >
                立即注册
              </Button>
            </p>
            <p v-if="unavailableMessage" class="mt-3 text-center text-xs text-sky-500" role="status">
              {{ unavailableMessage }}
            </p>
          </div>
        </div>
      </section>
    </div>

    <p
      class="pointer-events-none absolute inset-x-0 bottom-4 hidden text-center text-[11px] tracking-widest lg:block"
      :class="cn(isDark ? 'text-blue-100/35' : 'text-slate-400')"
    >
      DEVHUB · DEVELOPER RESOURCE HUB
    </p>
  </main>
</template>
