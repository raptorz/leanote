<script setup lang="ts">
import {computed,ref} from 'vue'
import brandMark from '../assets/gemsnote_s.png'
const props=defineProps<{admin?:boolean,user?:{Username?:string;Email?:string;Logo?:string},showSpacesButton?:boolean,workspace?:boolean}>()
const emit=defineEmits<{showSpaces:[]}>()
const menuOpen=ref(false)
const initials=computed(()=>(props.user?.Username||props.user?.Email||'用户').trim().slice(0,1).toUpperCase())
const avatarSrc=computed(()=>{const logo=props.user?.Logo||'';return /^[a-f\d]{24}$/i.test(logo)?`/api/file/getImage?fileId=${logo}`:logo})
function close(){menuOpen.value=false}
function focusOut(event:FocusEvent){if(!(event.currentTarget as HTMLElement).contains(event.relatedTarget as Node|null))close()}
</script>
<template><nav class="rail" aria-label="全局导航"><RouterLink to="/note" class="mark" title="珠玑笔记"><img :src="brandMark" alt="珠玑笔记"></RouterLink><button v-if="workspace && showSpacesButton" class="rail-link" title="展开我的空间" aria-label="笔记，展开我的空间" @click="emit('showSpaces')">笔记</button><RouterLink v-else to="/note" title="笔记">笔记</RouterLink><RouterLink v-if="admin" to="/admin" title="系统管理">管理</RouterLink><div class="rail-bottom"><div class="user-menu" @focusout="focusOut" @keydown.escape="close"><button class="avatar-button" :aria-expanded="menuOpen" aria-haspopup="menu" aria-label="用户菜单" @click="menuOpen=!menuOpen"><img v-if="avatarSrc" :src="avatarSrc" alt=""><span v-else>{{initials}}</span></button><div v-if="menuOpen" class="user-popover" role="menu"><RouterLink to="/member" role="menuitem" @click="close">账号</RouterLink><a href="/logout" role="menuitem">退出</a></div></div></div></nav></template>
