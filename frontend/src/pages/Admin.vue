<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { request } from '../api'
import Navigation from '../components/Navigation.vue'

const data = ref<any>({ Users: [], Settings: {} })
const message = ref(''), keywords = ref(''), page = ref(1)
const email = ref(''), pwd = ref(''), resetId = ref(''), resetPwd = ref('')

async function load() {
  try { data.value = await request('/web/adminData', { keywords: keywords.value, page: page.value }) }
  catch (error) { message.value = String(error) }
}
async function action(path: string, params: any) {
  try { await request(path, params); message.value = '操作成功'; await load() }
  catch (error) { message.value = String(error) }
}
onMounted(load)
</script>

<template>
  <div class="shell"><Navigation admin/><main class="settings">
    <h1>系统管理</h1><p role="status">{{ message }}</p>
    <section class="card"><h2>用户管理</h2>
      <form class="inline" @submit.prevent="page=1;load()"><input v-model="keywords" placeholder="搜索用户"><button>搜索</button></form>
      <table><thead><tr><th>用户名</th><th>邮箱</th><th>操作</th></tr></thead><tbody>
        <tr v-for="u in data.Users" :key="u.UserId"><td>{{u.Username}}</td><td>{{u.Email}}</td><td><button @click="resetId=u.UserId">重置密码</button></td></tr>
      </tbody></table>
      <div class="inline"><button :disabled="page===1" @click="page--;load()">上一页</button><span>{{page}}</span><button :disabled="data.Users.length<20" @click="page++;load()">下一页</button></div>
      <form v-if="resetId" class="inline" @submit.prevent="action('/web/adminResetPwd',{userId:resetId,pwd:resetPwd})"><label>新密码<input v-model="resetPwd" type="password" required></label><button>确认重置</button><button type="button" @click="resetId=''">取消</button></form>
    </section>
    <section class="card"><h2>添加用户</h2><form @submit.prevent="action('/web/adminRegister',{email,pwd})"><label>邮箱<input v-model="email" type="email" required></label><label>初始密码<input v-model="pwd" type="password" required></label><button>创建</button></form></section>
    <section class="card"><h2>站点与邮件</h2><form @submit.prevent="action('/web/adminSettings',data.Settings)">
      <label>站点 URL<input v-model="data.Settings.siteUrl" type="url"></label>
      <label>开放注册<select v-model="data.Settings.openRegister"><option value="1">是</option><option value="0">否</option></select></label>
      <label>邮件服务器<input v-model="data.Settings.emailHost"></label><label>邮件端口<input v-model="data.Settings.emailPort"></label>
      <label>邮件账号<input v-model="data.Settings.emailUsername"></label><label>邮件密码<input v-model="data.Settings.emailPassword" type="password" placeholder="留空保持原值"></label>
      <label>邮件 SSL<select v-model="data.Settings.emailSSL"><option value="1">启用</option><option value="0">禁用</option></select></label>
      <h3>上传限制（MB）</h3>
      <label>图片<input v-model="data.Settings.uploadImageSize" type="number" min="0" step="0.1"></label><label>头像<input v-model="data.Settings.uploadAvatarSize" type="number" min="0" step="0.1"></label><label>附件<input v-model="data.Settings.uploadAttachSize" type="number" min="0" step="0.1"></label>
      <h3>其他非博客设置</h3>
      <label>PDF 导出程序路径<input v-model="data.Settings.exportPdfBinPath"></label><label>Demo 用户名<input v-model="data.Settings.demoUsername"></label><label>Demo 密码<input v-model="data.Settings.demoPassword" type="password" placeholder="留空保持原值"></label>
      <button class="primary">保存设置</button>
    </form></section>
  </main></div>
</template>
