<script setup lang="ts">
import {ref,computed,onMounted,onBeforeUnmount,watch} from 'vue'
import {useRoute,useRouter,onBeforeRouteLeave} from 'vue-router'
import {marked} from 'marked'
import DOMPurify from 'dompurify'
import {request,objectId,upload} from '../api'
import Navigation from '../components/Navigation.vue'
import RichEditor from '../components/RichEditor.vue'
const route=useRoute(),router=useRouter(),boot=ref<any>({}),notes=ref<any[]>([]),current=ref<any>(null)
const notebook=ref(''),notebookSearch=ref(''),search=ref(''),tagFilter=ref(''),trash=ref(false),sharedOwner=ref(''),page=ref(1)
const title=ref(''),content=ref(''),tags=ref(''),dirty=ref(false),saving=ref(false),error=ref(''),info=ref(false),preview=ref(false),panel=ref(''),destination=ref(''),shareEmail=ref(''),sharePerm=ref('0'),histories=ref<any[]>([]),members=ref<any[]>([]),attachments=ref<any[]>([]),sort=ref('UpdatedTime'),sidebar=ref(false)
const notebooksVisible=ref(true),notesVisible=ref(true),notebooksWidth=ref(230),notesWidth=ref(290),bookMenu=ref(''),compact=ref(false)
const mobileNotesVisible=ref(false)
let layoutReady=false
function updateCompact(){compact.value=window.innerWidth<=1100}
const imageInput=ref<HTMLInputElement>(),attachInput=ref<HTMLInputElement>(),sharedQueueMsg=ref('')
const isSharedNote=computed(()=>!!current.value?.Note.IsShared)
const sharedCacheLabel=computed(()=>{if(!isSharedNote.value)return'';const state=current.value?.Note.CacheState as string|undefined;const at=current.value?.Note.CachedAt as string|undefined;switch(state){case 'ready':return at?`离线缓存 ${new Date(at).toLocaleString()}`:'离线缓存';case 'stale':return'缓存待更新';case 'pending':return'未下载，请先同步';case 'revoked':return'已撤销访问';default:return''}})
function cacheStateLabel(a:any){switch(a.CacheState){case 'ready':return'已缓存';case 'pending':return'未下载';case 'failed':return'下载失败';default:return''}}
async function queueSharedDownload(attachId:string){try{await request('/attach/queueSharedDownload',{attachId});sharedQueueMsg.value='已加入离线下载队列，稍后自动下载';attachmentsFor(current.value?.Note.NoteId||'')}catch(e){error.value=String(e)}}
let saveTimer:ReturnType<typeof setTimeout>|undefined,loadId=0,savePromise:Promise<boolean>|undefined
function flatten(items:any[],depth=0):any[]{return (items||[]).flatMap(n=>[{...n,depth},...flatten(n.Subs,depth+1)])}
const notebooks=computed(()=>flatten(boot.value.Notebooks||[]))
const filteredBooks=computed(()=>notebooks.value.filter(n=>n.Title.toLowerCase().includes(notebookSearch.value.toLowerCase())))
const sorted=computed(()=>[...notes.value].sort((a,b)=>sort.value==='Title'?a.Title.localeCompare(b.Title):String(b[sort.value]).localeCompare(String(a[sort.value]))))
const html=computed(()=>DOMPurify.sanitize(current.value?.Note.IsMarkdown?String(marked.parse(content.value,{async:false})):content.value))
const writable=computed(()=>!!current.value?.Writable)
const own=computed(()=>current.value?.Note.UserId===boot.value.User?.UserId)
function layoutKey(){return `pearlnote:workspace:${boot.value.User?.UserId||'guest'}`}
function persistLayout(){if(!layoutReady)return;localStorage.setItem(layoutKey(),JSON.stringify({notebooksVisible:notebooksVisible.value,notesVisible:notesVisible.value,notebooksWidth:notebooksWidth.value,notesWidth:notesWidth.value}))}
function restoreLayout(){
 if(layoutReady)return
 notebooksWidth.value=Math.min(520,Math.max(180,Number(boot.value.User?.NotebookWidth)||230))
 notesWidth.value=Math.min(520,Math.max(180,Number(boot.value.User?.NoteListWidth)||290))
 notebooksVisible.value=!boot.value.User?.LeftIsMin
 try{const saved=JSON.parse(localStorage.getItem(layoutKey())||'null');if(saved){notebooksVisible.value=saved.notebooksVisible!==false;notesVisible.value=saved.notesVisible!==false;notebooksWidth.value=Math.min(520,Math.max(180,Number(saved.notebooksWidth)||notebooksWidth.value));notesWidth.value=Math.min(520,Math.max(180,Number(saved.notesWidth)||notesWidth.value))}}catch{}
 layoutReady=true
}
function showNotebooks(){notebooksVisible.value=true;if(window.innerWidth<=1100)sidebar.value=true;persistLayout()}
function hideNotebooks(){notebooksVisible.value=false;sidebar.value=false;persistLayout()}
function hideNotes(){notesVisible.value=false;persistLayout()}
function showNotes(){notesVisible.value=true;mobileNotesVisible.value=window.innerWidth<=700;sidebar.value=false;persistLayout()}
async function bootstrap(){boot.value=await request('/web/bootstrap');if(!boot.value.User){await router.replace('/login');return false}restoreLayout();return true}
async function load(){try{notes.value=await request(sharedOwner.value?'/share/listShareNotes':'/web/notes',sharedOwner.value?{userId:sharedOwner.value,notebookId:notebook.value,page:page.value,sortField:sort.value,isAsc:false}:{notebookId:notebook.value,key:search.value,tag:tagFilter.value,trash:trash.value,page:page.value,sort:sort.value});notes.value ||= []}catch(e){error.value=String(e)}}
async function attachmentsFor(noteId:string){try{const r:any=await request('/attach/getAttachs',{noteId});attachments.value=r.List||[]}catch(e){attachments.value=[];error.value=String(e)}}
async function refreshDocument(){if(!current.value)return;current.value=await request('/web/document',{noteId:current.value.Note.NoteId})}
async function open(id:string){if(!await flush())return;const token=++loadId;try{const doc=await request('/web/document',{noteId:id});if(token!==loadId)return;current.value=doc;mobileNotesVisible.value=false;title.value=doc.Note.Title;content.value=doc.Content||'';tags.value=(doc.Note.Tags||[]).join(',');dirty.value=false;info.value=false;panel.value='';await attachmentsFor(id);await router.replace(`/note/${id}`)}catch(e){error.value=String(e)}}
async function select(id='',owner='',isTrash=false){if(!await flush())return;showNotes();notebook.value=id;sharedOwner.value=owner;trash.value=isTrash;page.value=1;search.value='';tagFilter.value='';bookMenu.value='';await load();sidebar.value=false}
async function selectBook(id:string){await select(id)}
let stopResize:(()=>void)|undefined
function resizeBy(panelName:'notebooks'|'notes',delta:number){if(panelName==='notebooks')notebooksWidth.value=Math.min(520,Math.max(180,notebooksWidth.value+delta));else notesWidth.value=Math.min(520,Math.max(180,notesWidth.value+delta));persistLayout()}
function resizePanel(panelName:'notebooks'|'notes',event:PointerEvent){
 const startX=event.clientX,startWidth=panelName==='notebooks'?notebooksWidth.value:notesWidth.value
 const move=(e:PointerEvent)=>{const width=Math.min(520,Math.max(180,startWidth+e.clientX-startX));if(panelName==='notebooks')notebooksWidth.value=width;else notesWidth.value=width}
 const stop=()=>{persistLayout();document.body.classList.remove('panel-resizing');window.removeEventListener('pointermove',move);window.removeEventListener('pointerup',stop);window.removeEventListener('pointercancel',stop);stopResize=undefined}
 stopResize=stop
 document.body.classList.add('panel-resizing')
 window.addEventListener('pointermove',move);window.addEventListener('pointerup',stop)
 window.addEventListener('pointercancel',stop)
}
function changed(){dirty.value=true;clearTimeout(saveTimer);saveTimer=setTimeout(()=>save(),1200)}
async function save():Promise<boolean>{
 if(savePromise)return savePromise
 if(!dirty.value||!current.value||!writable.value)return true
 const id=current.value.Note.NoteId,snapshot={title:title.value,content:content.value,tags:tags.value}
 saving.value=true;error.value=''
 savePromise=(async()=>{try{const doc=await request('/web/save',{noteId:id,usn:current.value.Note.Usn,...snapshot});if(current.value?.Note.NoteId===id){current.value=doc;dirty.value=title.value!==snapshot.title||content.value!==snapshot.content||tags.value!==snapshot.tags}await load();return true}catch(e){error.value=String(e).includes('conflict')?'此文章已在其他端修改。请先导出本地内容，再重新加载文章。':String(e);return false}finally{saving.value=false;savePromise=undefined}})()
 return savePromise
}
async function flush(){clearTimeout(saveTimer);if(!await save())return false;return dirty.value?save():true}
async function create(markdown:boolean){if(!await flush())return;let id=notebook.value||(!sharedOwner.value?notebooks.value[0]?.NotebookId:'');if(!id&&!sharedOwner.value){await addBook();id=notebooks.value[0]?.NotebookId}if(!id){error.value='请先选择一个可编辑的共享笔记本';return}try{const doc=await request('/web/save',{noteId:objectId(),notebookId:id,ownerId:sharedOwner.value,title:'未命名文章',content:'',tags:'',isNew:true,isMarkdown:markdown});await load();await open(doc.Note.NoteId)}catch(e){error.value=String(e)}}
async function addBook(){const name=prompt('笔记本名称');if(!name)return;try{await request('/notebook/addNotebook',{notebookId:objectId(),title:name,parentNotebookId:''});await bootstrap()}catch(e){error.value=String(e)}}
async function bookAction(id:string,remove=false){bookMenu.value='';try{if(remove){if(!confirm('删除此笔记本？请先移动其中的文章。'))return;await request('/notebook/deleteNotebook',{notebookId:id});if(notebook.value===id)notebook.value=''}else{const name=prompt('笔记本名称',notebooks.value.find(n=>n.NotebookId===id)?.Title);if(!name)return;await request('/notebook/updateNotebookTitle',{notebookId:id,title:name})}await bootstrap();await load()}catch(e){error.value=String(e)}}
async function remove(){if(!current.value||!confirm(own.value?(trash.value?'永久删除这篇文章？此操作不可撤销。':'将文章放入回收站？'):'从共享列表移除这篇文章？'))return;try{if(own.value&&trash.value)await request('/note/deleteTrash',{noteId:current.value.Note.NoteId});else if(own.value)await request('/note/deleteNote',{noteIds:[current.value.Note.NoteId],isShared:false});else await request('/share/deleteShareNoteBySharedUser',{noteId:current.value.Note.NoteId,fromUserId:current.value.Note.UserId});dirty.value=false;current.value=null;await load()}catch(e){error.value=String(e)}}
async function restore(){try{await request('/web/restore',{noteId:current.value.Note.NoteId});current.value=null;await load()}catch(e){error.value=String(e)}}
async function move(copy=false){if(!await flush()||!destination.value)return;try{await request(copy?'/note/copyNote':'/note/moveNote',{noteIds:[current.value.Note.NoteId],notebookId:destination.value});panel.value='';await load();await open(current.value.Note.NoteId)}catch(e){error.value=String(e)}}
async function share(){try{const result=await request('/share/addShareNote',{noteId:current.value.Note.NoteId,emails:[shareEmail.value],perm:Number(sharePerm.value)});const failures=Object.values(result).filter((r:any)=>!r.Ok);if(failures.length)throw new Error(JSON.stringify(failures));shareEmail.value='';await showShare()}catch(e){error.value=String(e)}}
async function showShare(){panel.value='share';try{const r=await request('/web/shareMembers',{noteId:current.value.Note.NoteId});members.value=r.Users||[]}catch(e){error.value=String(e)}}
async function revoke(id:string){try{await request('/share/deleteShareNote',{noteId:current.value.Note.NoteId,toUserId:id});await showShare()}catch(e){error.value=String(e)}}
async function history(){panel.value='history';try{histories.value=await request('/noteContentHistory/listHistories',{noteId:current.value.Note.NoteId})||[]}catch(e){error.value=String(e)}}
function chooseImage(){imageInput.value?.click()}
function chooseAttach(){attachInput.value?.click()}
function appendUpload(value:string){content.value+=(content.value&& !content.value.endsWith('\n')?'\n':'')+value;changed()}
async function uploadImage(event:Event){const input=event.target as HTMLInputElement,file=input.files?.[0];input.value='';if(!file||!current.value)return;try{const result:any=await upload('/file/pasteImage',file,{noteId:current.value.Note.NoteId});if(!result.Id)throw new Error(result.Msg||'图片上传失败');const src='/api/file/getImage?fileId='+encodeURIComponent(result.Id);appendUpload(current.value.Note.IsMarkdown?`![${file.name}](${src})`:`<img src="${src}" alt="${file.name}">`)}catch(e){error.value=String(e)}}
async function uploadAttach(event:Event){const input=event.target as HTMLInputElement,file=input.files?.[0];input.value='';if(!file||!current.value||!await flush())return;try{await upload('/attach/uploadAttach',file,{noteId:current.value.Note.NoteId});await refreshDocument();await attachmentsFor(current.value.Note.NoteId);panel.value='attachments'}catch(e){error.value=String(e)}}
async function deleteAttach(attachId:string){if(!confirm('删除此附件？')||!await flush())return;try{await request('/attach/deleteAttach',{attachId});await refreshDocument();await attachmentsFor(current.value.Note.NoteId)}catch(e){error.value=String(e)}}
function download(){const blob=new Blob([content.value],{type:current.value.Note.IsMarkdown?'text/markdown':'text/html'});const url=URL.createObjectURL(blob);const a=document.createElement('a');a.href=url;a.download=(title.value||'note')+(current.value.Note.IsMarkdown?'.md':'.html');a.click();URL.revokeObjectURL(url)}
function leave(e:BeforeUnloadEvent){if(dirty.value||saving.value){e.preventDefault();e.returnValue=''}}
function shortcut(e:KeyboardEvent){if(e.key==='Escape')bookMenu.value='';if((e.ctrlKey||e.metaKey)&&e.key==='s'){e.preventDefault();save()}}
function closeBookMenu(){bookMenu.value=''}
onMounted(async()=>{updateCompact();window.addEventListener('resize',updateCompact);window.addEventListener('beforeunload',leave);window.addEventListener('keydown',shortcut);window.addEventListener('pointerdown',closeBookMenu);const wails=(window as any).runtime;if(wails?.EventsOn){wails.EventsOn('shared-notes-revoked',(ids:string[])=>{if(ids?.includes(current.value?.Note?.NoteId)){error.value='该共享笔记已被撤销访问权限';current.value=null;attachments.value=[]}})}try{if(await bootstrap()){await load();if(route.params.noteId)await open(String(route.params.noteId))}}catch(e){error.value=String(e)}})
onBeforeUnmount(()=>{clearTimeout(saveTimer);stopResize?.();document.body.classList.remove('panel-resizing');window.removeEventListener('resize',updateCompact);window.removeEventListener('beforeunload',leave);window.removeEventListener('keydown',shortcut);window.removeEventListener('pointerdown',closeBookMenu)})
onBeforeRouteLeave(async()=>await flush())
watch(()=>route.params.noteId,id=>{if(id&&id!==current.value?.Note.NoteId)open(String(id));else if(!id)flush().then(ok=>{if(ok)current.value=null})})
</script>
<template>
<div class="shell workspace" :class="{'mobile-list-open':mobileNotesVisible}">
<Navigation workspace :admin="boot.IsAdmin" :user="boot.User" :show-spaces-button="!notebooksVisible || (compact && !sidebar)" @show-spaces="showNotebooks"/>
<aside v-if="notebooksVisible" class="notebooks" :class="{mobileOpen:sidebar}" :style="{width:notebooksWidth+'px'}">
<header>
<h2>我的空间</h2>
<div class="panel-actions"><button class="icon-button" @click="addBook" title="新建笔记本" aria-label="新建笔记本">＋</button><button class="icon-button desktop-only" @click="hideNotebooks" title="隐藏我的空间" aria-label="隐藏我的空间">‹</button></div>
</header>
<input v-model="notebookSearch" placeholder="搜索笔记本" aria-label="搜索笔记本">
<button :class="{selected:!notebook&&!trash&&!sharedOwner}" @click="select()">所有文章</button>
<div v-for="n in filteredBooks" :key="n.NotebookId" class="notebook-row" :class="{selected:notebook===n.NotebookId}" :style="{paddingLeft:4+n.depth*16+'px'}">
<button class="notebook-select" @click="selectBook(n.NotebookId)">▱ <span>{{n.Title}}</span><small>{{n.NumberNotes}}</small></button>
<button class="notebook-more" aria-haspopup="menu" :aria-expanded="bookMenu===n.NotebookId" :aria-label="`${n.Title} 菜单`" @pointerdown.stop @click.stop="bookMenu=bookMenu===n.NotebookId?'':n.NotebookId">⋯</button>
<div v-if="bookMenu===n.NotebookId" class="notebook-menu" role="menu" @pointerdown.stop><button role="menuitem" @click="bookAction(n.NotebookId)">重命名</button><button role="menuitem" @click="bookAction(n.NotebookId,true)">删除</button></div>
</div>
<h3>共享给我</h3>
<p v-if="boot.SharedCache==='unsupported'" class="muted">服务端不支持共享离线缓存</p>
<div v-for="(books,owner) in boot.SharedNotebooks" :key="String(owner)">
<button v-for="n in flatten(books)" :key="n.NotebookId" @click="select(n.IsDefault?'':n.NotebookId,String(owner))">♧ {{n.Title||'共享文章'}}</button>
</div>
<p v-if="!Object.keys(boot.SharedNotebooks||{}).length" class="muted">暂无共享内容</p>
<h3>标签</h3>
<div class="tags">
<button v-for="t in boot.Tags" :key="t.Tag" @click="tagFilter=t.Tag;load()">#{{t.Tag}}</button>
</div>
<button @click="select('','',true)">回收站</button>
<div class="resize-handle" role="separator" tabindex="0" aria-orientation="vertical" aria-label="调节我的空间宽度" :aria-valuenow="notebooksWidth" aria-valuemin="180" aria-valuemax="520" @keydown.arrow-left.prevent="resizeBy('notebooks',-10)" @keydown.arrow-right.prevent="resizeBy('notebooks',10)" @pointerdown.prevent="resizePanel('notebooks',$event)"></div>
</aside>
<section v-if="notesVisible" class="note-list" :style="{width:notesWidth+'px'}">
<header>
<button class="mobile-toggle" @click="sidebar=!sidebar">☰</button>
<h2>{{trash?'回收站':sharedOwner?'共享文章':'文章'}}</h2>
<button class="icon-button desktop-only" @click="hideNotes" title="隐藏文章栏" aria-label="隐藏文章栏">‹</button>
</header>
<form @submit.prevent="page=1;load()">
<label class="sort-control" title="排序"><span aria-hidden="true">⇅</span><select v-model="sort" aria-label="文章排序方式" @change="page=1;load()"><option value="UpdatedTime">最近修改</option><option value="CreatedTime">最近创建</option><option value="Title">标题</option></select></label>
<input v-model="search" placeholder="搜索文章">
<button class="icon-button" title="搜索" aria-label="搜索">⌕</button>
</form>
<div class="create-row">
<button class="icon-button" @click="create(false)" :disabled="trash" title="新建笔记" aria-label="新建笔记"><span aria-hidden="true">＋</span></button>
<button class="icon-button" @click="create(false)" :disabled="trash" title="新建富文本笔记" aria-label="新建富文本笔记"><span aria-hidden="true">T＋</span></button>
<button class="icon-button" @click="create(true)" :disabled="trash" title="新建 Markdown 笔记" aria-label="新建 Markdown 笔记"><span class="markdown-icon" aria-hidden="true">M＋</span></button>
</div>
<div class="list-items">
<button v-for="n in sorted" :key="n.NoteId" class="note-item" :class="{selected:current?.Note.NoteId===n.NoteId}" @click="open(n.NoteId)">
<strong>{{n.Title||'未命名'}}</strong>
<p>{{n.Desc||'暂无摘要'}}</p>
<small>{{new Date(n.UpdatedTime).toLocaleDateString()}} {{n.Perm===0?'· 只读':''}}</small>
</button>
<p v-if="!notes.length" class="empty">这里还没有文章</p>
</div>
<footer class="inline">
<button :disabled="page===1" @click="page--;load()">上一页</button>
<span>{{page}}</span>
<button :disabled="notes.length<100" @click="page++;load()">下一页</button>
</footer>
<div class="resize-handle" role="separator" tabindex="0" aria-orientation="vertical" aria-label="调节文章栏宽度" :aria-valuenow="notesWidth" aria-valuemin="180" aria-valuemax="520" @keydown.arrow-left.prevent="resizeBy('notes',-10)" @keydown.arrow-right.prevent="resizeBy('notes',10)" @pointerdown.prevent="resizePanel('notes',$event)"></div>
</section>
<main class="editor">
<p v-if="error" role="alert" class="error">{{error}} <button @click="error=''">关闭</button>
</p>
<template v-if="current">
<header class="toolbar">
<span role="status">{{saving?'保存中…':dirty?'未保存':writable?'已保存':'只读'}}</span>
<span v-if="isSharedNote" class="muted">{{sharedCacheLabel}}</span>
<button :disabled="!writable||saving" @click="save">保存</button>
<button @click="info=!info">ⓘ Info</button>
<button @click="preview=!preview">{{preview?'编辑':'预览'}}</button>
<button v-if="writable&&!trash" @click="chooseImage">图片</button>
<button v-if="writable&&!trash" @click="chooseAttach">附件{{attachments.length?` (${attachments.length})`:''}}</button>
<button v-else-if="attachments.length" @click="panel='attachments'">附件{{` (${attachments.length})`}}</button>
<button v-if="own&&!trash" @click="panel='move'">移动/复制</button>
<button v-if="own&&!trash" @click="showShare">共享</button>
<button v-if="own" @click="history">历史</button>
<button @click="download">导出</button>
<button v-if="trash&&own" @click="restore">恢复</button>
<button v-if="own||writable" @click="remove">删除</button>
<input ref="imageInput" class="visually-hidden" type="file" accept="image/*" @change="uploadImage">
<input ref="attachInput" class="visually-hidden" type="file" @change="uploadAttach">
</header>
<input class="title-input" v-model="title" :readonly="!writable" @input="changed" aria-label="文章标题">
<input class="tag-input" v-model="tags" :readonly="!writable" @input="changed" placeholder="标签，以逗号分隔">
<aside v-if="info" class="info-panel">
<h3>文章信息</h3>
<p>所属笔记本：{{notebooks.find(n=>n.NotebookId===current.Note.NotebookId)?.Title||'共享笔记本'}}</p>
<p>创建：{{new Date(current.Note.CreatedTime).toLocaleString()}}</p>
<p>修改：{{new Date(current.Note.UpdatedTime).toLocaleString()}}</p>
<p>{{content.length}} 字符 · {{current.Note.IsMarkdown?'Markdown':'富文本'}} · {{writable?'可编辑':'只读'}}</p>
<small>{{current.Note.NoteId}}</small>
</aside>
<div v-if="preview||!writable" class="rendered" v-html="html">
</div>
<RichEditor v-else-if="!current.Note.IsMarkdown" :key="current.Note.NoteId" v-model="content" @update:model-value="changed"/>
<textarea v-else v-model="content" class="content-editor" @input="changed" :aria-label="current.Note.IsMarkdown?'Markdown 正文':'HTML 正文'" spellcheck="false">
</textarea>
<section v-if="panel" class="dialog-backdrop" @click.self="panel=''">
<div class="dialog" role="dialog" aria-modal="true">
<button class="close" @click="panel=''">关闭</button>
<template v-if="panel==='move'">
<h2>移动或复制文章</h2>
<input v-model="notebookSearch" placeholder="搜索笔记本">
<select v-model="destination">
<option value="">选择笔记本</option>
<option v-for="n in filteredBooks" :value="n.NotebookId">{{n.Title}}</option>
</select>
<button @click="move()">移动</button>
<button @click="move(true)">复制</button>
</template>
<template v-if="panel==='share'">
<h2>共享文章</h2>
<form @submit.prevent="share">
<label>对方邮箱<input v-model="shareEmail" type="email" required>
</label>
<select v-model="sharePerm">
<option value="0">只读</option>
<option value="1">可编辑</option>
</select>
<button>添加共享</button>
</form>
<p v-for="m in members" :key="m.ToUserId">{{m.Email}} · {{m.Perm?'可编辑':'只读'}} <button v-if="!m.NotebookHasShared" @click="revoke(m.ToUserId)">取消共享</button>
</p>
</template>
<template v-if="panel==='attachments'">
<h2>附件</h2>
<p v-if="sharedQueueMsg" role="status" class="muted">{{sharedQueueMsg}}</p>
<p v-if="!attachments.length" class="muted">暂无附件</p>
<ul class="attachments">
<li v-for="attachment in attachments" :key="attachment.AttachId">
<a :href="`/attach/download?attachId=${attachment.AttachId}`">{{attachment.Title}}</a>
<small>{{Math.ceil((attachment.Size||0)/1024)}} KB</small>
<small v-if="isSharedNote" class="muted">{{cacheStateLabel(attachment)}}</small>
<button v-if="isSharedNote&&attachment.CacheState!=='ready'" @click="queueSharedDownload(attachment.AttachId)">下载供离线使用</button>
<button v-if="writable" @click="deleteAttach(attachment.AttachId)">删除</button>
</li>
</ul>
</template>
<template v-if="panel==='history'">
<h2>历史版本</h2>
<p v-if="!histories.length">暂无历史</p>
<details v-for="h in histories" :key="h.UpdatedTime">
<summary>{{new Date(h.UpdatedTime).toLocaleString()}}</summary>
<pre>{{h.Content}}</pre>
<button :disabled="!writable" @click="content=h.Content;changed();panel=''">恢复到编辑器</button>
</details>
</template>
</div>
</section>
</template>
<div v-else class="welcome">
<div class="pearl">◈</div>
<h1>让每一个想法有所归处</h1>
<p>选择一篇文章，或开始记录新的灵感。</p>
<button class="primary" @click="create(true)">新建文章</button>
</div>
</main>
</div>
</template>
<style scoped>.visually-hidden{position:absolute;width:1px;height:1px;padding:0;margin:-1px;overflow:hidden;clip:rect(0,0,0,0);white-space:nowrap;border:0}.attachments{padding:0;list-style:none}.attachments li{display:flex;align-items:center;gap:10px;padding:10px 0;border-bottom:1px solid #edf0e9}.attachments li a{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.attachments li small{margin-left:auto}.attachments li button{padding:4px 8px;font-size:12px}</style>
