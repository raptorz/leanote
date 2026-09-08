<script setup lang="ts">
defineProps<{ admin?: boolean }>()
</script>
<template><nav class="rail" aria-label="全局导航"><RouterLink to="/note" class="mark" title="珠玑笔记">◈</RouterLink><RouterLink to="/note" title="笔记">笔记</RouterLink><RouterLink v-if="admin" to="/admin" title="系统管理">管理</RouterLink><div class="rail-bottom"><RouterLink to="/member" title="账号管理">账号</RouterLink><a href="/logout">退出</a></div></nav></template>
