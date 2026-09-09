<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { request } from '../api'
const route=useRoute(), router=useRouter()
const email=ref(''), pwd=ref(''), captcha=ref(''), error=ref(''), message=ref(''), busy=ref(false), needCaptcha=ref(false), openRegister=ref(false), captchaUrl=ref('/captcha/get'), isDesktop=ref(false), server=ref('')
const mode=computed(()=>route.path.startsWith('/findPassword')?'reset':route.path==='/register'?'register':'login')
onMounted(async()=>{ try{const data=await request('/web/bootstrap');openRegister.value=data.OpenRegister;needCaptcha.value=data.NeedCaptcha;isDesktop.value=!!data.Desktop;server.value=data.Host||'';if(data.User)router.replace('/note')}catch(e){error.value=String(e)} })
async function submit(){busy.value=true;error.value='';message.value='';try{
 if(mode.value==='login'){await request('/doLogin',{email:email.value,pwd:pwd.value,captcha:captcha.value,host:desktopHost()});await router.replace('/note')}
 else if(mode.value==='register'){await request('/doRegister',{email:email.value,pwd:pwd.value,iu:String(route.query.iu||''),host:desktopHost()});await router.replace('/note')}
 else if(route.params.token){await request('/findPasswordUpdate',{token:String(route.params.token),pwd:pwd.value,host:desktopHost()});message.value='密码已修改，请重新登录。'}
 else {await request('/doFindPassword',{email:email.value,host:desktopHost()});message.value='请检查邮箱中的重置密码链接。'}
}catch(e){error.value=String(e);const data=await request('/web/bootstrap').catch(()=>({}));needCaptcha.value=data.NeedCaptcha}finally{busy.value=false}}
function desktopHost(){return isDesktop.value?server.value.trim():undefined}
</script>
<template><main class="auth-shell"><section class="auth-card"><div class="brand">◈ 珠玑笔记</div><p class="muted">Pearlnote · 记录、整理、沉淀</p><h1>{{mode==='login'?'欢迎回来':mode==='register'?'创建账号':'找回密码'}}</h1><form @submit.prevent="submit"><label v-if="isDesktop">服务器<input v-model="server" placeholder="https://your-pearlnote-server" autocomplete="url"></label><label v-if="!route.params.token">邮箱或用户名<input v-model="email" autocomplete="username" required></label><label v-if="mode!=='reset'||route.params.token">密码<input v-model="pwd" type="password" :autocomplete="mode==='login'?'current-password':'new-password'" required></label><label v-if="needCaptcha&&mode==='login'">验证码<img :src="captchaUrl" alt="验证码"><input v-model="captcha" required></label><p v-if="error" role="alert" class="error">{{error}}</p><p v-if="message" role="status">{{message}}</p><button class="primary" :disabled="busy">{{busy?'处理中…':'继续'}}</button></form><nav class="inline"><RouterLink to="/login">登录</RouterLink><RouterLink v-if="openRegister" to="/register">注册</RouterLink><RouterLink to="/findPassword">找回密码</RouterLink></nav></section></main></template>
