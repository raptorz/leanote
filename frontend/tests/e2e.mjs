import { chromium } from '@playwright/test'
import { existsSync } from 'node:fs'

const baseURL = process.env.PEARLNOTE_E2E_URL
const username = process.env.PEARLNOTE_E2E_USERNAME
const password = process.env.PEARLNOTE_E2E_PASSWORD
if (!baseURL || !username || !password) {
  throw new Error('PEARLNOTE_E2E_URL, PEARLNOTE_E2E_USERNAME and PEARLNOTE_E2E_PASSWORD are required')
}

const executablePath = [process.env.PEARLNOTE_CHROMIUM_PATH, '/usr/bin/chromium', '/usr/bin/chromium-browser', '/usr/bin/google-chrome'].find(path => path && existsSync(path))
const browser = await chromium.launch({ headless: true, ...(executablePath ? { executablePath } : {}) })
const page = await browser.newPage()
page.setDefaultTimeout(15_000)
try {
  await page.goto(`${baseURL}/login`)
  await page.getByLabel('邮箱或用户名').fill(username)
  await page.getByLabel('密码').fill(password)
  await page.getByRole('button', { name: '继续' }).click()
  await page.waitForURL(/\/note/)

  await page.getByTitle('新建 Markdown').click()
  await page.getByLabel('文章标题').fill('Vue 端到端验收')
  await page.getByLabel('Markdown 正文').fill('# Pearlnote\n\nVue Web UI 正常保存。')
  await page.getByRole('button', { name: '保存', exact: true }).click()
  await page.getByText('已保存', { exact: true }).waitFor()
  await page.getByRole('button', { name: /Info/ }).click()
  await page.getByRole('heading', { name: '文章信息' }).waitFor()

  const attachmentName = `pearlnote-e2e-${Date.now()}.txt`
  await page.locator('input[type="file"]').nth(1).setInputFiles({
    name: attachmentName,
    mimeType: 'text/plain',
    buffer: Buffer.from('Pearlnote attachment E2E'),
  })
  await page.getByRole('heading', { name: '附件' }).waitFor()
  await page.getByRole('link', { name: attachmentName }).waitFor()
  const deleteResponse = page.waitForResponse(response => response.url().includes('/attach/deleteAttach'))
  page.once('dialog', dialog => dialog.accept())
  await page.getByRole('listitem').filter({ hasText: attachmentName }).getByRole('button', { name: '删除' }).click()
  const deleteResult = await (await deleteResponse).json()
  if (!deleteResult.Ok) throw new Error(`attachment delete failed: ${deleteResult.Msg}`)
  await page.getByRole('link', { name: attachmentName }).waitFor({ state: 'detached' })
  await page.getByRole('button', { name: '关闭' }).click()

  // Attachment mutations increment the note USN; saving afterwards verifies
  // that the UI refreshed its optimistic-concurrency metadata.
  await page.getByLabel('Markdown 正文').fill('# Pearlnote\n\n附件后继续保存正常。')
  await page.getByRole('button', { name: '保存', exact: true }).click()
  await page.getByText('已保存', { exact: true }).waitFor()

  await page.getByRole('link', { name: '账号' }).click()
  await page.getByRole('heading', { name: '账号管理' }).waitFor()
  await page.getByRole('link', { name: '管理' }).click()
  await page.getByRole('heading', { name: '系统管理' }).waitFor()

  const removed = await page.request.get(`${baseURL}/blog`)
  if (removed.status() !== 410) throw new Error(`blog returned ${removed.status()}, expected 410`)
  console.log('Pearlnote Vue Web E2E passed')
} finally {
  await browser.close()
}
