<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { request } from '../api'

const route = useRoute()
const message = ref('正在验证…')

onMounted(async () => {
  try {
    const result = await request('/web/verifyEmail', {
      token: String(route.query.token || ''),
      change: route.path === '/user/updateEmail',
    })
    message.value = `邮箱 ${result.Email} 验证成功，请重新登录。`
  } catch (error) {
    message.value = String(error)
  }
})
</script>

<template><main class="auth-shell"><section class="auth-card"><h1>邮箱验证</h1><p role="status">{{ message }}</p><RouterLink to="/login">返回登录</RouterLink></section></main></template>
