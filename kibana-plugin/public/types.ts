import { NavigationPublicPluginStart } from '../../../src/plugins/navigation/public';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface ScorestackPluginSetup {}
// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface ScorestackPluginStart {}

export interface AppPluginStartDependencies {
  navigation: NavigationPublicPluginStart;
}
