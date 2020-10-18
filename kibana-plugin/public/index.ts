import './index.scss';

import { ScorestackPlugin } from './plugin';

// This exports static code and TypeScript types,
// as well as, Kibana Platform `plugin()` initializer.
export function plugin() {
  return new ScorestackPlugin();
}
export { ScorestackPluginSetup, ScorestackPluginStart } from './types';
