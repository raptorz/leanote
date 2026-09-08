export type Params = Record<string, string | number | boolean | string[] | undefined>
export async function request<T = any>(path: string, params?: Params): Promise<T> {
  const body = new URLSearchParams()
  for (const [key, value] of Object.entries(params || {})) {
    if (value === undefined) continue
    if (Array.isArray(value)) value.forEach((v, i) => body.append(`${key}[${i}]`, v))
    else body.append(key, String(value))
  }
  const response = await fetch(path, { method: params ? 'POST' : 'GET', credentials: 'same-origin', headers: { 'X-Requested-With': 'XMLHttpRequest' }, body: params ? body : undefined })
  const data = await response.json().catch(() => null)
  if (!response.ok || data === false || data?.Ok === false) {
    if (data?.Msg === 'NOTLOGIN') window.dispatchEvent(new Event('session-expired'))
    throw new Error(data?.Msg || `请求失败 (${response.status})`)
  }
  return data
}
export function objectId(): string {
  return Math.floor(Date.now() / 1000).toString(16).padStart(8, '0') + Array.from(crypto.getRandomValues(new Uint8Array(8)), b => b.toString(16).padStart(2, '0')).join('')
}

export async function upload(path: string, file: File, params: Record<string, string> = {}) {
  const body = new FormData()
  body.append('file', file)
  Object.entries(params).forEach(([key, value]) => body.append(key, value))
  const response = await fetch(path, { method: 'POST', credentials: 'same-origin', headers: { 'X-Requested-With': 'XMLHttpRequest' }, body })
  const data = await response.json().catch(() => null)
  if (!response.ok || !data || data.Ok === false) throw new Error(data?.Msg || '上传失败')
  return data
}
