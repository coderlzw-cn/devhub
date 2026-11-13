<script setup lang="ts">
import { cn } from '@/utils/cn'
import { computed, ref } from 'vue'
import { useQueryClient } from '@tanstack/vue-query'
import {
  ArrowLeft,
  ArrowUpRight,
  BarChart3,
  BookOpen,
  CloudUpload,
  Download,
  Folder,
  Pencil,
  Plus,
  Settings2,
  UserRound,
} from 'lucide-vue-next'
import type { ProfileData } from '@/data/profile'
import profileBanner from '@/assets/profile-space-banner.png'
import Select from '@/components/ui/Select.vue'
import { useResourcesQuery } from '@/composables/useResourcesQuery'
import { useProfileQuery, profileQueryKey } from '@/composables/useProfileQuery'
import Skeleton from '@/components/ui/Skeleton.vue'
import ProfileEditorDialog, { type ProfileFormValues } from '@/pages/profile/components/ProfileEditorDialog.vue'

const props = defineProps<{
  themeMode: string
  defaultPage: string
  isDark: boolean
}>()
const queryClient = useQueryClient()
const resourcesQuery = useResourcesQuery()
const profileQuery = useProfileQuery()
const resources = computed(() => resourcesQuery.data.value ?? [])
const profile = computed<ProfileData>(
  () => profileQuery.data.value ?? { name: '', email: '', bio: '', avatar: '', joinedAt: '' },
)
const isPending = computed(() => resourcesQuery.isPending.value || profileQuery.isPending.value)
const isError = computed(() => resourcesQuery.isError.value || profileQuery.isError.value)
const isFetching = computed(() => resourcesQuery.isFetching.value || profileQuery.isFetching.value)
const emit = defineEmits<{
  (event: 'update-preference', key: 'themeMode' | 'defaultPage', value: string): void
  (event: 'navigate', page: string): void
  (event: 'notify', message: string): void
}>()
const showEditor = ref(false)
const categories = computed(() => [...new Set(resources.value.map((resource) => resource.category))])
const statCardClasses: Record<string, string> = {
  purple: 'bg-gradient-to-br from-purple-950 to-slate-900 text-purple-400',
  blue: 'bg-gradient-to-br from-blue-950 to-slate-900 text-sky-400',
  green: 'bg-gradient-to-br from-emerald-950 to-slate-900 text-emerald-400',
  orange: 'bg-gradient-to-br from-orange-950 to-slate-900 text-orange-400',
}
const statCardLightClasses: Record<string, string> = {
  purple: 'border-violet-100 !bg-none bg-violet-50 text-violet-600',
  blue: 'border-blue-100 !bg-none bg-blue-50 text-blue-600',
  green: 'border-emerald-100 !bg-none bg-emerald-50 text-emerald-600',
  orange: 'border-orange-100 !bg-none bg-orange-50 text-orange-600',
}
function editProfile() {
  showEditor.value = true
}
function saveProfile(values: ProfileFormValues) {
  const updatedProfile = { ...profile.value, ...values }
  queryClient.setQueryData(profileQueryKey, updatedProfile)
  localStorage.setItem('devhub-local-profile-v2', JSON.stringify(updatedProfile))
  showEditor.value = false
}
function downloadData() {
  const payload = {
    exportedAt: new Date().toISOString(),
    profile: profile.value,
    resources: resources.value,
    preferences: { themeMode: props.themeMode, defaultPage: props.defaultPage },
  }
  const blob = new Blob([JSON.stringify(payload, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'devhub-backup.json'
  link.click()
  URL.revokeObjectURL(url)
  emit('notify', '数据备份已导出')
}
function retryQueries() {
  void Promise.all([resourcesQuery.refetch(), profileQuery.refetch()])
}
</script>

<template>
  <div
    v-if="isPending"
    class="mx-auto grid min-h-[50vh] max-w-[1600px] content-start gap-4 px-3.5 pt-5 sm:px-6 sm:pt-7 lg:px-[42px]"
    aria-label="个人中心加载中"
  >
    <Skeleton class="h-60 w-full rounded-xl" />
    <div class="grid gap-3 xl:grid-cols-2">
      <Skeleton v-for="item in 4" :key="item" class="h-48 rounded-xl" />
    </div>
  </div>
  <div v-else-if="isError" class="grid min-h-[50vh] place-content-center gap-3 text-center">
    <p class="text-sm text-rose-300">个人资料加载失败，请检查连接后重试。</p>
    <button
      class="mx-auto rounded-md border border-slate-700 px-3 py-1.5 text-xs text-slate-200 disabled:opacity-50"
      :disabled="isFetching"
      @click="retryQueries"
    >
      {{ isFetching ? '正在重试…' : '重试' }}
    </button>
  </div>
  <div v-else>
    <div class="mx-auto max-w-[1600px] px-3.5 pt-3 pb-[30px] sm:px-6 sm:pt-[29px] sm:pb-[50px] lg:px-[42px]">
      <button
        class="mb-3 inline-flex items-center gap-2 rounded-lg px-2.5 py-2 text-xs text-slate-300 transition hover:bg-slate-800 hover:text-white sm:mb-4 light:text-slate-600 light:hover:bg-slate-100 light:hover:text-slate-900"
        @click="emit('navigate', '首页')"
      >
        <ArrowLeft :size="16" />返回资源中心
      </button>
      <section
        class="relative flex min-h-[270px] items-start overflow-hidden rounded-[15px] border border-blue-900 bg-slate-950 p-4 sm:h-[238px] sm:min-h-0 sm:items-center sm:p-7"
      >
        <img :src="profileBanner" alt="" class="absolute inset-0 size-full object-cover object-[center_47%]" />
        <div class="absolute inset-0 bg-linear-to-r from-[#050e20]/95 via-[#051025]/80 to-[#051025]/10"></div>
        <div
          class="relative z-10 flex min-w-0 flex-col items-start gap-[13px] sm:flex-row sm:items-center sm:gap-[25px]"
        >
          <div
            class="relative grid size-[75px] shrink-0 place-items-center rounded-full border-4 border-sky-100/30 bg-linear-to-br from-sky-300 via-blue-900 to-violet-500 text-[28px] font-extrabold text-white shadow-[0_8px_35px_#1b9dff55] sm:size-[116px] sm:text-[40px]"
          >
            {{ profile.avatar || profile.name.slice(0, 1)
            }}<button
              class="absolute right-0 bottom-0 grid size-[31px] place-items-center rounded-full border-2 border-blue-950 bg-violet-600 text-white"
              title="编辑头像"
              @click="editProfile"
            >
              <Pencil :size="13" />
            </button>
          </div>
          <div class="min-w-0">
            <h1 class="mb-2 text-[23px] tracking-tight text-white sm:text-[30px]">{{ profile.name || '本地资料' }}</h1>
            <p class="text-[11px] text-blue-100 sm:text-sm">{{ profile.bio }}</p>
            <div class="mt-[17px] flex gap-[15px] sm:mt-[23px] sm:gap-[27px]">
              <div class="grid grid-cols-[22px_auto] items-center gap-x-2">
                <Folder class="row-span-2 size-[21px] text-sky-400" /><strong
                  class="text-base/tight text-white sm:text-xl"
                  >{{ resources.length }}</strong
                ><span class="mt-[3px] text-[9px] text-blue-100/70 sm:text-[10px]">浏览资源</span>
              </div>
            </div>
          </div>
        </div>
        <button
          class="absolute right-3 bottom-3 z-10 flex items-center gap-2 rounded-[10px] border border-violet-300/40 bg-indigo-900/70 px-2.5 py-2 text-[10px] text-white backdrop-blur-sm sm:right-[22px] sm:bottom-5 sm:px-4 sm:py-[11px] sm:text-xs"
          @click="editProfile"
        >
          <Pencil :size="15" />编辑资料
        </button>
      </section>
      <div
        class="mt-[10px] grid gap-[10px] xl:mt-[14px] xl:grid-cols-[minmax(0,1.25fr)_minmax(360px,.9fr)] xl:gap-[14px]"
      >
        <div class="flex min-w-0 flex-col gap-[10px] xl:gap-[14px]">
          <section
            class="relative min-w-0 overflow-hidden rounded-[14px] border border-slate-700 bg-linear-to-br from-slate-800 to-[#121b2a] p-3.5 sm:p-[17px_18px] light:border-slate-200! light:bg-white! light:bg-none!"
          >
            <div class="mb-[15px] flex items-center justify-between gap-2.5">
              <div class="flex items-center gap-2.5">
                <BarChart3 class="size-[22px] text-violet-400" />
                <h2 class="text-sm font-bold text-slate-100 light:text-slate-800">数据概览</h2>
              </div>
              <button
                class="flex items-center gap-1 text-[10px] whitespace-nowrap text-indigo-300 light:text-indigo-700"
                @click="emit('navigate', '资源管理')"
              >
                查看全部 <ArrowUpRight :size="14" />
              </button>
            </div>
            <div class="grid grid-cols-2 gap-2.5 sm:grid-cols-3">
              <div
                v-for="item in [
                  { label: '全部资源', value: resources.length, icon: Folder, tone: 'blue' },
                  {
                    label: '置顶站点',
                    value: resources.filter((resource) => resource.isPinned).length,
                    icon: Plus,
                    tone: 'orange',
                  },
                ]"
                :key="item.label"
                class="relative flex h-[105px] flex-col items-start gap-[5px] overflow-hidden rounded-xl border border-slate-700 p-[13px] sm:h-[122px] light:border-slate-200"
                :class="cn(statCardClasses[item.tone], !isDark && statCardLightClasses[item.tone])"
              >
                <component :is="item.icon" class="mb-0.5 size-[27px]" /><strong
                  class="text-[21px] leading-none text-slate-100 light:text-slate-800"
                  >{{ item.value }}</strong
                ><span class="text-[10px] text-slate-300 light:text-slate-600">{{ item.label }}</span
                ><i
                  class="absolute -right-6 -bottom-6 h-11 w-[110px] rotate-[-25deg] rounded-[50%] border-t border-current opacity-50"
                ></i>
              </div>
            </div>
          </section>
          <section
            class="relative min-w-0 overflow-hidden rounded-[14px] border border-slate-700 bg-linear-to-br from-slate-800 to-[#121b2a] p-3.5 sm:p-[17px_18px] light:border-slate-200! light:bg-white! light:bg-none!"
          >
            <div class="mb-[15px] flex items-center justify-between">
              <div class="flex items-center gap-2.5">
                <Folder class="size-[22px] text-violet-400" />
                <h2 class="text-sm font-bold text-slate-100 light:text-slate-800">资源分类</h2>
              </div>
            </div>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="category in categories"
                :key="category"
                class="rounded-full border border-indigo-900 bg-indigo-950 px-[13px] py-[7px] text-[10px] text-indigo-200 light:border-indigo-100 light:bg-indigo-50 light:text-indigo-700"
                >{{ category }}</span
              ><span v-if="categories.length === 0" class="text-[10px] text-slate-400">暂无分类</span>
            </div>
          </section>
        </div>
        <div class="flex min-w-0 flex-col gap-[10px] xl:gap-[14px]">
          <section
            class="relative min-w-0 overflow-hidden rounded-[14px] border border-slate-700 bg-linear-to-br from-slate-800 to-[#121b2a] p-3.5 sm:min-h-[222px] sm:p-[17px_18px] light:border-slate-200! light:bg-white! light:bg-none!"
          >
            <div class="mb-[15px] flex items-center justify-between gap-2.5">
              <div class="flex items-center gap-2.5">
                <UserRound class="size-[22px] text-violet-400" />
                <h2 class="text-sm font-bold text-slate-100 light:text-slate-800">本地个人信息</h2>
              </div>
              <button
                class="flex items-center gap-2 rounded-lg border border-indigo-900 bg-indigo-950 px-2.5 py-[7px] text-[10px] text-white light:border-indigo-200 light:bg-indigo-50 light:text-indigo-700"
                @click="editProfile"
              >
                <Pencil :size="13" />修改资料
              </button>
            </div>
            <dl class="mt-[18px] flex flex-col gap-3.5">
              <div class="grid grid-cols-[82px_minmax(0,1fr)] items-start gap-2.5 sm:grid-cols-[100px_minmax(0,1fr)]">
                <dt class="text-[11px] text-slate-400 light:text-slate-500">用户名</dt>
                <dd class="m-0 min-w-0 text-[11px] wrap-break-word text-slate-300 light:text-slate-700">
                  {{ profile.name }}
                </dd>
              </div>
              <div class="grid grid-cols-[82px_minmax(0,1fr)] items-start gap-2.5 sm:grid-cols-[100px_minmax(0,1fr)]">
                <dt class="text-[11px] text-slate-400 light:text-slate-500">邮箱</dt>
                <dd class="m-0 min-w-0 text-[11px] wrap-break-word text-slate-300 light:text-slate-700">
                  {{ profile.email || '尚未填写' }}
                </dd>
              </div>
              <div class="grid grid-cols-[82px_minmax(0,1fr)] items-start gap-2.5 sm:grid-cols-[100px_minmax(0,1fr)]">
                <dt class="text-[11px] text-slate-400 light:text-slate-500">个人简介</dt>
                <dd class="m-0 min-w-0 text-[11px] wrap-break-word text-slate-300 light:text-slate-700">
                  {{ profile.bio }}
                </dd>
              </div>
              <div class="grid grid-cols-[82px_minmax(0,1fr)] items-start gap-2.5 sm:grid-cols-[100px_minmax(0,1fr)]">
                <dt class="text-[11px] text-slate-400 light:text-slate-500">注册时间</dt>
                <dd class="m-0 min-w-0 text-[11px] wrap-break-word text-slate-300 light:text-slate-700">
                  {{ profile.joinedAt }}
                </dd>
              </div>
            </dl>
          </section>
          <section
            class="relative min-w-0 overflow-hidden rounded-[14px] border border-slate-700 bg-linear-to-br from-slate-800 to-[#121b2a] p-3.5 sm:min-h-[194px] sm:p-[17px_18px] light:border-slate-200! light:bg-white! light:bg-none!"
          >
            <div class="mb-3 flex items-center gap-2.5">
              <Settings2 class="size-[22px] text-violet-400" />
              <h2 class="text-sm font-bold text-slate-100 light:text-slate-800">偏好设置</h2>
            </div>
            <div
              class="mt-[9px] grid grid-cols-[92px_minmax(0,1fr)] items-center gap-2.5 sm:grid-cols-[118px_minmax(0,1fr)]"
            >
              <label class="flex items-center gap-2 text-[10px] whitespace-nowrap text-slate-300 light:text-slate-600"
                ><UserRound :size="16" class="text-indigo-300" />主题模式</label
              >
              <div
                class="flex gap-[3px] rounded-[10px] border border-slate-700 bg-slate-800 p-[3px] light:border-slate-200 light:bg-slate-100"
              >
                <button
                  v-for="option in [
                    { label: '跟随系统', value: 'system' },
                    { label: '深色', value: 'dark' },
                    { label: '浅色', value: 'light' },
                  ]"
                  :key="option.value"
                  class="h-[29px] flex-1 rounded-lg text-[9px] whitespace-nowrap text-slate-300 sm:text-[10px] light:text-slate-600"
                  :class="cn(themeMode === option.value && 'bg-linear-to-r from-violet-600 to-indigo-600 text-white')"
                  @click="emit('update-preference', 'themeMode', option.value)"
                >
                  {{ option.label }}
                </button>
              </div>
            </div>
            <div
              class="mt-[9px] grid grid-cols-[92px_minmax(0,1fr)] items-center gap-2.5 sm:grid-cols-[118px_minmax(0,1fr)]"
            >
              <label class="flex items-center gap-2 text-[10px] whitespace-nowrap text-slate-300 light:text-slate-600"
                ><BookOpen :size="16" class="text-indigo-300" />默认打开页</label
              ><Select
                :model-value="defaultPage"
                :options="['首页', '资源管理', '分类管理']"
                class="w-full text-[10px]"
                @update:model-value="emit('update-preference', 'defaultPage', String($event))"
              />
            </div>
          </section>
          <section
            class="relative flex min-h-[105px] flex-wrap items-center gap-3 overflow-hidden rounded-[14px] border border-slate-700 bg-linear-to-br from-slate-800 to-[#121b2a] p-3.5 sm:flex-nowrap sm:p-[17px_18px] light:border-slate-200! light:bg-white! light:bg-none!"
          >
            <div
              class="grid size-[45px] shrink-0 place-items-center rounded-xl bg-linear-to-br from-blue-900 to-violet-700 text-blue-200"
            >
              <CloudUpload :size="30" />
            </div>
            <div class="min-w-[calc(100%-65px)] flex-1 sm:min-w-0">
              <h2 class="mb-1.5 text-[13px] font-semibold text-slate-100 light:text-slate-800">数据导出与备份</h2>
              <p class="text-[10px] text-slate-400 light:text-slate-500">导出个人资料、资源和偏好设置</p>
            </div>
            <button
              class="ml-[58px] flex items-center gap-2 rounded-lg border border-violet-400 bg-linear-to-r from-violet-600 to-indigo-600 px-[13px] py-2.5 text-[10px] whitespace-nowrap text-white sm:ml-0"
              @click="downloadData"
            >
              <Download :size="16" />导出数据
            </button>
          </section>
        </div>
      </div>
      <ProfileEditorDialog v-if="showEditor" :profile="profile" @close="showEditor = false" @save="saveProfile" />
    </div>
  </div>
</template>
