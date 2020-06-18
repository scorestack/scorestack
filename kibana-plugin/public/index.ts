import './index.scss';

import { ScoreStackPlugin } from './plugin';

// This exports static code and TypeScript types,
// as well as, Kibana Platform `plugin()` initializer.
export function plugin() {
  return new ScoreStackPlugin();
}
export { ScoreStackPluginSetup, ScoreStackPluginStart } from './types';
