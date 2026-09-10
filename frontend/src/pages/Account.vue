<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { request, upload } from '../api'
import Navigation from '../components/Navigation.vue'

type GroupUser = { UserId: string; Username: string; Email: string }
type Group = { GroupId: string; UserId: string; Title: string; Users?: GroupUser[] }

const user = ref<any>({})
const admin = ref(false)
const message = ref('')
const username = ref('')
const email = ref('')
const emailPwd = ref('')
const oldPwd = ref('')
const pwd = ref('')
const groups = ref<Group[]>([])
const newGroupTitle = ref('')
const addUserEmails = ref<Record<string, string>>({})
const editingGroupId = ref('')
const editingGroupTitle = ref('')

function showError(error: unknown) { message.value = error instanceof Error ? error.message : String(error) }
function owned(group: Group) { return group.UserId === user.value.UserId }
async function loadGroups() { try { groups.value = await request<Group[]>('/web/groups') } catch (error) { showError(error) } }

onMounted(async () => {
  try {
    const bootstrap = await request('/web/bootstrap')
    if (!bootstrap.User) { location.href = '/login'; return }
    user.value = bootstrap.User; username.value = bootstrap.User.Username; email.value = bootstrap.User.Email; admin.value = bootstrap.IsAdmin
    await loadGroups()
  } catch (error) { showError(error) }
})

async function update(path: string, data: any) {
  try { await request(path, data); message.value = '已更新'; if (path.endsWith('updatePwd')) location.href = '/login' } catch (error) { showError(error) }
}
async function avatar(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  try { const result = await upload('/file/uploadAvatar', file); user.value.Logo = result.Id; message.value = '头像已更新' } catch (error) { showError(error) }
}
async function addGroup() {
  const title = newGroupTitle.value.trim()
  if (!title) { message.value = '请输入分组名称'; return }
  try { await request('/member/group/addGroup', { title }); newGroupTitle.value = ''; message.value = '分组已创建'; await loadGroups() } catch (error) { showError(error) }
}
function startRename(group: Group) { editingGroupId.value = group.GroupId; editingGroupTitle.value = group.Title }
async function renameGroup(group: Group) {
  const title = editingGroupTitle.value.trim()
  if (!title) { message.value = '请输入分组名称'; return }
  try { await request('/member/group/updateGroupTitle', { groupId: group.GroupId, title }); editingGroupId.value = ''; message.value = '分组名称已更新'; await loadGroups() } catch (error) { showError(error) }
}
async function deleteGroup(group: Group) {
  if (!confirm(`确定删除分组“${group.Title}”吗？`)) return
  try { await request('/member/group/deleteGroup', { groupId: group.GroupId }); message.value = '分组已删除'; await loadGroups() } catch (error) { showError(error) }
}
async function addUser(group: Group) {
  const memberEmail = (addUserEmails.value[group.GroupId] || '').trim()
  if (!memberEmail) { message.value = '请输入用户邮箱'; return }
  try { await request('/member/group/addUser', { groupId: group.GroupId, email: memberEmail }); addUserEmails.value[group.GroupId] = ''; message.value = '用户已加入分组'; await loadGroups() } catch (error) { showError(error) }
}
async function deleteUser(group: Group, member: GroupUser) {
  if (!confirm(`确定将 ${member.Email || member.Username} 移出此分组吗？`)) return
  try { await request('/member/group/deleteUser', { groupId: group.GroupId, userId: member.UserId }); message.value = '用户已移出分组'; await loadGroups() } catch (error) { showError(error) }
}
</script>

<template>
  <div class="shell">
    <Navigation :admin="admin" :user="user" />
    <main class="settings">
      <h1>账号管理</h1><p class="muted">{{ user.Email }}</p><p v-if="message" role="status" class="message">{{ message }}</p>
      <section class="card"><h2>个人资料</h2><img v-if="user.Logo" :src="user.Logo" class="avatar" alt="当前头像"><label>更换头像<input type="file" accept="image/*" @change="avatar"></label><form @submit.prevent="update('/user/updateUsername', { username })"><label>用户名<input v-model="username" required></label><button>更新用户名</button></form><form @submit.prevent="update('/web/emailChange', { email, pwd: emailPwd })"><label>新邮箱<input v-model="email" type="email" required></label><label>当前密码<input v-model="emailPwd" type="password" required></label><button>发送邮箱验证邮件</button></form><button @click="update('/user/reSendActiveEmail', {})">重新发送当前邮箱验证邮件</button></section>
      <section class="card"><h2>修改密码</h2><form @submit.prevent="update('/user/updatePwd', { oldPwd, pwd })"><label>原密码<input v-model="oldPwd" type="password" autocomplete="current-password" required></label><label>新密码<input v-model="pwd" type="password" autocomplete="new-password" required></label><button class="primary">修改并重新登录</button></form></section>
      <section class="card groups">
        <h2>用户分组</h2><p class="muted">将已有账号加入分组后，可在分享笔记或笔记本时选择整个分组。</p>
        <form class="inline-form" @submit.prevent="addGroup"><input v-model="newGroupTitle" maxlength="100" placeholder="新分组名称" aria-label="新分组名称"><button class="primary">新建分组</button></form>
        <p v-if="!groups.length" class="muted">尚无分组。</p>
        <article v-for="group in groups" :key="group.GroupId" class="group">
          <header>
            <form v-if="editingGroupId === group.GroupId" class="inline-form" @submit.prevent="renameGroup(group)"><input v-model="editingGroupTitle" maxlength="100" aria-label="分组名称"><button>保存</button><button type="button" @click="editingGroupId = ''">取消</button></form>
            <template v-else><h3>{{ group.Title }}</h3><span class="muted">{{ group.Users?.length || 0 }} 位成员</span><span v-if="!owned(group)" class="muted">（他人分组）</span><div v-if="owned(group)" class="actions"><button @click="startRename(group)">改名</button><button class="danger" @click="deleteGroup(group)">删除</button></div></template>
          </header>
          <ul v-if="group.Users?.length" class="members"><li v-for="member in group.Users" :key="member.UserId"><span>{{ member.Username || member.Email }}</span><span class="muted">{{ member.Email }}</span><button v-if="owned(group)" class="danger" @click="deleteUser(group, member)">移除</button></li></ul>
          <p v-else class="muted">暂无成员。</p>
          <form v-if="owned(group)" class="inline-form" @submit.prevent="addUser(group)"><input v-model="addUserEmails[group.GroupId]" type="email" placeholder="已有用户的邮箱" :aria-label="`${group.Title} 的成员邮箱`"><button>添加用户</button></form>
        </article>
      </section>
    </main>
  </div>
</template>

<style scoped>
.avatar { width: 72px; height: 72px; border-radius: 50%; object-fit: cover }
.message { color: var(--accent, #1769aa) }
.inline-form, .actions { display: flex; align-items: center; gap: .5rem; flex-wrap: wrap }
.groups { max-width: 760px }.group { padding: 1rem 0; border-top: 1px solid #e5e7eb }.group header { display: flex; align-items: center; gap: .5rem; flex-wrap: wrap }.group h3 { margin: 0; margin-right: .25rem }.actions { margin-left: auto }.members { list-style: none; padding: 0; margin: .75rem 0 }.members li { display: flex; align-items: center; gap: .75rem; padding: .35rem 0 }.members .muted { flex: 1 }.danger { color: #b42318 }
</style>
