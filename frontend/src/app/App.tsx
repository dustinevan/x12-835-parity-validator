import { SelectionProvider } from './selection'
import { Shell } from './Shell'

// The delivered validator app: the shell (Shell.tsx), with the two
// submission pickers and the sections in the left navigation, under the
// one selection every page and link shares (selection.tsx). The reference
// pages and the 835 compare page of example files are in the local docs
// site.
export default function App() {
  return (
    <SelectionProvider>
      <Shell />
    </SelectionProvider>
  )
}
