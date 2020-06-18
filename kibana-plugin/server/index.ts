import { PluginInitializerContext } from '../../../src/core/server';
import { ScoreStackPlugin } from './plugin';

//  This exports static code and TypeScript types,
//  as well as, Kibana Platform `plugin()` initializer.

export function plugin(initializerContext: PluginInitializerContext) {
  return new ScoreStackPlugin(initializerContext);
}

export { ScoreStackPluginSetup, ScoreStackPluginStart } from './types';
