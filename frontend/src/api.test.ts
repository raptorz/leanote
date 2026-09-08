import { afterEach, describe, expect, it, vi } from 'vitest'
import { objectId, request, upload } from './api'

afterEach(()=>vi.unstubAllGlobals())
describe('legacy API adapter',()=>{
 it('encodes indexed arrays and leaves password text unchanged',async()=>{
  const fetch=vi.fn().mockResolvedValue({ok:true,json:async()=>({Ok:true})});vi.stubGlobal('fetch',fetch)
  await request('/share/addShareNote',{emails:['a@b.test','c@d.test'],pwd:'a&b+c'})
  const options=fetch.mock.calls[0][1]
  expect(options.body.get('emails[1]')).toBe('c@d.test')
  expect(options.body.get('pwd')).toBe('a&b+c')
  expect(options.credentials).toBe('same-origin')
 })
  it('does not treat an unsuccessful HTTP-200 write as saved',async()=>{
   vi.stubGlobal('fetch',vi.fn().mockResolvedValue({ok:true,status:200,json:async()=>({Ok:false,Msg:'conflict'})}))
   await expect(request('/web/save',{})).rejects.toThrow('conflict')
  })
  it('surfaces a friendly error when the backend answers with non-JSON',async()=>{
   vi.stubGlobal('fetch',vi.fn().mockResolvedValue({ok:false,status:404,json:async()=>{throw new SyntaxError('Unexpected token <')}}))
   await expect(request('/web/document',{noteId:'x'})).rejects.toThrow('请求失败 (404)')
  })
  it('rejects bare-false upload responses instead of reporting success',async()=>{
   vi.stubGlobal('fetch',vi.fn().mockResolvedValue({ok:true,status:200,json:async()=>false}))
   await expect(upload('/attach/uploadAttach',new File(['x'],'a.txt'),{noteId:'x'})).rejects.toThrow('上传失败')
  })
 it('generates Mongo-compatible identifiers without duplicates',()=>{
  const ids=Array.from({length:100},objectId)
  expect(new Set(ids).size).toBe(100)
  ids.forEach(id=>expect(id).toMatch(/^[0-9a-f]{24}$/))
 })
})
