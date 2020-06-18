import { NavigationPublicPluginStart } from '../../../src/plugins/navigation/public';

/* eslint-disable @typescript-eslint/no-empty-interface, prettier/prettier */
export interface ScoreStackPluginSetup { }
export interface ScoreStackPluginStart { }
/* eslint-enable @typescript-eslint/no-empty-interface, prettier/prettier */

export interface AppPluginStartDependencies {
  navigation: NavigationPublicPluginStart;
}
