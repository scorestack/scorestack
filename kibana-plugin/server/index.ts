import { PluginInitializerContext } from 'kibana/server';
import { ScorestackPlugin } from './plugin';

//  This exports static code and TypeScript types,
//  as well as, Kibana Platform `plugin()` initializer.

export function plugin(initializerContext: PluginInitializerContext) {
  return new ScorestackPlugin(initializerContext);
}

export { ScorestackPluginSetup, ScorestackPluginStart } from './types';
