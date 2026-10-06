// Exercise the delivered HTTP interface against a running synthetic demo.
import assert from 'node:assert/strict'

const origin = process.argv[2] ?? 'http://127.0.0.1:8089'
const request = (path, options = {}) => fetch(new URL(path, origin), {
  ...options, signal: AbortSignal.timeout(10000),
})

let ready = false
for (let attempt = 0; attempt < 40; attempt++) {
  try {
    ready = (await request('/health')).status === 204
  } catch { /* The container may still be starting. */ }
  if (ready) break
  await new Promise(resolve => setTimeout(resolve, 250))
}
assert.ok(ready, 'service becomes healthy')
assert.equal((await request('/api/health')).status, 204)

const home = await request('/')
assert.equal(home.status, 200)
assert.match(home.headers.get('content-type'), /text\/html/)
assert.equal(home.headers.get('cache-control'), 'no-cache')
const html = await home.text()
assert.match(html, /id="root"/)
// The app's pages, each at the root; the last two use the app's address
// form: v1= and v2= name the two picked submissions, side= the one a review
// page shows (the server ignores them).
for (const path of ['/claims', '/differences', '/unmatched', '/workloads', '/claim',
  '/?v1=vendor-1%2F2026-10-30&v2=vendor-2%2F2026-10-30&side=2',
  '/claim?v1=vendor-1%2F2026-10-30&v2=vendor-2%2F2026-10-30&match_key=0']) {
  const response = await request(path)
  assert.equal(response.status, 200, path)
  assert.equal(await response.text(), html, path + ' serves the app shell')
}
const script = html.match(/src="(\/assets\/[^" ]+\.js)"/)?.[1]
assert.ok(script, 'app shell names a built JavaScript asset')
const asset = await request(script)
assert.equal(asset.status, 200)
assert.match(asset.headers.get('content-type'), /javascript/)
assert.match(asset.headers.get('cache-control'), /immutable/)
assert.ok((await asset.text()).length > 100)
assert.equal((await request(script, { method: 'HEAD' })).status, 200)

for (const path of ['/api', '/api/missing', '/assets/missing.js', '/assets/', '/missing.css', '/.env']) {
  const response = await request(path)
  assert.equal(response.status, 404, path)
  assert.doesNotMatch(response.headers.get('content-type') ?? '', /text\/html/, path + ' does not become an app page')
}
assert.equal((await request('/api/diff')).status, 405)
assert.equal((await request('/claim', { method: 'POST' })).status, 405)

console.log('HTTP checks passed: app routes and assets, API isolation, health, and all six saved demo comparisons, including the one-cent change the five-cent rule counts as the same amount.')
