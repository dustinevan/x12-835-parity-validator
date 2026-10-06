// Compact loading and error states shared by the views.

export function Loading({ what = 'Loading' }: { what?: string }) {
  return (
    <div className="flex items-center gap-2 px-5 py-6 font-mono text-xs text-ink-2" role="status">
      <span className="h-2 w-2 animate-pulse rounded-full bg-ink-2" />
      {what}…
    </div>
  )
}

export function Failed({ error }: { error: Error }) {
  return (
    <div className="flex items-start gap-3 border-l-4 border-warn bg-warn/15 px-5 py-4" role="alert">
      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-warn text-xs font-bold text-black">!</span>
      <div>
        <div className="text-sm font-semibold text-ink">Could not load</div>
        <div className="font-mono text-[11px] text-ink-2">{error.message}</div>
      </div>
    </div>
  )
}

export function Empty({ children }: { children: string }) {
  return <div className="px-5 py-6 font-mono text-xs text-ink-2">{children}</div>
}

// Renders the first failure, else a loading line while any input is loading,
// else the empty message.
export function Pending({ errors, loading, empty = 'Nothing loaded yet' }: { errors: (Error | undefined)[]; loading: boolean; empty?: string }) {
  const error = errors.find((e) => e !== undefined)
  if (error) return <Failed error={error} />
  if (loading) return <Loading />
  return <Empty>{empty}</Empty>
}
